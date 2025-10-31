package tests

import (
	"os"
	"testing"
	"time"

	"github.com/atprof/license-server/kms/internal/storage"
	"github.com/atprof/license-server/kms/pkg/errors"
)

// TestBoltStore tests BoltDB storage operations
func TestBoltStore(t *testing.T) {
	// Create temporary database file
	dbPath := "/tmp/test_kms.db"
	defer os.Remove(dbPath)

	// Create store
	store, err := storage.NewBoltStore(dbPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Test StoreKey
	now := time.Now().UTC()
	expiresAt := now.Add(365 * 24 * time.Hour)

	key := &storage.Key{
		ID:                 "test-key-id",
		KeyType:            storage.KeyTypeSymmetric,
		PublicKey:          nil,
		EncryptedPrivateKey: []byte("encrypted-key-data"),
		ExpiresAt:          expiresAt,
		CreatedAt:          now,
		Status:             storage.KeyStatusActive,
		Version:            1,
	}

	err = store.StoreKey(key)
	if err != nil {
		t.Fatalf("Failed to store key: %v", err)
	}

	// Test GetKey
	retrieved, err := store.GetKey("test-key-id")
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}

	if retrieved.ID != key.ID {
		t.Errorf("Expected ID %s, got %s", key.ID, retrieved.ID)
	}

	if retrieved.KeyType != key.KeyType {
		t.Errorf("Expected KeyType %s, got %s", key.KeyType, retrieved.KeyType)
	}

	if retrieved.Status != key.Status {
		t.Errorf("Expected Status %s, got %s", key.Status, retrieved.Status)
	}

	// Test GetKey - not found
	_, err = store.GetKey("non-existent-key")
	if err == nil {
		t.Error("Expected error for non-existent key")
	}
	if err != errors.ErrKeyNotFound {
		// Check if it matches our custom error
		if err.Error() != "key not found" {
			t.Errorf("Expected key not found error, got: %v", err)
		}
	}

	// Test UpdateKeyExpiry
	newExpiresAt := time.Now().UTC().Add(730 * 24 * time.Hour)
	err = store.UpdateKeyExpiry("test-key-id", newExpiresAt)
	if err != nil {
		t.Fatalf("Failed to update expiry: %v", err)
	}

	retrieved, err = store.GetKey("test-key-id")
	if err != nil {
		t.Fatalf("Failed to get key after update: %v", err)
	}

	// Check that expiry was updated
	if !retrieved.ExpiresAt.Equal(newExpiresAt) {
		t.Errorf("Expected expiry %v, got %v", newExpiresAt, retrieved.ExpiresAt)
	}

	// Check that version was incremented
	if retrieved.Version != 2 {
		t.Errorf("Expected version 2, got %d", retrieved.Version)
	}

	// Test RevokeKey
	err = store.RevokeKey("test-key-id")
	if err != nil {
		t.Fatalf("Failed to revoke key: %v", err)
	}

	retrieved, err = store.GetKey("test-key-id")
	if err != nil {
		t.Fatalf("Failed to get key after revoke: %v", err)
	}

	if retrieved.Status != storage.KeyStatusRevoked {
		t.Errorf("Expected status %s, got %s", storage.KeyStatusRevoked, retrieved.Status)
	}

	if retrieved.Version != 3 {
		t.Errorf("Expected version 3, got %d", retrieved.Version)
	}
}

// TestBoltStoreAsymmetricKey tests storing and retrieving asymmetric keys
func TestBoltStoreAsymmetricKey(t *testing.T) {
	// Create temporary database file
	dbPath := "/tmp/test_kms_asym.db"
	defer os.Remove(dbPath)

	store, err := storage.NewBoltStore(dbPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	publicKey := []byte("public-key-data-32-bytes-long")
	now := time.Now().UTC()

	key := &storage.Key{
		ID:                 "test-asym-key-id",
		KeyType:            storage.KeyTypeAsymmetric,
		PublicKey:          publicKey,
		EncryptedPrivateKey: []byte("encrypted-private-key"),
		ExpiresAt:          now.Add(365 * 24 * time.Hour),
		CreatedAt:          now,
		Status:             storage.KeyStatusActive,
		Version:            1,
	}

	err = store.StoreKey(key)
	if err != nil {
		t.Fatalf("Failed to store asymmetric key: %v", err)
	}

	retrieved, err := store.GetKey("test-asym-key-id")
	if err != nil {
		t.Fatalf("Failed to get asymmetric key: %v", err)
	}

	if len(retrieved.PublicKey) != len(publicKey) {
		t.Errorf("Expected public key length %d, got %d", len(publicKey), len(retrieved.PublicKey))
	}

	if retrieved.KeyType != storage.KeyTypeAsymmetric {
		t.Errorf("Expected KeyType %s, got %s", storage.KeyTypeAsymmetric, retrieved.KeyType)
	}
}

