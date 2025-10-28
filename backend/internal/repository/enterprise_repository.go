package repository

import (
	"taskmaster-license/internal/models"
	"time"
)

// CreateEnterprise creates a new enterprise
func (r *Repository) CreateEnterprise(enterprise *models.Enterprise) error {
	query := `INSERT INTO enterprises (id, name, org_id, enterprise_key, created_at)
			  VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.Connection.Exec(query, enterprise.ID, enterprise.Name, enterprise.OrgID,
		enterprise.EnterpriseKey, time.Now().Format("2006-01-02 15:04:05"))
	return err
}

// GetEnterprise retrieves enterprise by org_id
func (r *Repository) GetEnterprise(orgID string) (*models.Enterprise, error) {
	query := `SELECT id, name, org_id, enterprise_key, created_at FROM enterprises WHERE org_id = ?`
	row := r.db.Connection.QueryRow(query, orgID)

	var e models.Enterprise
	var createdAt string
	err := row.Scan(&e.ID, &e.Name, &e.OrgID, &e.EnterpriseKey, &createdAt)
	if err != nil {
		return nil, err
	}

	e.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	return &e, nil
}

// ListEnterprises lists all enterprises
func (r *Repository) ListEnterprises(orgID string) ([]*models.Enterprise, error) {
	query := `SELECT id, name, org_id, enterprise_key, created_at FROM enterprises 
			  WHERE org_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Connection.Query(query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enterprises []*models.Enterprise
	for rows.Next() {
		var e models.Enterprise
		var createdAt string
		if err := rows.Scan(&e.ID, &e.Name, &e.OrgID, &e.EnterpriseKey, &createdAt); err != nil {
			continue
		}
		e.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		enterprises = append(enterprises, &e)
	}

	return enterprises, rows.Err()
}

// UpdateEnterprise updates an enterprise
func (r *Repository) UpdateEnterprise(enterprise *models.Enterprise) error {
	query := `UPDATE enterprises SET name = ?, enterprise_key = ? WHERE id = ?`
	_, err := r.db.Connection.Exec(query, enterprise.Name, enterprise.EnterpriseKey, enterprise.ID)
	return err
}

