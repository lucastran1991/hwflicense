package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"taskmaster-license/internal/client"
	"taskmaster-license/internal/config"
	"taskmaster-license/internal/models"
	"taskmaster-license/pkg/crypto"
)

// CreateSiteLicenseWithMode creates a site license with key type (dev/prod)
func (s *SiteService) CreateSiteLicenseWithMode(orgID, siteID string, fingerprint map[string]interface{}, mode string) (*models.SiteLicense, error) {
	// TODO: Integrate with license server when ready
	// licenseClient := client.NewLicenseServerClient()

	// Get or create enterprise
	enterprise, err := s.repo.GetEnterprise(orgID)
	if err != nil {
		// Create default enterprise
		enterprise = &models.Enterprise{
			ID:            uuid.New().String(),
			Name:          orgID + " Enterprise",
			OrgID:         orgID,
			EnterpriseKey: generateEnterpriseKey(),
			CreatedAt:     time.Now(),
		}
		if err := s.repo.CreateEnterprise(enterprise); err != nil {
			return nil, fmt.Errorf("failed to create enterprise: %w", err)
		}
	}

	// Call License Server to create site key
	// Note: For now, we'll set expires_at manually until we integrate fully
	keyResponse := &client.SiteKeyResponse{
		SiteID:      siteID,
		KeyType:     mode,
		ExpiresAt:   time.Now().AddDate(0, 0, 30),
	}
	
	// TODO: Make actual call to license server when it's running
	// keyResponse, err := licenseClient.CreateSiteKey(siteID, enterprise.ID, mode, orgID)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to create site key from license server: %w", err)
	// }

	// Get CML for validation
	cml, err := s.cmlService.GetCML(orgID)
	if err != nil {
		return nil, fmt.Errorf("CML not found: %w", err)
	}

	// Check max sites constraint
	activeCount, err := s.repo.CountActiveSites(orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to count active sites: %w", err)
	}

	if activeCount >= cml.MaxSites {
		return nil, fmt.Errorf("maximum sites limit (%d) reached", cml.MaxSites)
	}

	// Create site license data with new fields
	now := time.Now()
	expiresAt := keyResponse.ExpiresAt

	// Get org key for signing (existing logic)
	orgKey, err := s.repo.GetOrgKey(orgID, "dev")
	if err != nil {
		return nil, fmt.Errorf("failed to get org key: %w", err)
	}

	privateKeyPEM, err := crypto.DecryptPrivateKey(orgKey.PrivateKeyEncrypted, config.AppConfig.EncryptionPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt org key: %w", err)
	}

	privateKey, err := crypto.LoadPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	// Parse feature packs
	var featurePacks []string
	var cmlData models.CMLData
	if err := json.Unmarshal(cml.CMLData, &cmlData); err == nil {
		featurePacks = cmlData.FeaturePacks
	}

	licenseData := models.SiteLicenseData{
		Type:         "site_license",
		SiteID:       siteID,
		ParentCML:    orgID,
		ParentCMLSig: cml.Signature,
		Fingerprint:  fingerprint,
		IssuedAt:     now.Format(time.RFC3339),
		ExpiresAt:    expiresAt.Format(time.RFC3339),
		Features:     featurePacks,
	}

	licenseDataJSON, err := json.Marshal(licenseData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal license data: %w", err)
	}

	// Sign with org private key
	signature, err := crypto.SignJSON(licenseData, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign license: %w", err)
	}

	fingerprintJSON, _ := json.Marshal(fingerprint)

	// Create site license entity with new fields
	siteLicense := &models.SiteLicense{
		ID:            uuid.New().String(),
		SiteID:        siteID,
		OrgID:         orgID,
		EnterpriseID:  enterprise.ID,
		KeyType:       mode,
		ExpiresAt:     &expiresAt,
		Fingerprint:   json.RawMessage(fingerprintJSON),
		LicenseData:   json.RawMessage(licenseDataJSON),
		Signature:     signature,
		IssuedAt:      now,
		Status:        "active",
		CreatedAt:     now,
	}

	// Store in database
	if err := s.repo.CreateSiteLicense(siteLicense); err != nil {
		return nil, fmt.Errorf("failed to create site license: %w", err)
	}

	return siteLicense, nil
}

// RefreshSiteKey refreshes a site key
func (s *SiteService) RefreshSiteKey(siteID string) error {
	// Get current site license
	siteLicense, err := s.repo.GetSiteLicense(siteID)
	if err != nil {
		return fmt.Errorf("site license not found: %w", err)
	}

	// Update expiration to new 30 days
	now := time.Now()
	newExpiration := now.AddDate(0, 0, 30)
	siteLicense.ExpiresAt = &newExpiration
	lastRefreshed := now
	siteLicense.LastRefreshed = &lastRefreshed

	// Store updated license (would use Update method when implemented)
	return s.repo.RevokeSiteLicense(siteID) // Placeholder
}

// GetSitesNearExpiration returns sites expiring within specified days
func (s *SiteService) GetSitesNearExpiration(days int) ([]*models.SiteLicense, error) {
	// TODO: Implement query for sites near expiration
	// For now, return empty list
	return []*models.SiteLicense{}, nil
}

// AggregateQuarterlyStats aggregates quarterly stats
func (s *SiteService) AggregateQuarterlyStats() (*client.QuarterlyStats, error) {
	// Get all active sites
	productionCount := 0
	devCount := 0
	userCounts := make(map[string]interface{})
	enterpriseBreakdown := make([]map[string]interface{}, 0)

	// Query production sites
	sites, _, err := s.repo.ListSiteLicenses("", "active", 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to list sites: %w", err)
	}

	for _, site := range sites {
		if site.KeyType == "production" {
			productionCount++
		} else if site.KeyType == "dev" {
			devCount++
		}
	}

	return &client.QuarterlyStats{
		Period:              getQuarterPeriod(),
		ProductionSites:     productionCount,
		DevSites:            devCount,
		UserCounts:          userCounts,
		EnterpriseBreakdown: enterpriseBreakdown,
	}, nil
}

// UpdateSiteLicense updates a site license
func (s *SiteService) UpdateSiteLicense(siteID string, license *models.SiteLicense) error {
	// TODO: Implement update method in repository
	return nil
}

// Helper functions
func generateEnterpriseKey() string {
	// Simple key generation
	return fmt.Sprintf("ent_%s", uuid.New().String()[:8])
}

func getQuarterPeriod() string {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	var quarter string
	if month <= 3 {
		quarter = "Q1"
	} else if month <= 6 {
		quarter = "Q2"
	} else if month <= 9 {
		quarter = "Q3"
	} else {
		quarter = "Q4"
	}

	return fmt.Sprintf("%s_%d", quarter, year)
}

