package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// Salt size for PBKDF2 key derivation
	saltSize = 32
	// Nonce size for AES-GCM
	nonceSize = 12
	// PBKDF2 iterations for key derivation (CPU-intensive to prevent brute force)
	pbkdf2Iterations = 100000
)

// EncryptPrivateKey encrypts a private key using AES-256-GCM with PBKDF2 key derivation
// Returns base64-encoded encrypted data
func EncryptPrivateKey(keyBytes []byte, password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("password cannot be empty")
	}

	// Generate a random salt
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key from password using PBKDF2
	derivedKey := pbkdf2.Key([]byte(password), salt, pbkdf2Iterations, 32, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate a random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the data
	ciphertext := gcm.Seal(nil, nonce, keyBytes, nil)

	// Combine salt, nonce, and ciphertext
	// Format: base64(salt + nonce + ciphertext)
	combined := make([]byte, len(salt)+len(nonce)+len(ciphertext))
	copy(combined[0:], salt)
	copy(combined[len(salt):], nonce)
	copy(combined[len(salt)+len(nonce):], ciphertext)

	// Encode to base64 for storage
	encrypted := base64.StdEncoding.EncodeToString(combined)
	return encrypted, nil
}

// DecryptPrivateKey decrypts a private key using AES-256-GCM
// Expects base64-encoded encrypted data
func DecryptPrivateKey(encrypted string, password string) ([]byte, error) {
	if len(password) == 0 {
		return nil, errors.New("password cannot be empty")
	}

	// Decode from base64
	combined, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted data: %w", err)
	}

	// Extract salt, nonce, and ciphertext
	if len(combined) < saltSize+nonceSize {
		return nil, errors.New("encrypted data too short")
	}

	salt := combined[0:saltSize]
	nonce := combined[saltSize : saltSize+nonceSize]
	ciphertext := combined[saltSize+nonceSize:]

	// Derive key from password using same PBKDF2 parameters
	derivedKey := pbkdf2.Key([]byte(password), salt, pbkdf2Iterations, 32, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt the data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// ValidatePassword checks if the password meets minimum requirements
func ValidatePassword(password string) error {
	if len(password) < 16 {
		return errors.New("password must be at least 16 characters")
	}
	return nil
}

