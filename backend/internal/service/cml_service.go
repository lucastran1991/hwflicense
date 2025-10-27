package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

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
	return s.repo.GetCML(orgID)
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
