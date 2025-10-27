package repository

import (
	"database/sql"
	"time"

	"taskmaster-license/internal/database"
)


type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// Helper methods for time conversion
func timeToString(t time.Time) string {
	return t.Format(time.RFC3339)
}

func stringToTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func nullTimeToString(t sql.NullTime) *string {
	if !t.Valid {
		return nil
	}
	s := t.Time.Format(time.RFC3339)
	return &s
}

func stringToNullTime(s *string) (sql.NullTime, error) {
	if s == nil {
		return sql.NullTime{Valid: false}, nil
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return sql.NullTime{}, err
	}
	return sql.NullTime{Time: t, Valid: true}, nil
}
