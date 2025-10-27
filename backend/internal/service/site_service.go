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

type SiteService struct {
	repo       *repository.Repository
	cmlService *CMLService
}

func NewSiteService(repo *repository.Repository, cmlService *CMLService) *SiteService {
	return &SiteService{
		repo:       repo,
		cmlService: cmlService,
	}
}

func (s *SiteService) CreateSiteLicense(orgID, siteID string, fingerprint map[string]interface{}) (*models.SiteLicense, error) {
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

	// Parse feature packs from CML
	var featurePacks []string
	var cmlData models.CMLData
	if err := json.Unmarshal(cml.CMLData, &cmlData); err == nil {
		featurePacks = cmlData.FeaturePacks
	}

	// Get org key for signing
	orgKey, err := s.repo.GetOrgKey(orgID, "dev") // Default to dev key
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

	// Create site license data
	now := time.Now()
	licenseData := models.SiteLicenseData{
		Type:         "site_license",
		SiteID:       siteID,
		ParentCML:    orgID,
		ParentCMLSig: cml.Signature,
		Fingerprint:  fingerprint,
		IssuedAt:     now.Format(time.RFC3339),
		ExpiresAt:    cml.Validity.Format(time.RFC3339),
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

	// Create site license entity
	siteLicense := &models.SiteLicense{
		ID:           uuid.New().String(),
		SiteID:       siteID,
		OrgID:        orgID,
		Fingerprint:  json.RawMessage(fingerprintJSON),
		LicenseData:  json.RawMessage(licenseDataJSON),
		Signature:    signature,
		IssuedAt:     now,
		Status:       "active",
		CreatedAt:    now,
	}

	// Store in database
	if err := s.repo.CreateSiteLicense(siteLicense); err != nil {
		return nil, fmt.Errorf("failed to create site license: %w", err)
	}

	return siteLicense, nil
}

func (s *SiteService) GetSiteLicense(siteID string) (*models.SiteLicense, error) {
	return s.repo.GetSiteLicense(siteID)
}

func (s *SiteService) ListSiteLicenses(orgID, status string, limit, offset int) ([]models.SiteLicense, int, error) {
	return s.repo.ListSiteLicenses(orgID, status, limit, offset)
}

func (s *SiteService) UpdateHeartbeat(siteID string) error {
	return s.repo.UpdateSiteHeartbeat(siteID)
}

func (s *SiteService) RevokeSiteLicense(siteID string) error {
	return s.repo.RevokeSiteLicense(siteID)
}

func (s *SiteService) ValidateLicense(licenseDataStr string, fingerprint map[string]interface{}) (*ValidationResult, error) {
	// Parse license data
	var licenseData models.SiteLicenseData
	if err := json.Unmarshal([]byte(licenseDataStr), &licenseData); err != nil {
		return nil, fmt.Errorf("invalid license data: %w", err)
	}

	result := &ValidationResult{
		Valid:     true,
		ExpiresAt: licenseData.ExpiresAt,
		Features:  licenseData.Features,
	}

	// Check expiration with 30-day grace period
	expires, err := time.Parse(time.RFC3339, licenseData.ExpiresAt)
	if err == nil {
		now := time.Now()
		gracePeriod := now.AddDate(0, 0, 30) // 30-day grace period

		if expires.Before(now) {
			if expires.Before(gracePeriod.AddDate(0, 0, -30)) {
				result.Valid = false
				result.Message = "License has expired and grace period has passed"
			} else {
				result.Message = "License expired but within grace period (30 days)"
			}
		}
	}

	// Verify signature chain of trust
	if !s.verifySignatureChain(licenseData) {
		result.Valid = false
		if result.Message == "" {
			result.Message = "Signature chain verification failed"
		} else {
			result.Message += "; signature chain verification failed"
		}
	}

	// Check fingerprint matching (optional but recommended)
	if fingerprint != nil && len(fingerprint) > 0 {
		if !s.matchFingerprint(licenseData.Fingerprint, fingerprint) {
			result.Message = "Fingerprint mismatch detected"
			// Note: This doesn't invalidate the license, just logs a warning
		}
	}

	return result, nil
}

// verifySignatureChain verifies the signature chain: Site License → CML → Root
func (s *SiteService) verifySignatureChain(licenseData models.SiteLicenseData) bool {
	// Get CML for parent validation
	cml, err := s.cmlService.GetCML(licenseData.ParentCML)
	if err != nil {
		return false
	}

	// Verify site signature with org public key
	// Get org key to get public key
	orgKey, err := s.repo.GetOrgKey(licenseData.ParentCML, "dev") // Default to dev
	if err != nil {
		return false
	}

	// Load public key from PEM
	publicKey, err := crypto.LoadPublicKeyFromPEM(orgKey.PublicKey)
	if err != nil {
		return false
	}

	// For now, we need to get the signature from the database
	// Get site license to access signature
	siteLicense, err := s.repo.GetSiteLicense(licenseData.SiteID)
	if err != nil {
		return false
	}

	// Verify site signature
	valid, err := crypto.VerifyJSON(licenseData, siteLicense.Signature, publicKey)
	if err != nil || !valid {
		return false
	}

	// Verify CML signature with root public key (if root key is configured)
	if config.AppConfig.RootPublicKey != "" {
		rootPublicKey, err := crypto.LoadPublicKeyFromPEM(config.AppConfig.RootPublicKey)
		if err == nil {
			var cmlData models.CMLData
			if err := json.Unmarshal(cml.CMLData, &cmlData); err == nil {
				cmlValid, err := crypto.VerifyJSON(cmlData, cml.Signature, rootPublicKey)
				if err != nil || !cmlValid {
					return false
				}
			}
		}
	}

	return true
}

// matchFingerprint checks if the license fingerprint matches the provided fingerprint
func (s *SiteService) matchFingerprint(licenseFingerprint map[string]interface{}, providedFingerprint map[string]interface{}) bool {
	for key, value := range providedFingerprint {
		licenseValue, exists := licenseFingerprint[key]
		if !exists || fmt.Sprintf("%v", licenseValue) != fmt.Sprintf("%v", value) {
			return false
		}
	}
	return true
}

type ValidationResult struct {
	Valid     bool     `json:"valid"`
	Message   string   `json:"message,omitempty"`
	Features  []string `json:"features"`
	ExpiresAt string   `json:"expires_at"`
}
