package models

import "time"

// License server models - merged from license-server

// SiteKey represents a site key entry
type SiteKey struct {
	ID            string     `json:"id"`
	SiteID        string     `json:"site_id"`
	EnterpriseID  string     `json:"enterprise_id"`
	KeyType       string     `json:"key_type"`
	KeyValue      string     `json:"key_value"`
	IssuedAt      time.Time  `json:"issued_at"`
	ExpiresAt     time.Time  `json:"expires_at"`
	Status        string     `json:"status"`
	LastValidated *time.Time `json:"last_validated,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

// Enterprise represents an enterprise entry
type Enterprise struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	OrgID         string    `json:"org_id"`
	EnterpriseKey string    `json:"enterprise_key"`
	CreatedAt     time.Time `json:"created_at"`
}

// KeyRefreshLog represents a key refresh audit entry
type KeyRefreshLog struct {
	ID          string    `json:"id"`
	SiteID      string    `json:"site_id"`
	OldKey      string    `json:"old_key"`
	NewKey      string    `json:"new_key"`
	RefreshedAt time.Time `json:"refreshed_at"`
	Reason      string    `json:"reason,omitempty"`
}

// QuarterlyStats represents quarterly aggregated statistics
type QuarterlyStats struct {
	ID                  string                   `json:"id"`
	Period              string                   `json:"period"`
	ProductionSites     int                      `json:"production_sites"`
	DevSites            int                      `json:"dev_sites"`
	UserCounts          map[string]interface{}   `json:"user_counts"`
	EnterpriseBreakdown []map[string]interface{} `json:"enterprise_breakdown"`
	CreatedAt           time.Time                `json:"created_at"`
}

// ValidationCache represents a cached validation token
type ValidationCache struct {
	ID        string    `json:"id"`
	SiteID    string    `json:"site_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// Alert represents an alert entry
type Alert struct {
	ID           string    `json:"id"`
	SiteID       string    `json:"site_id"`
	AlertType    string    `json:"alert_type"`
	Message      string    `json:"message"`
	Timestamp    time.Time `json:"timestamp"`
	SentToAStack bool      `json:"sent_to_astack"`
	CreatedAt    time.Time `json:"created_at"`
}

// ValidationResponse represents the response from key validation
type ValidationResponse struct {
	Valid     bool   `json:"valid"`
	Token     string `json:"token,omitempty"`
	ExpiresIn int    `json:"expires_in,omitempty"`
	Message   string `json:"message"`
}

