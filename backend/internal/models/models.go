package models

import (
	"encoding/json"
	"time"
)

// CML represents Customer Master License
type CML struct {
	ID           string          `json:"id"`
	OrgID        string          `json:"org_id"`
	MaxSites     int             `json:"max_sites"`
	Validity     time.Time       `json:"validity"`
	FeaturePacks []string        `json:"feature_packs"`
	DevKeyPublic string          `json:"dev_key_public"`
	ProdKeyPublic string          `json:"prod_key_public"`
	CMLData      json.RawMessage `json:"cml_data"`
	Signature    string          `json:"signature"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// SiteLicense represents a site license
type SiteLicense struct {
	ID           string          `json:"id"`
	SiteID       string          `json:"site_id"`
	OrgID        string          `json:"org_id"`
	Fingerprint  json.RawMessage `json:"fingerprint"`
	LicenseData  json.RawMessage `json:"license_data"`
	Signature    string          `json:"signature"`
	IssuedAt     time.Time       `json:"issued_at"`
	LastSeen     *time.Time      `json:"last_seen,omitempty"`
	Status       string          `json:"status"`
	CreatedAt    time.Time       `json:"created_at"`
}

// UsageLedgerEntry represents a ledger entry
type UsageLedgerEntry struct {
	ID        string          `json:"id"`
	OrgID     string          `json:"org_id"`
	EntryType string          `json:"entry_type"`
	SiteID    string          `json:"site_id,omitempty"`
	Data      json.RawMessage `json:"data"`
	Signature string          `json:"signature,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

// UsageStats represents aggregated statistics
type UsageStats struct {
	ID                string          `json:"id"`
	OrgID             string          `json:"org_id"`
	Period            string          `json:"period"`
	UserStats         json.RawMessage `json:"user_stats"`
	SiteStats         json.RawMessage `json:"site_stats"`
	TotalActiveSites  int             `json:"total_active_sites"`
	CreatedAt         time.Time       `json:"created_at"`
}

// UsageManifest represents usage manifest
type UsageManifest struct {
	ID           string          `json:"id"`
	OrgID        string          `json:"org_id"`
	Period       string          `json:"period"`
	ManifestData json.RawMessage `json:"manifest_data"`
	Signature    string          `json:"signature"`
	SentToAStack bool            `json:"sent_to_astack"`
	SentAt       *time.Time      `json:"sent_at,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
}

// OrgKey represents organization keys
type OrgKey struct {
	ID                string    `json:"id"`
	OrgID             string    `json:"org_id"`
	KeyType           string    `json:"key_type"`
	PrivateKeyEncrypted string    `json:"private_key_encrypted"`
	PublicKey         string    `json:"public_key"`
	CreatedAt         time.Time `json:"created_at"`
}

// CMLData represents the actual CML JSON data
type CMLData struct {
	Type          string   `json:"type"`
	OrgID         string   `json:"org_id"`
	MaxSites      int      `json:"max_sites"`
	Validity      string   `json:"validity"`
	FeaturePacks  []string `json:"feature_packs"`
	KeyType       string   `json:"key_type"`
	IssuedBy      string   `json:"issued_by"`
	IssuerPublicKey string `json:"issuer_public_key"`
	IssuedAt      string   `json:"issued_at"`
}

// SiteLicenseData represents the actual site license JSON data
type SiteLicenseData struct {
	Type         string                 `json:"type"`
	SiteID       string                 `json:"site_id"`
	ParentCML    string                 `json:"parent_cml"`
	ParentCMLSig string                 `json:"parent_cml_sig"`
	Fingerprint  map[string]interface{} `json:"fingerprint"`
	IssuedAt     string                 `json:"issued_at"`
	ExpiresAt    string                 `json:"expires_at"`
	Features     []string               `json:"features"`
}

// Fingerprint represents site fingerprint data
type Fingerprint struct {
	Address       string `json:"address,omitempty"`
	DNSSuffix     string `json:"dns_suffix,omitempty"`
	DeploymentTag string `json:"deployment_tag,omitempty"`
}
