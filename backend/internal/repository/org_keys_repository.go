package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"taskmaster-license/internal/models"
)

// CreateOrgKey creates a new organization key in the database
func (r *Repository) CreateOrgKey(orgKey *models.OrgKey) error {
	query := `INSERT INTO org_keys (id, org_id, key_type, private_key_encrypted, public_key, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err := r.db.Connection.Exec(query,
		orgKey.ID,
		orgKey.OrgID,
		orgKey.KeyType,
		orgKey.PrivateKeyEncrypted,
		orgKey.PublicKey,
		timeToString(orgKey.CreatedAt),
	)
	return err
}

// GetOrgKey retrieves an organization key by org_id and key_type
func (r *Repository) GetOrgKey(orgID, keyType string) (*models.OrgKey, error) {
	query := `SELECT id, org_id, key_type, private_key_encrypted, public_key, created_at
		FROM org_keys WHERE org_id = ? AND key_type = ?`

	var orgKey models.OrgKey
	var createdAtStr string

	err := r.db.Connection.QueryRow(query, orgID, keyType).Scan(
		&orgKey.ID,
		&orgKey.OrgID,
		&orgKey.KeyType,
		&orgKey.PrivateKeyEncrypted,
		&orgKey.PublicKey,
		&createdAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("org key not found: org_id=%s, key_type=%s", orgID, keyType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get org key: %w", err)
	}

	createdAt, err := stringToTime(createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}
	orgKey.CreatedAt = createdAt

	return &orgKey, nil
}

// GetOrgKeyByID retrieves an organization key by its ID
func (r *Repository) GetOrgKeyByID(id string) (*models.OrgKey, error) {
	query := `SELECT id, org_id, key_type, private_key_encrypted, public_key, created_at
		FROM org_keys WHERE id = ?`

	var orgKey models.OrgKey
	var createdAtStr string

	err := r.db.Connection.QueryRow(query, id).Scan(
		&orgKey.ID,
		&orgKey.OrgID,
		&orgKey.KeyType,
		&orgKey.PrivateKeyEncrypted,
		&orgKey.PublicKey,
		&createdAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("org key not found: id=%s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get org key: %w", err)
	}

	createdAt, err := stringToTime(createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}
	orgKey.CreatedAt = createdAt

	return &orgKey, nil
}

// ListOrgKeys retrieves all organization keys for a specific org_id
func (r *Repository) ListOrgKeys(orgID string) ([]models.OrgKey, error) {
	query := `SELECT id, org_id, key_type, private_key_encrypted, public_key, created_at
		FROM org_keys WHERE org_id = ? ORDER BY created_at DESC`

	rows, err := r.db.Connection.Query(query, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to query org keys: %w", err)
	}
	defer rows.Close()

	var orgKeys []models.OrgKey
	for rows.Next() {
		var orgKey models.OrgKey
		var createdAtStr string

		if err := rows.Scan(
			&orgKey.ID,
			&orgKey.OrgID,
			&orgKey.KeyType,
			&orgKey.PrivateKeyEncrypted,
			&orgKey.PublicKey,
			&createdAtStr,
		); err != nil {
			return nil, fmt.Errorf("failed to scan org key: %w", err)
		}

		createdAt, err := stringToTime(createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}
		orgKey.CreatedAt = createdAt

		orgKeys = append(orgKeys, orgKey)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating org keys: %w", err)
	}

	return orgKeys, nil
}

// DeleteOrgKey deletes an organization key by ID
func (r *Repository) DeleteOrgKey(id string) error {
	query := `DELETE FROM org_keys WHERE id = ?`
	
	result, err := r.db.Connection.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete org key: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("org key not found: id=%s", id)
	}

	return nil
}

// CreateOrgKeyWithID creates an org key with a specific ID (useful for testing/imports)
func (r *Repository) CreateOrgKeyWithID(orgKey *models.OrgKey) error {
	if orgKey.ID == "" {
		orgKey.ID = uuid.New().String()
	}
	return r.CreateOrgKey(orgKey)
}

