package repository

import (
	"database/sql"
	"fmt"
	"time"

	"taskmaster-license/internal/models"
)

func (r *Repository) CreateCML(cml *models.CML) error {
	query := `INSERT INTO cml (id, org_id, max_sites, validity, feature_packs, dev_key_public, 
		prod_key_public, cml_data, signature, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Connection.Exec(query,
		cml.ID,
		cml.OrgID,
		cml.MaxSites,
		timeToString(cml.Validity),
		toJSON(cml.FeaturePacks),
		cml.DevKeyPublic,
		cml.ProdKeyPublic,
		string(cml.CMLData),
		cml.Signature,
		timeToString(cml.CreatedAt),
		timeToString(cml.UpdatedAt),
	)
	return err
}

func (r *Repository) GetCML(orgID string) (*models.CML, error) {
	query := `SELECT id, org_id, max_sites, validity, feature_packs, dev_key_public, 
		prod_key_public, cml_data, signature, created_at, updated_at
		FROM cml WHERE org_id = ?`

	var cml models.CML
	var validityStr, createdAtStr, updatedAtStr string
	var featurePacksJSON string

	err := r.db.Connection.QueryRow(query, orgID).Scan(
		&cml.ID,
		&cml.OrgID,
		&cml.MaxSites,
		&validityStr,
		&featurePacksJSON,
		&cml.DevKeyPublic,
		&cml.ProdKeyPublic,
		&cml.CMLData,
		&cml.Signature,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("CML not found for org_id: %s", orgID)
	}
	if err != nil {
		return nil, err
	}

	// Parse timestamps
	validity, err := stringToTime(validityStr)
	if err != nil {
		return nil, err
	}
	cml.Validity = validity

	createdAt, err := stringToTime(createdAtStr)
	if err != nil {
		return nil, err
	}
	cml.CreatedAt = createdAt

	updatedAt, err := stringToTime(updatedAtStr)
	if err != nil {
		return nil, err
	}
	cml.UpdatedAt = updatedAt

	// Parse feature packs
	cml.FeaturePacks = fromJSONArray(featurePacksJSON)

	return &cml, nil
}

func (r *Repository) UpdateCML(cml *models.CML) error {
	query := `UPDATE cml SET max_sites = ?, validity = ?, feature_packs = ?, 
		dev_key_public = ?, prod_key_public = ?, cml_data = ?, signature = ?, updated_at = ?
		WHERE org_id = ?`

	now := time.Now()
	_, err := r.db.Connection.Exec(query,
		cml.MaxSites,
		timeToString(cml.Validity),
		toJSON(cml.FeaturePacks),
		cml.DevKeyPublic,
		cml.ProdKeyPublic,
		string(cml.CMLData),
		cml.Signature,
		timeToString(now),
		cml.OrgID,
	)
	return err
}

func (r *Repository) GetAllCML() ([]models.CML, error) {
	query := `SELECT id, org_id, max_sites, validity, feature_packs, dev_key_public, 
		prod_key_public, cml_data, signature, created_at, updated_at FROM cml`

	rows, err := r.db.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cmls []models.CML
	for rows.Next() {
		var cml models.CML
		var validityStr, createdAtStr, updatedAtStr string
		var featurePacksJSON string

		err := rows.Scan(
			&cml.ID,
			&cml.OrgID,
			&cml.MaxSites,
			&validityStr,
			&featurePacksJSON,
			&cml.DevKeyPublic,
			&cml.ProdKeyPublic,
			&cml.CMLData,
			&cml.Signature,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, err
		}

		// Parse timestamps
		cml.Validity, _ = stringToTime(validityStr)
		cml.CreatedAt, _ = stringToTime(createdAtStr)
		cml.UpdatedAt, _ = stringToTime(updatedAtStr)
		cml.FeaturePacks = fromJSONArray(featurePacksJSON)

		cmls = append(cmls, cml)
	}

	return cmls, nil
}
