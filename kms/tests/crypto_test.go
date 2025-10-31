package tests

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"

	"github.com/atprof/license-server/kms/internal/crypto"
)

// TestGenerateSymmetricKey tests symmetric key generation
func TestGenerateSymmetricKey(t *testing.T) {
	key, err := crypto.GenerateSymmetricKey()
	if err != nil {
		t.Fatalf("Failed to generate symmetric key: %v", err)
	}

	if len(key) != crypto.SymmetricKeySize {
		t.Fatalf("Expected key size %d, got %d", crypto.SymmetricKeySize, len(key))
	}

	// Generate another key and ensure it's different
	key2, err := crypto.GenerateSymmetricKey()
	if err != nil {
		t.Fatalf("Failed to generate second symmetric key: %v", err)
	}

	// Keys should be different
	if crypto.CompareKeys(key, key2) {
		t.Error("Generated keys should be different")
	}
}

// TestEnvelopeEncryption tests envelope encryption and decryption
func TestEnvelopeEncryption(t *testing.T) {
	// Generate a master key (32 bytes)
	masterKey := make([]byte, 32)
	rand.Read(masterKey)

	// Generate plaintext
	plaintext := []byte("test data for encryption")
	if len(plaintext) < 1 {
		t.Fatal("Plaintext too short")
	}

	// Encrypt
	ciphertext, err := crypto.EncryptKey(masterKey, plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	if len(ciphertext) <= len(plaintext) {
		t.Error("Ciphertext should be longer than plaintext (includes nonce)")
	}

	// Decrypt
	decrypted, err := crypto.DecryptKey(masterKey, ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	// Verify
	if !crypto.CompareKeys(plaintext, decrypted) {
		t.Error("Decrypted text does not match plaintext")
	}
}

// TestEnvelopeEncryptionDifferentKeys tests that decryption fails with wrong key
func TestEnvelopeEncryptionDifferentKeys(t *testing.T) {
	// Generate master keys
	masterKey1 := make([]byte, 32)
	rand.Read(masterKey1)

	masterKey2 := make([]byte, 32)
	rand.Read(masterKey2)

	plaintext := []byte("test data")

	// Encrypt with key1
	ciphertext, err := crypto.EncryptKey(masterKey1, plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	// Try to decrypt with key2 - should fail
	_, err = crypto.DecryptKey(masterKey2, ciphertext)
	if err == nil {
		t.Error("Decryption with wrong key should fail")
	}
}

// TestGenerateAsymmetricKeyPair tests asymmetric key pair generation
func TestGenerateAsymmetricKeyPair(t *testing.T) {
	publicKey, privateKey, err := crypto.GenerateAsymmetricKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	if len(publicKey) != crypto.Ed25519PublicKeySize {
		t.Fatalf("Expected public key size %d, got %d", crypto.Ed25519PublicKeySize, len(publicKey))
	}

	if len(privateKey) != crypto.Ed25519PrivateKeySize {
		t.Fatalf("Expected private key size %d, got %d", crypto.Ed25519PrivateKeySize, len(privateKey))
	}

	// Test signing and verification
	message := []byte("test message")
	signature := ed25519.Sign(privateKey, message)

	// Verify signature
	valid, err := crypto.ValidateSignature(publicKey, message, signature)
	if err != nil {
		t.Fatalf("Failed to validate signature: %v", err)
	}

	if !valid {
		t.Error("Valid signature should be verified successfully")
	}

	// Test with wrong message - should fail
	valid, err = crypto.ValidateSignature(publicKey, []byte("wrong message"), signature)
	if err != nil {
		t.Fatalf("Failed to validate signature: %v", err)
	}

	if valid {
		t.Error("Invalid signature should not be verified")
	}
}

// TestValidateSymmetricKey tests symmetric key validation
func TestValidateSymmetricKey(t *testing.T) {
	// Generate master key
	masterKey := make([]byte, 32)
	rand.Read(masterKey)

	// Generate symmetric key
	originalKey, err := crypto.GenerateSymmetricKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Encrypt the key
	encryptedKey, err := crypto.EncryptKey(masterKey, originalKey)
	if err != nil {
		t.Fatalf("Failed to encrypt key: %v", err)
	}

	// Validate with correct key
	valid, err := crypto.ValidateSymmetricKey(masterKey, encryptedKey, originalKey)
	if err != nil {
		t.Fatalf("Failed to validate key: %v", err)
	}

	if !valid {
		t.Error("Valid key should be validated successfully")
	}

	// Generate wrong key
	wrongKey, err := crypto.GenerateSymmetricKey()
	if err != nil {
		t.Fatalf("Failed to generate wrong key: %v", err)
	}

	// Validate with wrong key
	valid, err = crypto.ValidateSymmetricKey(masterKey, encryptedKey, wrongKey)
	if err != nil {
		t.Fatalf("Failed to validate key: %v", err)
	}

	if valid {
		t.Error("Invalid key should not be validated")
	}
}

