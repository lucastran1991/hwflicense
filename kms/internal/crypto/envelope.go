package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"github.com/atprof/license-server/kms/pkg/errors"
)

// EncryptKey encrypts plaintext using AES-256-GCM with envelope encryption
// The masterKey must be 32 bytes (256 bits)
// Returns ciphertext with nonce prepended
func EncryptKey(masterKey, plaintext []byte) ([]byte, error) {
	if len(masterKey) != 32 {
		return nil, fmt.Errorf("master key must be 32 bytes")
	}

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrEncryptionFailed, err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrEncryptionFailed, err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("%w: failed to generate nonce: %v", errors.ErrEncryptionFailed, err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptKey decrypts ciphertext using AES-256-GCM
// The masterKey must be 32 bytes (256 bits)
// The ciphertext is expected to have nonce prepended
// The plaintext is zeroed out after decryption
func DecryptKey(masterKey, ciphertext []byte) ([]byte, error) {
	if len(masterKey) != 32 {
		return nil, fmt.Errorf("master key must be 32 bytes")
	}

	if len(ciphertext) < 12 { // Minimum size: nonce (12 bytes) + some ciphertext
		return nil, fmt.Errorf("%w: ciphertext too short", errors.ErrDecryptionFailed)
	}

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrDecryptionFailed, err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrDecryptionFailed, err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("%w: ciphertext too short", errors.ErrDecryptionFailed)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrDecryptionFailed, err)
	}

	return plaintext, nil
}

