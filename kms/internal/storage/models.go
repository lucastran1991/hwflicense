package storage

import "time"

// KeyType represents the type of cryptographic key
type KeyType string

const (
	// KeyTypeSymmetric represents a symmetric key (AES)
	KeyTypeSymmetric KeyType = "symmetric"
	// KeyTypeAsymmetric represents an asymmetric key pair (Ed25519)
	KeyTypeAsymmetric KeyType = "asymmetric"
)

// KeyStatus represents the status of a key
type KeyStatus string

const (
	// KeyStatusActive indicates the key is active and can be used
	KeyStatusActive KeyStatus = "active"
	// KeyStatusRevoked indicates the key has been revoked
	KeyStatusRevoked KeyStatus = "revoked"
)

// Key represents a cryptographic key stored in the system
type Key struct {
	ID                 string     `json:"id"`
	KeyType            KeyType    `json:"key_type"`
	PublicKey          []byte     `json:"public_key,omitempty"`          // Only for asymmetric keys
	EncryptedPrivateKey []byte    `json:"encrypted_private_key"`        // AES-GCM encrypted
	ExpiresAt          time.Time  `json:"expires_at"`
	CreatedAt          time.Time  `json:"created_at"`
	Status             KeyStatus  `json:"status"`
	Version            int        `json:"version"`
}

// IsExpired checks if the key has expired
func (k *Key) IsExpired() bool {
	return time.Now().After(k.ExpiresAt)
}

// IsRevoked checks if the key has been revoked
func (k *Key) IsRevoked() bool {
	return k.Status == KeyStatusRevoked
}

// IsValid checks if the key is valid (active and not expired)
func (k *Key) IsValid() bool {
	return k.Status == KeyStatusActive && !k.IsExpired()
}

