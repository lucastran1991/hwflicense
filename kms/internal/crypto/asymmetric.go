package crypto

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/atprof/license-server/kms/pkg/errors"
)

const (
	// Ed25519PublicKeySize is the size of Ed25519 public keys in bytes
	Ed25519PublicKeySize = 32
	// Ed25519PrivateKeySize is the size of Ed25519 private keys in bytes
	Ed25519PrivateKeySize = 64
	// Ed25519SignatureSize is the size of Ed25519 signatures in bytes
	Ed25519SignatureSize = 64
)

// GenerateAsymmetricKeyPair generates a new Ed25519 key pair
// Returns public key (32 bytes) and private key (64 bytes)
// The private key must be encrypted before storage
func GenerateAsymmetricKeyPair() (publicKey, privateKey []byte, err error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	
	// Ed25519.Public() returns 32 bytes
	// Ed25519.PrivateKey is 64 bytes (seed + public key)
	return pub, priv, nil
}

// ValidateSignature validates an Ed25519 signature
// Returns true if signature is valid, false otherwise
func ValidateSignature(publicKey, message, signature []byte) (bool, error) {
	if len(publicKey) != Ed25519PublicKeySize {
		return false, errors.ErrInvalidKeyMaterial
	}
	
	if len(signature) != Ed25519SignatureSize {
		return false, errors.ErrInvalidSignature
	}
	
	// Use ed25519.Verify function
	valid := ed25519.Verify(publicKey, message, signature)
	return valid, nil
}

