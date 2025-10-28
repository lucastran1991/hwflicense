package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"taskmaster-license/internal/models"
)

type LicenseRepository struct {
	db *sql.DB
}

func NewLicenseRepository(db *sql.DB) *LicenseRepository {
	return &LicenseRepository{db: db}
}

// CreateSiteKey creates a new site key
func (r *LicenseRepository) CreateSiteKey(siteKey *models.SiteKey) error {
	query := `INSERT INTO site_keys (id, site_id, enterprise_id, key_type, key_value, issued_at, expires_at, status, created_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := r.db.Exec(query, siteKey.ID, siteKey.SiteID, siteKey.EnterpriseID, siteKey.KeyType,
		siteKey.KeyValue, licenseTimeToString(siteKey.IssuedAt), licenseTimeToString(siteKey.ExpiresAt), siteKey.Status,
		licenseTimeToString(siteKey.CreatedAt))
	
	return err
}

// GetSiteKey retrieves a site key by site_id
func (r *LicenseRepository) GetSiteKey(siteID string) (*models.SiteKey, error) {
	query := `SELECT id, site_id, enterprise_id, key_type, key_value, issued_at, expires_at, status, last_validated, created_at 
			  FROM site_keys WHERE site_id = ?`
	
	var key models.SiteKey
	var issuedAt, expiresAt, createdAt, lastValidated sql.NullString
	
	err := r.db.QueryRow(query, siteID).Scan(
		&key.ID, &key.SiteID, &key.EnterpriseID, &key.KeyType, &key.KeyValue,
		&issuedAt, &expiresAt, &key.Status, &lastValidated, &createdAt)
	
	if err != nil {
		return nil, err
	}
	
	key.IssuedAt, _ = licenseStringToTime(issuedAt.String)
	key.ExpiresAt, _ = licenseStringToTime(expiresAt.String)
	key.CreatedAt, _ = licenseStringToTime(createdAt.String)
	if lastValidated.Valid {
		t, _ := licenseStringToTime(lastValidated.String)
		key.LastValidated = &t
	}
	
	return &key, nil
}

// ListSiteKeys retrieves all site keys, optionally filtered by enterprise_id
func (r *LicenseRepository) ListSiteKeys(enterpriseID string) ([]*models.SiteKey, error) {
	var query string
	var args []interface{}
	
	if enterpriseID != "" {
		query = `SELECT id, site_id, enterprise_id, key_type, key_value, issued_at, expires_at, status, 
				 last_validated, created_at FROM site_keys WHERE enterprise_id = ? ORDER BY created_at DESC`
		args = []interface{}{enterpriseID}
	} else {
		query = `SELECT id, site_id, enterprise_id, key_type, key_value, issued_at, expires_at, status, 
				 last_validated, created_at FROM site_keys ORDER BY created_at DESC`
	}
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var keys []*models.SiteKey
	for rows.Next() {
		var key models.SiteKey
		var issuedAt, expiresAt, createdAt, lastValidated sql.NullString
		
		err := rows.Scan(&key.ID, &key.SiteID, &key.EnterpriseID, &key.KeyType, &key.KeyValue,
			&issuedAt, &expiresAt, &key.Status, &lastValidated, &createdAt)
		if err != nil {
			return nil, err
		}
		
		key.IssuedAt, _ = licenseStringToTime(issuedAt.String)
		key.ExpiresAt, _ = licenseStringToTime(expiresAt.String)
		key.CreatedAt, _ = licenseStringToTime(createdAt.String)
		if lastValidated.Valid {
			t, _ := licenseStringToTime(lastValidated.String)
			key.LastValidated = &t
		}
		
		keys = append(keys, &key)
	}
	
	return keys, nil
}

// UpdateSiteKey updates a site key's status or other fields
func (r *LicenseRepository) UpdateSiteKey(siteID string, updates map[string]interface{}) error {
	query := `UPDATE site_keys SET `
	var args []interface{}
	var setParts []string
	
	for field, value := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = ?", field))
		args = append(args, value)
	}
	
	if len(setParts) == 0 {
		return fmt.Errorf("no updates provided")
	}
	
	query += fmt.Sprintf("%s WHERE site_id = ?", joinStrings(setParts, ", "))
	args = append(args, siteID)
	
	_, err := r.db.Exec(query, args...)
	return err
}

// RefreshSiteKey creates a new key for an existing site and logs the refresh
func (r *LicenseRepository) RefreshSiteKey(siteID, oldKey, newKey string, newExpiresAt time.Time) error {
	// Update the site key with new key value and expiration
	query := `UPDATE site_keys SET key_value = ?, expires_at = ?, issued_at = ?, status = 'active' 
			  WHERE site_id = ? AND key_value = ?`
	
	result, err := r.db.Exec(query, newKey, licenseTimeToString(newExpiresAt), licenseTimeToString(time.Now()), siteID, oldKey)
	if err != nil {
		return err
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("site key not found or old key mismatch")
	}
	
	// Log the refresh
	logQuery := `INSERT INTO key_refresh_log (id, site_id, old_key, new_key, refreshed_at, reason) 
				 VALUES (?, ?, ?, ?, ?, ?)`
	_, err = r.db.Exec(logQuery, fmt.Sprintf("log_%d", time.Now().UnixNano()), siteID, oldKey, newKey,
		licenseTimeToString(time.Now()), "Monthly refresh")
	
	return err
}

// ValidateSiteKey checks if a key is valid
func (r *LicenseRepository) ValidateSiteKey(siteID, key string) (*models.SiteKey, error) {
	query := `SELECT id, site_id, enterprise_id, key_type, key_value, issued_at, expires_at, status, 
			  last_validated, created_at FROM site_keys 
			  WHERE site_id = ? AND key_value = ? AND status = 'active'`
	
	var siteKey models.SiteKey
	var issuedAt, expiresAt, createdAt, lastValidated sql.NullString
	
	err := r.db.QueryRow(query, siteID, key).Scan(
		&siteKey.ID, &siteKey.SiteID, &siteKey.EnterpriseID, &siteKey.KeyType, &siteKey.KeyValue,
		&issuedAt, &expiresAt, &siteKey.Status, &lastValidated, &createdAt)
	
	if err != nil {
		return nil, err
	}
	
	siteKey.IssuedAt, _ = licenseStringToTime(issuedAt.String)
	siteKey.ExpiresAt, _ = licenseStringToTime(expiresAt.String)
	siteKey.CreatedAt, _ = licenseStringToTime(createdAt.String)
	if lastValidated.Valid {
		t, _ := licenseStringToTime(lastValidated.String)
		siteKey.LastValidated = &t
	}
	
	// Check if key has expired
	if siteKey.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("key has expired")
	}
	
	// Update last_validated
	r.db.Exec(`UPDATE site_keys SET last_validated = ? WHERE site_id = ?`,
		licenseTimeToString(time.Now()), siteID)
	
	return &siteKey, nil
}

// SaveQuarterlyStats saves quarterly statistics
func (r *LicenseRepository) SaveQuarterlyStats(stats *models.QuarterlyStats) error {
	query := `INSERT INTO quarterly_stats (id, period, production_sites, dev_sites, user_counts, enterprise_breakdown, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?)
			  ON CONFLICT(period) DO UPDATE SET 
			  production_sites = excluded.production_sites,
			  dev_sites = excluded.dev_sites,
			  user_counts = excluded.user_counts,
			  enterprise_breakdown = excluded.enterprise_breakdown`
	
	userCountsJSON, _ := json.Marshal(stats.UserCounts)
	breakdownJSON, _ := json.Marshal(stats.EnterpriseBreakdown)
	
	_, err := r.db.Exec(query, stats.ID, stats.Period, stats.ProductionSites, stats.DevSites,
		string(userCountsJSON), string(breakdownJSON), licenseTimeToString(stats.CreatedAt))
	
	return err
}

// SaveAlert saves an alert
func (r *LicenseRepository) SaveAlert(alert *models.Alert) error {
	query := `INSERT INTO alerts (id, site_id, alert_type, message, alert_timestamp, sent_to_astack, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	_, err := r.db.Exec(query, alert.ID, alert.SiteID, alert.AlertType, alert.Message,
		licenseTimeToString(alert.Timestamp), alert.SentToAStack, licenseTimeToString(alert.CreatedAt))
	
	return err
}

// Helper functions
func licenseTimeToString(t time.Time) string {
	return t.Format(time.RFC3339)
}

func licenseStringToTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

