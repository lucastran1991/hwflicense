package licenses

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/atprof/license-server/kms/internal/storage"
)

// GenerateLicense generates a license file for a given key
// Returns the LicenseFile struct and raw JSON bytes
func GenerateLicense(key *storage.Key, licenseType string, metadata map[string]string, masterKey []byte) (*LicenseFile, []byte, error) {
	// Validate key is active
	if !key.IsValid() {
		if key.IsExpired() {
			return nil, nil, fmt.Errorf("cannot generate license for expired key")
		}
		if key.IsRevoked() {
			return nil, nil, fmt.Errorf("cannot generate license for revoked key")
		}
		return nil, nil, fmt.Errorf("key is not valid")
	}

	// Create license structure
	license := &LicenseFile{
		LicenseID:   uuid.New().String(),
		LicenseType: licenseType,
		KeyID:       key.ID,
		KeyType:     string(key.KeyType),
		IssuedAt:    time.Now().UTC(),
		ExpiresAt:   key.ExpiresAt,
		Metadata:    metadata,
	}

	// Add public key if asymmetric
	if key.KeyType == storage.KeyTypeAsymmetric && len(key.PublicKey) > 0 {
		license.PublicKey = base64.StdEncoding.EncodeToString(key.PublicKey)
	}

	// Serialize to JSON without signature first
	tempLicense := *license
	tempLicense.Signature = ""
	jsonWithoutSig, err := json.Marshal(tempLicense)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal license: %w", err)
	}

	// Sign the JSON content (without signature)
	signature, err := SignLicense(jsonWithoutSig, masterKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to sign license: %w", err)
	}

	// Add signature to license
	license.Signature = signature

	// Serialize final JSON with signature
	finalJSON, err := json.Marshal(license)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal final license: %w", err)
	}

	return license, finalJSON, nil
}
