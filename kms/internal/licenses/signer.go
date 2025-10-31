package licenses

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/atprof/license-server/kms/pkg/errors"
)

// SignLicense signs license content using HMAC-SHA256 with the master key
// Returns base64-encoded signature
func SignLicense(content []byte, masterKey []byte) (string, error) {
	if len(masterKey) != 32 {
		return "", errors.ErrInvalidKeyMaterial
	}

	mac := hmac.New(sha256.New, masterKey)
	mac.Write(content)
	signature := mac.Sum(nil)

	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifyLicenseSignature verifies the HMAC-SHA256 signature of license content
// Uses constant-time comparison to prevent timing attacks
func VerifyLicenseSignature(content []byte, signatureBase64 string, masterKey []byte) (bool, error) {
	if len(masterKey) != 32 {
		return false, errors.ErrInvalidKeyMaterial
	}

	// Decode the signature
	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return false, errors.ErrInvalidSignature
	}

	// Compute expected signature
	mac := hmac.New(sha256.New, masterKey)
	mac.Write(content)
	expectedSignature := mac.Sum(nil)

	// Constant-time comparison
	if len(signature) != len(expectedSignature) {
		return false, errors.ErrInvalidSignature
	}

	// Use constant-time comparison
	result := 0
	for i := 0; i < len(signature); i++ {
		result |= int(signature[i] ^ expectedSignature[i])
	}

	return result == 0, nil
}

