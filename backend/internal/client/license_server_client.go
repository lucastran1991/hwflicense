package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type LicenseServerClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewLicenseServerClient() *LicenseServerClient {
	baseURL := getEnv("LICENSE_SERVER_URL", "http://localhost:8081")
	
	return &LicenseServerClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// SiteKeyResponse represents response from license server
type SiteKeyResponse struct {
	ID           string    `json:"id"`
	SiteID       string    `json:"site_id"`
	EnterpriseID string    `json:"enterprise_id"`
	KeyType      string    `json:"key_type"`
	KeyValue     string    `json:"key_value"`
	IssuedAt     time.Time `json:"issued_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	Status       string    `json:"status"`
}

// ValidationResponse represents validation response
type ValidationResponse struct {
	Valid     bool   `json:"valid"`
	Token     string `json:"token,omitempty"`
	ExpiresIn int    `json:"expires_in,omitempty"`
	Message   string `json:"message"`
}

// CreateSiteKeyRequest represents request to create site key
type CreateSiteKeyRequest struct {
	SiteID       string `json:"site_id"`
	EnterpriseID string `json:"enterprise_id"`
	Mode         string `json:"mode"` // "production" or "dev"
	OrgID        string `json:"org_id"`
}

// RefreshKeyRequest represents request to refresh key
type RefreshKeyRequest struct {
	SiteID string `json:"site_id"`
	OldKey string `json:"old_key"`
}

// ValidateKeyRequest represents request to validate key
type ValidateKeyRequest struct {
	SiteID string `json:"site_id"`
	Key    string `json:"key"`
}

// QuarterlyStats represents quarterly stats  
type QuarterlyStats struct {
	Period              string                 `json:"period"`
	ProductionSites     int                    `json:"production_sites"`
	DevSites            int                    `json:"dev_sites"`
	UserCounts          map[string]interface{} `json:"user_counts"`
	EnterpriseBreakdown []map[string]interface{} `json:"enterprise_breakdown"`
}

// AlertRequest represents alert to license server
type AlertRequest struct {
	SiteID    string    `json:"site_id"`
	AlertType string    `json:"alert_type"` // "key_expired" or "key_invalid"
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// CreateSiteKey creates a site key on license server (API 1)
func (c *LicenseServerClient) CreateSiteKey(siteID, enterpriseID, mode, orgID string) (*SiteKeyResponse, error) {
	reqBody := CreateSiteKeyRequest{
		SiteID:       siteID,
		EnterpriseID: enterpriseID,
		Mode:         mode,
		OrgID:        orgID,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/api/v1/sites/create", c.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call license server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("license server returned status %d: %s", resp.StatusCode, body)
	}

	var response SiteKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// RefreshKey refreshes a site key (API 4)
func (c *LicenseServerClient) RefreshKey(siteID, oldKey string) (*SiteKeyResponse, error) {
	reqBody := RefreshKeyRequest{
		SiteID: siteID,
		OldKey: oldKey,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/api/v1/keys/refresh", c.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call license server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("license server returned status %d: %s", resp.StatusCode, body)
	}

	var response SiteKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// ValidateKey validates a site key (API 6)
func (c *LicenseServerClient) ValidateKey(siteID, key string) (*ValidationResponse, error) {
	reqBody := ValidateKeyRequest{
		SiteID: siteID,
		Key:    key,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/api/v1/keys/validate", c.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call license server: %w", err)
	}
	defer resp.Body.Close()

	var response ValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// SendStats sends quarterly stats to license server (API 5)
func (c *LicenseServerClient) SendStats(stats *QuarterlyStats) error {
	jsonData, err := json.Marshal(stats)
	if err != nil {
		return fmt.Errorf("failed to marshal stats: %w", err)
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/api/v1/stats/aggregate", c.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to call license server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("license server returned status %d: %s", resp.StatusCode, body)
	}

	return nil
}

// SendAlert sends an alert to license server (API 7)
func (c *LicenseServerClient) SendAlert(alert *AlertRequest) error {
	jsonData, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("failed to marshal alert: %w", err)
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/api/v1/alerts", c.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to call license server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("license server returned status %d: %s", resp.StatusCode, body)
	}

	return nil
}

