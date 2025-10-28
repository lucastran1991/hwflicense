package service

import (
	"fmt"
	"time"

	"taskmaster-license/internal/models"
	"taskmaster-license/internal/repository"

	"github.com/google/uuid"
)

type LicenseAlertService struct {
	repo *repository.LicenseRepository
}

func NewLicenseAlertService(repo *repository.LicenseRepository) *LicenseAlertService {
	return &LicenseAlertService{repo: repo}
}

// SaveAlert saves an alert
func (s *LicenseAlertService) SaveAlert(alert *models.Alert) error {
	// Validate alert type
	validTypes := []string{"key_expired", "key_invalid"}
	if !containsString(validTypes, alert.AlertType) {
		return fmt.Errorf("invalid alert type: %s. Valid types are: key_expired, key_invalid", alert.AlertType)
	}
	
	// Set ID if not set
	if alert.ID == "" {
		alert.ID = uuid.New().String()
	}
	
	// Set created_at if not set
	if alert.CreatedAt.IsZero() {
		alert.CreatedAt = time.Now()
	}
	
	// Set timestamp if not set
	if alert.Timestamp.IsZero() {
		alert.Timestamp = time.Now()
	}
	
	return s.repo.SaveAlert(alert)
}

// containsString checks if a slice contains a string
func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

