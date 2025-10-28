package service

import (
	"fmt"
	"time"

	"taskmaster-license/internal/models"
	"taskmaster-license/internal/repository"

	"github.com/google/uuid"
)

type LicenseSiteService struct {
	repo *repository.LicenseRepository
}

func NewLicenseSiteService(repo *repository.LicenseRepository) *LicenseSiteService {
	return &LicenseSiteService{repo: repo}
}

// CreateSiteKey creates a new site key
func (s *LicenseSiteService) CreateSiteKey(siteID, enterpriseID, mode, orgID string) (*models.SiteKey, error) {
	// Generate unique key value
	keyValue := generateKeyValue()
	
	// Calculate expiration (30 days from now)
	issuedAt := time.Now()
	expiresAt := issuedAt.AddDate(0, 0, 30)
	
	siteKey := &models.SiteKey{
		ID:           uuid.New().String(),
		SiteID:       siteID,
		EnterpriseID: enterpriseID,
		KeyType:      mode, // "production" or "dev"
		KeyValue:     keyValue,
		IssuedAt:     issuedAt,
		ExpiresAt:     expiresAt,
		Status:       "active",
		CreatedAt:    time.Now(),
	}
	
	if err := s.repo.CreateSiteKey(siteKey); err != nil {
		return nil, fmt.Errorf("failed to create site key: %w", err)
	}
	
	return siteKey, nil
}

// RefreshSiteKey refreshes an existing site key
func (s *LicenseSiteService) RefreshSiteKey(siteID, oldKey string) (*models.SiteKey, error) {
	// Validate old key
	siteKey, err := s.repo.GetSiteKey(siteID)
	if err != nil {
		return nil, fmt.Errorf("site key not found: %w", err)
	}
	
	if siteKey.KeyValue != oldKey {
		return nil, fmt.Errorf("old key mismatch")
	}
	
	if siteKey.Status != "active" {
		return nil, fmt.Errorf("key is not active")
	}
	
	// Generate new key
	newKey := generateKeyValue()
	newExpiresAt := time.Now().AddDate(0, 0, 30)
	
	// Update the key
	if err := s.repo.RefreshSiteKey(siteID, oldKey, newKey, newExpiresAt); err != nil {
		return nil, fmt.Errorf("failed to refresh key: %w", err)
	}
	
	// Get updated key
	updatedKey, err := s.repo.GetSiteKey(siteID)
	if err != nil {
		return nil, err
	}
	
	return updatedKey, nil
}

// ValidateSiteKey validates a key and returns validation info
func (s *LicenseSiteService) ValidateSiteKey(siteID, key string) (*models.SiteKey, error) {
	siteKey, err := s.repo.ValidateSiteKey(siteID, key)
	if err != nil {
		return nil, fmt.Errorf("key validation failed: %w", err)
	}
	
	// Check expiration
	if siteKey.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("key has expired")
	}
	
	return siteKey, nil
}

// UpdateSiteKeyStatus updates a site key's status
func (s *LicenseSiteService) UpdateSiteKeyStatus(siteID, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	return s.repo.UpdateSiteKey(siteID, updates)
}

// ListSiteKeys lists all site keys, optionally filtered by enterprise
func (s *LicenseSiteService) ListSiteKeys(enterpriseID string) ([]*models.SiteKey, error) {
	return s.repo.ListSiteKeys(enterpriseID)
}

// Helper function to generate unique key value
func generateKeyValue() string {
	return fmt.Sprintf("LS-%s", uuid.New().String())
}

