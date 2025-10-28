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

type CMLService struct {
	repo *repository.Repository
}

func NewCMLService(repo *repository.Repository) *CMLService {
	return &CMLService{repo: repo}
}

func (s *CMLService) UploadCML(cmlDataStr, signature string, publicKeyPEM string) (*models.CML, error) {
	// Parse CML data
	var cmlData models.CMLData
	if err := json.Unmarshal([]byte(cmlDataStr), &cmlData); err != nil {
		return nil, fmt.Errorf("failed to parse CML data: %w", err)
	}

	// Verify signature with the provided public key
	publicKey, err := crypto.PEMToPublicKey(publicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	// Verify signature
	valid, err := crypto.VerifySignature([]byte(cmlDataStr), signature, publicKey)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid CML signature: %w", err)
	}

	// Create CML entity
	cml := &models.CML{
		ID:            uuid.New().String(),
		OrgID:         cmlData.OrgID,
		MaxSites:      cmlData.MaxSites,
		Validity:      parseTime(cmlData.Validity),
		FeaturePacks:  cmlData.FeaturePacks,
		DevKeyPublic:  publicKeyPEM,
		ProdKeyPublic: publicKeyPEM,
		CMLData:       json.RawMessage(cmlDataStr),
		Signature:     signature,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Store in database
	if err := s.repo.CreateCML(cml); err != nil {
		return nil, fmt.Errorf("failed to store CML: %w", err)
	}

	return cml, nil
}

func (s *CMLService) GetCML(orgID string) (*models.CML, error) {
	cml, err := s.repo.GetCML(orgID)
	if err != nil {
		// Return default CML if not found in database
		return s.createDefaultCML(orgID), nil
	}
	return cml, nil
}

func (s *CMLService) createDefaultCML(orgID string) *models.CML {
	// Calculate validity date (1 year from now)
	validity := time.Now().AddDate(1, 0, 0)
	
	// Get default values from config
	maxSites := config.AppConfig.DefaultMaxSites
	if maxSites == 0 {
		maxSites = 100 // Fallback default
	}
	
	featurePacks := config.AppConfig.DefaultFeaturePacks
	if len(featurePacks) == 0 {
		featurePacks = []string{"basic", "standard"}
	}
	
	// Try to get org's public key if available
	publicKey := ""
	orgKey, err := s.repo.GetOrgKey(orgID, "dev")
	if err == nil && orgKey != nil {
		publicKey = orgKey.PublicKey
	}
	
	// Create default CML data
	cmlData := models.CMLData{
		Type:            "default_cml",
		OrgID:           orgID,
		MaxSites:        maxSites,
		Validity:        validity.Format(time.RFC3339),
		FeaturePacks:    featurePacks,
		KeyType:         "dev",
		IssuedBy:        "system",
		IssuerPublicKey: publicKey,
		IssuedAt:        time.Now().Format(time.RFC3339),
	}
	
	cmlDataJSON, _ := json.Marshal(cmlData)
	
	// Create in-memory CML (not persisted to database)
	return &models.CML{
		ID:            fmt.Sprintf("default_%s", orgID),
		OrgID:         orgID,
		MaxSites:      maxSites,
		Validity:      validity,
		FeaturePacks:  featurePacks,
		DevKeyPublic:  publicKey,
		ProdKeyPublic: publicKey,
		CMLData:       json.RawMessage(cmlDataJSON),
		Signature:     "default_cml_no_signature",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func (s *CMLService) RefreshCML(orgID string, cmlDataStr, signature string) error {
	cml, err := s.repo.GetCML(orgID)
	if err != nil {
		return err
	}

	// Update CML data
	cml.CMLData = json.RawMessage(cmlDataStr)
	cml.Signature = signature
	cml.UpdatedAt = time.Now()

	// Parse and update validity
	var cmlData models.CMLData
	if err := json.Unmarshal([]byte(cmlDataStr), &cmlData); err == nil {
		cml.Validity = parseTime(cmlData.Validity)
		cml.MaxSites = cmlData.MaxSites
		cml.FeaturePacks = cmlData.FeaturePacks
	}

	return s.repo.UpdateCML(cml)
}

func parseTime(timeStr string) time.Time {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		// Try alternative format
		t, _ = time.Parse("2006-01-02", timeStr)
	}
	return t
}
