package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/bbolt"

	"github.com/atprof/license-server/kms/pkg/errors"
)

const (
	// KeysBucket is the name of the bucket storing keys
	KeysBucket = "keys"
)

// BoltStore implements the storage interface using BoltDB
type BoltStore struct {
	db *bbolt.DB
}

// NewBoltStore creates a new BoltDB storage instance
func NewBoltStore(dbPath string) (*BoltStore, error) {
	db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	store := &BoltStore{db: db}

	// Initialize the keys bucket
	if err := store.initBucket(); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

// Close closes the database connection
func (s *BoltStore) Close() error {
	return s.db.Close()
}

// initBucket initializes the keys bucket if it doesn't exist
func (s *BoltStore) initBucket() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(KeysBucket))
		return err
	})
}

// StoreKey stores a key in the database
func (s *BoltStore) StoreKey(key *Key) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(KeysBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", KeysBucket)
		}

		data, err := json.Marshal(key)
		if err != nil {
			return fmt.Errorf("failed to marshal key: %w", err)
		}

		return bucket.Put([]byte(key.ID), data)
	})
}

// GetKey retrieves a key from the database by ID
func (s *BoltStore) GetKey(keyID string) (*Key, error) {
	var key *Key
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(KeysBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", KeysBucket)
		}

		data := bucket.Get([]byte(keyID))
		if data == nil {
			return errors.ErrKeyNotFound
		}

		var k Key
		if err := json.Unmarshal(data, &k); err != nil {
			return fmt.Errorf("failed to unmarshal key: %w", err)
		}

		key = &k
		return nil
	})

	return key, err
}

// UpdateKeyExpiry updates the expiry time of a key
func (s *BoltStore) UpdateKeyExpiry(keyID string, newExpiry time.Time) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(KeysBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", KeysBucket)
		}

		data := bucket.Get([]byte(keyID))
		if data == nil {
			return errors.ErrKeyNotFound
		}

		var key Key
		if err := json.Unmarshal(data, &key); err != nil {
			return fmt.Errorf("failed to unmarshal key: %w", err)
		}

		key.ExpiresAt = newExpiry
		key.Version++

		data, err := json.Marshal(&key)
		if err != nil {
			return fmt.Errorf("failed to marshal key: %w", err)
		}

		return bucket.Put([]byte(keyID), data)
	})
}

// RevokeKey revokes a key by setting its status to revoked
func (s *BoltStore) RevokeKey(keyID string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(KeysBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", KeysBucket)
		}

		data := bucket.Get([]byte(keyID))
		if data == nil {
			return errors.ErrKeyNotFound
		}

		var key Key
		if err := json.Unmarshal(data, &key); err != nil {
			return fmt.Errorf("failed to unmarshal key: %w", err)
		}

		key.Status = KeyStatusRevoked
		key.Version++

		data, err := json.Marshal(&key)
		if err != nil {
			return fmt.Errorf("failed to marshal key: %w", err)
		}

		return bucket.Put([]byte(keyID), data)
	})
}

// ListKeys lists all keys in the database (for audit purposes)
// Returns keys without private key material
func (s *BoltStore) ListKeys() ([]*Key, error) {
	var keys []*Key
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(KeysBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", KeysBucket)
		}

		return bucket.ForEach(func(k, v []byte) error {
			var key Key
			if err := json.Unmarshal(v, &key); err != nil {
				return fmt.Errorf("failed to unmarshal key: %w", err)
			}

			// Clear private key material for audit
			key.EncryptedPrivateKey = nil
			keys = append(keys, &key)
			return nil
		})
	})

	return keys, err
}

