package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"taskmaster-license/internal/models"
)

func (r *Repository) CreateSiteLicense(site *models.SiteLicense) error {
	query := `INSERT INTO site_licenses (id, site_id, org_id, fingerprint, license_data, 
		signature, issued_at, last_seen, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var lastSeen sql.NullString
	if site.LastSeen != nil {
		lastSeen = sql.NullString{String: timeToString(*site.LastSeen), Valid: true}
	}

	_, err := r.db.Connection.Exec(query,
		site.ID,
		site.SiteID,
		site.OrgID,
		string(site.Fingerprint),
		string(site.LicenseData),
		site.Signature,
		timeToString(site.IssuedAt),
		lastSeen,
		site.Status,
		timeToString(site.CreatedAt),
	)
	return err
}

func (r *Repository) GetSiteLicense(siteID string) (*models.SiteLicense, error) {
	query := `SELECT id, site_id, org_id, fingerprint, license_data, signature, 
		issued_at, last_seen, status, created_at FROM site_licenses WHERE site_id = ?`

	var site models.SiteLicense
	var issuedAtStr, createdAtStr string
	var fingerprintStr, licenseDataStr string
	var lastSeen sql.NullTime

	err := r.db.Connection.QueryRow(query, siteID).Scan(
		&site.ID,
		&site.SiteID,
		&site.OrgID,
		&fingerprintStr,
		&licenseDataStr,
		&site.Signature,
		&issuedAtStr,
		&lastSeen,
		&site.Status,
		&createdAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("site license not found: %s", siteID)
	}
	if err != nil {
		return nil, err
	}

	// Parse timestamps
	site.IssuedAt, _ = stringToTime(issuedAtStr)
	site.CreatedAt, _ = stringToTime(createdAtStr)
	if lastSeen.Valid {
		t := lastSeen.Time
		site.LastSeen = &t
	}

	// Parse JSON fields
	site.Fingerprint = json.RawMessage(fingerprintStr)
	site.LicenseData = json.RawMessage(licenseDataStr)

	return &site, nil
}

func (r *Repository) ListSiteLicenses(orgID, status string, limit, offset int) ([]models.SiteLicense, int, error) {
	// Build query with filters
	query := `SELECT id, site_id, org_id, fingerprint, license_data, signature, 
		issued_at, last_seen, status, created_at FROM site_licenses WHERE 1=1`
	args := []interface{}{}

	if orgID != "" {
		query += " AND org_id = ?"
		args = append(args, orgID)
	}
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Connection.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sites []models.SiteLicense
	for rows.Next() {
		var site models.SiteLicense
		var issuedAtStr, createdAtStr string
		var fingerprintStr, licenseDataStr string
		var lastSeen sql.NullTime

		err := rows.Scan(
			&site.ID,
			&site.SiteID,
			&site.OrgID,
			&fingerprintStr,
			&licenseDataStr,
			&site.Signature,
			&issuedAtStr,
			&lastSeen,
			&site.Status,
			&createdAtStr,
		)
		if err != nil {
			return nil, 0, err
		}

		site.IssuedAt, _ = stringToTime(issuedAtStr)
		site.CreatedAt, _ = stringToTime(createdAtStr)
		if lastSeen.Valid {
			t := lastSeen.Time
			site.LastSeen = &t
		}
		site.Fingerprint = json.RawMessage(fingerprintStr)
		site.LicenseData = json.RawMessage(licenseDataStr)

		sites = append(sites, site)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM site_licenses WHERE 1=1`
	countArgs := []interface{}{}
	if orgID != "" {
		countQuery += " AND org_id = ?"
		countArgs = append(countArgs, orgID)
	}
	if status != "" {
		countQuery += " AND status = ?"
		countArgs = append(countArgs, status)
	}

	var total int
	err = r.db.Connection.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return sites, total, nil
}

func (r *Repository) UpdateSiteHeartbeat(siteID string) error {
	query := `UPDATE site_licenses SET last_seen = ? WHERE site_id = ?`
	now := time.Now()
	_, err := r.db.Connection.Exec(query, timeToString(now), siteID)
	return err
}

func (r *Repository) RevokeSiteLicense(siteID string) error {
	query := `UPDATE site_licenses SET status = ? WHERE site_id = ?`
	_, err := r.db.Connection.Exec(query, "revoked", siteID)
	return err
}

func (r *Repository) CountActiveSites(orgID string) (int, error) {
	query := `SELECT COUNT(*) FROM site_licenses WHERE org_id = ? AND status = 'active'`
	var count int
	err := r.db.Connection.QueryRow(query, orgID).Scan(&count)
	return count, err
}

func (r *Repository) DeleteSiteLicense(siteID string) error {
	query := `DELETE FROM site_licenses WHERE site_id = ?`
	_, err := r.db.Connection.Exec(query, siteID)
	return err
}
