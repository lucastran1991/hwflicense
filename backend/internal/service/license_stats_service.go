package service

import (
	"fmt"
	"time"

	"taskmaster-license/internal/models"
	"taskmaster-license/internal/repository"

	"github.com/google/uuid"
)

type LicenseStatsService struct {
	repo *repository.LicenseRepository
}

func NewLicenseStatsService(repo *repository.LicenseRepository) *LicenseStatsService {
	return &LicenseStatsService{repo: repo}
}

// SaveQuarterlyStats saves quarterly aggregated statistics
func (s *LicenseStatsService) SaveQuarterlyStats(stats *models.QuarterlyStats) error {
	// Validate period format (should be like Q1_2025, Q2_2025, etc.)
	if !isValidPeriodFormat(stats.Period) {
		return fmt.Errorf("invalid period format. Expected format: Q1_YYYY, Q2_YYYY, etc.")
	}
	
	// Set ID if not set
	if stats.ID == "" {
		stats.ID = uuid.New().String()
	}
	
	// Set created_at if not set
	if stats.CreatedAt.IsZero() {
		stats.CreatedAt = time.Now()
	}
	
	return s.repo.SaveQuarterlyStats(stats)
}

// isValidPeriodFormat validates the period format
func isValidPeriodFormat(period string) bool {
	// Should be in format Q1_2025, Q2_2025, Q3_2025, Q4_2025
	if len(period) < 4 {
		return false
	}
	
	if period[0] != 'Q' {
		return false
	}
	
	if period[1] < '1' || period[1] > '4' {
		return false
	}
	
	if period[2] != '_' {
		return false
	}
	
	// Check if year is valid (4 digits)
	if len(period) < 7 {
		return false
	}
	
	return true
}

