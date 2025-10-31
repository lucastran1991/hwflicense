package licenses

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/atprof/license-server/kms/internal/storage"
	"github.com/atprof/license-server/kms/pkg/errors"
)

// ValidateLicense validates a license file
// Returns validation result with license information
func ValidateLicense(fileContent []byte, store *storage.BoltStore, masterKey []byte) (*ValidationResult, error) {
	// Parse license file
	var license LicenseFile
	if err := json.Unmarshal(fileContent, &license); err != nil {
		return &ValidationResult{
			Valid: false,
			Error: fmt.Sprintf("failed to parse license file: %v", err),
		}, nil
	}

	// Extract signature for verification
	signature := license.Signature
	if signature == "" {
		return &ValidationResult{
			Valid: false,
			Error: "license file missing signature",
		}, nil
	}

	// Remove signature for verification
	tempLicense := license
	tempLicense.Signature = ""
	jsonWithoutSig, err := json.Marshal(tempLicense)
	if err != nil {
		return &ValidationResult{
			Valid: false,
			Error: fmt.Sprintf("failed to marshal license for verification: %v", err),
		}, nil
	}

	// Verify signature
	validSig, err := VerifyLicenseSignature(jsonWithoutSig, signature, masterKey)
	if err != nil {
		return &ValidationResult{
			Valid: false,
			Error: fmt.Sprintf("signature verification failed: %v", err),
		}, nil
	}

	if !validSig {
		return &ValidationResult{
			Valid: false,
			Error: "invalid license signature",
		}, nil
	}

	// Check expiry
	expired := time.Now().After(license.ExpiresAt)
	if expired {
		return &ValidationResult{
			Valid:    false,
			Expired:  true,
			LicenseID: license.LicenseID,
			KeyID:    license.KeyID,
		}, nil
	}

	// Verify key exists in database and is not revoked
	key, err := store.GetKey(license.KeyID)
	if err != nil {
		if err == errors.ErrKeyNotFound {
			return &ValidationResult{
				Valid:    false,
				Error:    "key not found in database",
				LicenseID: license.LicenseID,
				KeyID:    license.KeyID,
			}, nil
		}
		return &ValidationResult{
			Valid:    false,
			Error:    fmt.Sprintf("failed to retrieve key: %v", err),
			LicenseID: license.LicenseID,
			KeyID:    license.KeyID,
		}, nil
	}

	// Check if key is revoked
	if key.IsRevoked() {
		return &ValidationResult{
			Valid:    false,
			Revoked:  true,
			LicenseID: license.LicenseID,
			KeyID:    license.KeyID,
		}, nil
	}

	// License is valid
	return &ValidationResult{
		Valid:      true,
		LicenseID:  license.LicenseID,
		LicenseType: license.LicenseType,
		KeyID:      license.KeyID,
		ExpiresAt:  license.ExpiresAt,
		Expired:    false,
		Revoked:    false,
		Metadata:   license.Metadata,
	}, nil
}
