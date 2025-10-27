package repository

import (
	"database/sql"

	"taskmaster-license/internal/models"
)

func (r *Repository) CreateLedgerEntry(entry *models.UsageLedgerEntry) error {
	query := `INSERT INTO usage_ledger (id, org_id, entry_type, site_id, data, signature, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Connection.Exec(query,
		entry.ID,
		entry.OrgID,
		entry.EntryType,
		entry.SiteID,
		string(entry.Data),
		entry.Signature,
		timeToString(entry.CreatedAt),
	)
	return err
}

func (r *Repository) GetLedgerEntries(orgID string, limit, offset int) ([]models.UsageLedgerEntry, int, error) {
	// Get entries
	query := `SELECT id, org_id, entry_type, site_id, data, signature, created_at 
		FROM usage_ledger WHERE org_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Connection.Query(query, orgID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var entries []models.UsageLedgerEntry
	for rows.Next() {
		var entry models.UsageLedgerEntry
		var createdAtStr string
		var siteID sql.NullString
		var dataStr sql.NullString
		var signature sql.NullString

		err := rows.Scan(
			&entry.ID,
			&entry.OrgID,
			&entry.EntryType,
			&siteID,
			&dataStr,
			&signature,
			&createdAtStr,
		)
		if err != nil {
			return nil, 0, err
		}

		entry.CreatedAt, _ = stringToTime(createdAtStr)
		if siteID.Valid {
			entry.SiteID = siteID.String
		}
		if dataStr.Valid {
			entry.Data = []byte(dataStr.String)
		}
		if signature.Valid {
			entry.Signature = signature.String
		}

		entries = append(entries, entry)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM usage_ledger WHERE org_id = ?`
	var total int
	err = r.db.Connection.QueryRow(countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

func (r *Repository) GetActiveSiteCount(orgID string) (int, error) {
	query := `SELECT COUNT(*) FROM site_licenses WHERE org_id = ? AND status = 'active'`
	var count int
	err := r.db.Connection.QueryRow(query, orgID).Scan(&count)
	return count, err
}
