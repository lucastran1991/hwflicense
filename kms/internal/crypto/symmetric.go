package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"fmt"

	"github.com/atprof/license-server/kms/pkg/errors"
)

const (
	// SymmetricKeySize is the size of symmetric keys in bytes (256 bits)
	SymmetricKeySize = 32
)

// GenerateSymmetricKey generates a new random symmetric key
// Returns a 32-byte (256-bit) key suitable for AES-256
func GenerateSymmetricKey() ([]byte, error) {
	key := make([]byte, SymmetricKeySize)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate symmetric key: %w", err)
	}
	return key, nil
}

// CompareKeys compares two keys in constant time to prevent timing attacks
// Returns true if keys are equal, false otherwise
func CompareKeys(key1, key2 []byte) bool {
	if len(key1) != len(key2) {
		return false
	}
	return subtle.ConstantTimeCompare(key1, key2) == 1
}

// ValidateSymmetricKey validates a provided symmetric key against a stored encrypted key
// It decrypts the stored key and compares it with the provided key in constant time
func ValidateSymmetricKey(masterKey []byte, storedEncryptedKey, providedKey []byte) (bool, error) {
	// Decrypt the stored key
	decryptedKey, err := DecryptKey(masterKey, storedEncryptedKey)
	if err != nil {
		return false, err
	}
	defer func() {
		// Zero out the decrypted key after use
		for i := range decryptedKey {
			decryptedKey[i] = 0
		}
	}()

	// Compare keys in constant time
	if len(decryptedKey) != len(providedKey) {
		return false, errors.ErrInvalidKeyMaterial
	}

	valid := CompareKeys(decryptedKey, providedKey)
	
	// Zero out provided key if it matches (defensive)
	if valid {
		for i := range providedKey {
			providedKey[i] = 0
		}
	}
	
	return valid, nil
}

