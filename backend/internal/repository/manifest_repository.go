package repository

import (
	"database/sql"
	"fmt"
	"time"

	"taskmaster-license/internal/models"
)

func (r *Repository) CreateManifest(manifest *models.UsageManifest) error {
	query := `INSERT INTO usage_manifests (id, org_id, period, manifest_data, signature, 
		sent_to_astack, sent_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	sentInt := 0
	if manifest.SentToAStack {
		sentInt = 1
	}

	var sentAtStr sql.NullString
	if manifest.SentAt != nil {
		sentAtStr = sql.NullString{String: timeToString(*manifest.SentAt), Valid: true}
	}

	_, err := r.db.Connection.Exec(query,
		manifest.ID,
		manifest.OrgID,
		manifest.Period,
		string(manifest.ManifestData),
		manifest.Signature,
		sentInt,
		sentAtStr,
		timeToString(manifest.CreatedAt),
	)
	return err
}

func (r *Repository) GetManifest(id string) (*models.UsageManifest, error) {
	query := `SELECT id, org_id, period, manifest_data, signature, sent_to_astack, sent_at, created_at
		FROM usage_manifests WHERE id = ?`

	var manifest models.UsageManifest
	var createdAtStr string
	var sentToAStackInt int
	var sentAt sql.NullString

	err := r.db.Connection.QueryRow(query, id).Scan(
		&manifest.ID,
		&manifest.OrgID,
		&manifest.Period,
		&manifest.ManifestData,
		&manifest.Signature,
		&sentToAStackInt,
		&sentAt,
		&createdAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("manifest not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	manifest.SentToAStack = sentToAStackInt == 1
	manifest.CreatedAt, _ = stringToTime(createdAtStr)

	if sentAt.Valid {
		t, _ := stringToTime(sentAt.String)
		manifest.SentAt = &t
	}

	return &manifest, nil
}

func (r *Repository) ListManifests(orgID, period string) ([]models.UsageManifest, error) {
	query := `SELECT id, org_id, period, manifest_data, signature, sent_to_astack, sent_at, created_at
		FROM usage_manifests WHERE 1=1`
	args := []interface{}{}

	if orgID != "" {
		query += " AND org_id = ?"
		args = append(args, orgID)
	}
	if period != "" {
		query += " AND period = ?"
		args = append(args, period)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Connection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var manifests []models.UsageManifest
	for rows.Next() {
		var manifest models.UsageManifest
		var createdAtStr string
		var sentToAStackInt int
		var sentAt sql.NullString

		err := rows.Scan(
			&manifest.ID,
			&manifest.OrgID,
			&manifest.Period,
			&manifest.ManifestData,
			&manifest.Signature,
			&sentToAStackInt,
			&sentAt,
			&createdAtStr,
		)
		if err != nil {
			return nil, err
		}

		manifest.SentToAStack = sentToAStackInt == 1
		manifest.CreatedAt, _ = stringToTime(createdAtStr)

		if sentAt.Valid {
			t, _ := stringToTime(sentAt.String)
			manifest.SentAt = &t
		}

		manifests = append(manifests, manifest)
	}

	return manifests, nil
}

func (r *Repository) UpdateManifestSent(id string, sentToAStack bool, sentAt *time.Time) error {
	sentInt := 0
	if sentToAStack {
		sentInt = 1
	}

	var sentAtStr sql.NullString
	if sentAt != nil {
		sentAtStr = sql.NullString{String: timeToString(*sentAt), Valid: true}
	}

	query := `UPDATE usage_manifests SET sent_to_astack = ?, sent_at = ? WHERE id = ?`
	_, err := r.db.Connection.Exec(query, sentInt, sentAtStr, id)
	return err
}
