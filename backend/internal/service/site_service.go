package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"taskmaster-license/internal/models"
	"taskmaster-license/internal/repository"
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

	// Get or generate org key for signing
	// For now, we'll sign with a temporary approach (in production, load encrypted org key)
	// TODO: Load org key from database

	// Create site license data
	now := time.Now()
	licenseData := models.SiteLicenseData{
		Type:        "site_license",
		SiteID:      siteID,
		ParentCML:   orgID,
		ParentCMLSig: cml.Signature,
		Fingerprint: fingerprint,
		IssuedAt:    now.Format(time.RFC3339),
		ExpiresAt:   cml.Validity.Format(time.RFC3339),
		Features:    featurePacks,
	}

	licenseDataJSON, err := json.Marshal(licenseData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal license data: %w", err)
	}

	// TODO: Sign with org key
	// For now, create a placeholder signature
	signature := "TODO: sign with org key"

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

	// Check expiration
	expires, err := time.Parse(time.RFC3339, licenseData.ExpiresAt)
	if err == nil {
		now := time.Now()
		gracePeriod := now.AddDate(0, 0, 30) // 30-day grace period

		if expires.Before(now) && expires.Before(gracePeriod) {
			result.Valid = false
			result.Message = "License has expired"
		}
	}

	// TODO: Verify chain of trust
	// TODO: Verify signature
	// TODO: Check fingerprint matching (optional)

	return result, nil
}

type ValidationResult struct {
	Valid     bool     `json:"valid"`
	Message   string   `json:"message,omitempty"`
	Features  []string `json:"features"`
	ExpiresAt string   `json:"expires_at"`
}
