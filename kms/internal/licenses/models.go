package licenses

import "time"

// LicenseFile represents a license file structure
type LicenseFile struct {
	LicenseID   string            `json:"license_id"`
	LicenseType string            `json:"license_type"`
	KeyID       string            `json:"key_id"`
	KeyType     string            `json:"key_type"`
	PublicKey   string            `json:"public_key,omitempty"` // Base64 encoded, only for asymmetric keys
	IssuedAt    time.Time         `json:"issued_at"`
	ExpiresAt   time.Time         `json:"expires_at"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Signature   string            `json:"signature"` // Base64 encoded HMAC-SHA256 signature
}

// GenerateLicenseRequest represents a request to generate a license file
type GenerateLicenseRequest struct {
	KeyID       string            `json:"key_id" binding:"required"`
	LicenseType string            `json:"license_type" binding:"required"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// GenerateLicenseResponse represents a response from generating a license file
type GenerateLicenseResponse struct {
	LicenseFile string `json:"license_file"` // Base64 encoded license file content
	Filename    string `json:"filename"`     // Suggested filename (e.g., "enterprise.lic")
	LicenseID   string `json:"license_id"`
}

// ValidateLicenseRequest represents a request to validate a license file
// Can be either multipart file upload or JSON with base64 content
type ValidateLicenseRequest struct {
	LicenseContent string `json:"license_content,omitempty"` // Base64 encoded license file (for JSON body)
}

// ValidationResult represents the result of license validation
type ValidationResult struct {
	Valid      bool              `json:"valid"`
	LicenseID  string            `json:"license_id,omitempty"`
	LicenseType string           `json:"license_type,omitempty"`
	KeyID      string            `json:"key_id,omitempty"`
	ExpiresAt  time.Time         `json:"expires_at,omitempty"`
	Expired    bool              `json:"expired"`
	Revoked    bool              `json:"revoked"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	Error      string            `json:"error,omitempty"` // Error message if validation failed
}

// ValidateLicenseResponse represents a response from validating a license file
type ValidateLicenseResponse struct {
	Valid       bool              `json:"valid"`
	LicenseID   string            `json:"license_id,omitempty"`
	LicenseType string            `json:"license_type,omitempty"`
	KeyID       string            `json:"key_id,omitempty"`
	ExpiresAt   time.Time         `json:"expires_at,omitempty"`
	Expired     bool              `json:"expired"`
	Revoked     bool              `json:"revoked"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Error       string            `json:"error,omitempty"`
}

