package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"taskmaster-license/internal/config"
	"taskmaster-license/internal/models"
	"taskmaster-license/internal/repository"
	"taskmaster-license/pkg/crypto"
)

type ManifestService struct {
	repo  *repository.Repository
	orgID string
}

func NewManifestService(repo *repository.Repository) *ManifestService {
	return &ManifestService{
		repo:  repo,
		orgID: "default", // In production, get from context
	}
}

func (s *ManifestService) GenerateManifest(period string) (*models.UsageManifest, error) {
	// Get active sites for the period
	sites, _, err := s.repo.ListSiteLicenses(s.orgID, "active", 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sites: %w", err)
	}

	// Get active site count
	activeCount, err := s.repo.GetActiveSiteCount(s.orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active site count: %w", err)
	}

	// Generate stats (mock data for now)
	stats := map[string]interface{}{
		"users": map[string]interface{}{
			"admin_users": 5,
			"total_users":  100,
		},
		"sites": map[string]interface{}{
			"boost": map[string]interface{}{
				"active": len(sites) / 2,
				"basic":  len(sites) / 4,
			},
			"hwf": map[string]interface{}{
				"active": len(sites) / 2,
				"basic":  len(sites) / 4,
			},
		},
		"total_active_sites": activeCount,
	}

	// Create manifest data
	manifestData := map[string]interface{}{
		"type":         "usage_manifest",
		"org_id":       s.orgID,
		"period":       period,
		"generated_at": time.Now().Format(time.RFC3339),
		"stats":        stats,
		"active_sites": s.formatActiveSites(sites),
	}

	manifestJSON, err := json.Marshal(manifestData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal manifest: %w", err)
	}

	// Get org key for signing
	orgKey, err := s.repo.GetOrgKey(s.orgID, "dev") // Default to dev key
	if err != nil {
		return nil, fmt.Errorf("failed to get org key: %w", err)
	}

	// Decrypt private key
	privateKeyPEM, err := crypto.DecryptPrivateKey(orgKey.PrivateKeyEncrypted, config.AppConfig.EncryptionPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt org key: %w", err)
	}

	// Load private key
	privateKey, err := crypto.LoadPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	// Sign manifest with org private key
	signature, err := crypto.SignJSON(manifestData, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign manifest: %w", err)
	}

	// Create manifest entity
	manifest := &models.UsageManifest{
		ID:           uuid.New().String(),
		OrgID:        s.orgID,
		Period:       period,
		ManifestData: manifestJSON,
		Signature:    signature,
		SentToAStack: false,
		CreatedAt:    time.Now(),
	}

	// Store in database
	if err := s.repo.CreateManifest(manifest); err != nil {
		return nil, fmt.Errorf("failed to store manifest: %w", err)
	}

	return manifest, nil
}

func (s *ManifestService) formatActiveSites(sites []models.SiteLicense) []map[string]interface{} {
	activeSites := []map[string]interface{}{}

	for _, site := range sites {
		licenseData := models.SiteLicenseData{}
		json.Unmarshal(site.LicenseData, &licenseData)

		siteInfo := map[string]interface{}{
			"site_id":   site.SiteID,
			"issued_at": licenseData.IssuedAt,
			"last_seen": s.formatLastSeen(site.LastSeen),
			"status":    site.Status,
		}

		activeSites = append(activeSites, siteInfo)
	}

	return activeSites
}

func (s *ManifestService) formatLastSeen(lastSeen *time.Time) string {
	if lastSeen == nil {
		return time.Now().Format(time.RFC3339)
	}
	return lastSeen.Format(time.RFC3339)
}

func (s *ManifestService) GetManifest(id string) (*models.UsageManifest, error) {
	return s.repo.GetManifest(id)
}

func (s *ManifestService) ListManifests(period string) ([]models.UsageManifest, error) {
	return s.repo.ListManifests(s.orgID, period)
}

func (s *ManifestService) MarkManifestSent(id string) error {
	now := time.Now()
	return s.repo.UpdateManifestSent(id, true, &now)
}
