# Organization Key Management Implementation

## Overview

Successfully implemented secure organization key management with AES-256-GCM encryption at rest. This is Task #1 of the license management system completion.

## What Was Implemented

### 1. Encryption Functions (`backend/pkg/crypto/encryption.go`)

**New Functions:**
- `EncryptPrivateKey(keyBytes []byte, password string) (string, error)`
  - Encrypts private key data using AES-256-GCM
  - Uses PBKDF2 key derivation with 100,000 iterations
  - Returns base64-encoded encrypted data
  - Includes salt and nonce in encrypted output

- `DecryptPrivateKey(encrypted string, password string) ([]byte, error)`
  - Decrypts encrypted private key data
  - Extracts salt and nonce, derives key, and decrypts
  - Returns plaintext private key bytes

- `ValidatePassword(password string) error`
  - Validates password meets minimum requirements (16+ characters)

**Security Features:**
- AES-256-GCM (Authenticated Encryption with Associated Data)
- PBKDF2-SHA256 with 100,000 iterations
- Random salt per encryption (32 bytes)
- Random nonce per encryption (12 bytes)
- No plaintext private keys stored in database

### 2. Repository Layer (`backend/internal/repository/org_keys_repository.go`)

**New Methods:**
- `CreateOrgKey(orgKey *models.OrgKey) error` - Create new org key
- `GetOrgKey(orgID, keyType string) (*models.OrgKey, error)` - Retrieve by org and type
- `GetOrgKeyByID(id string) (*models.OrgKey, error)` - Retrieve by ID
- `ListOrgKeys(orgID string) ([]models.OrgKey, error)` - List all keys for org
- `DeleteOrgKey(id string) error` - Delete org key
- `CreateOrgKeyWithID(orgKey *models.OrgKey) error` - Create with specific ID

**Features:**
- Proper error handling with context
- UUID-based unique IDs
- Timestamp tracking (created_at)
- SQL injection prevention via parameterized queries

### 3. Configuration Updates (`backend/internal/config/config.go`)

**Added:**
- `EncryptionPassword` field to Config struct
- Environment variable: `ENCRYPTION_PASSWORD`
- Default value: `change-this-password-in-production-12345`
- Loaded from `.env` file

**Usage:**
```bash
export ENCRYPTION_PASSWORD="your-secure-password-minimum-16-characters"
```

### 4. Enhanced Key Generation Utility (`backend/cmd/genkeys/main.go`)

**New Commands:**
- `go run cmd/genkeys/main.go root` - Generate root keys (file storage)
- `go run cmd/genkeys/main.go org <org-id> <dev|prod>` - Generate org keys (database storage)

**Features:**
- Generates ECDSA P-256 key pairs
- Encrypts private key with AES-256-GCM
- Stores in database with org_id and key_type
- Prevents duplicate key generation
- Displays public key for verification

## Usage Examples

### Generate Root Keys (File Storage)
```bash
cd backend
go run cmd/genkeys/main.go root
```

### Generate Organization Keys (Database Storage)
```bash
# Generate dev key for organization
go run cmd/genkeys/main.go org org_12345 dev

# Generate prod key for organization
go run cmd/genkeys/main.go org org_12345 prod
```

### Using Generated Keys

The org keys are now available for use in services:

```go
// In service layer (to be implemented in next task)
repo := repository.NewRepository(db)
orgKey, err := repo.GetOrgKey(orgID, "dev")
if err != nil {
    return err
}

// Decrypt private key
privateKeyPEM, err := crypto.DecryptPrivateKey(orgKey.PrivateKeyEncrypted, config.AppConfig.EncryptionPassword)
if err != nil {
    return err
}

// Load private key for signing
privateKey, err := crypto.LoadPrivateKeyFromPEM(privateKeyPEM)
if err != nil {
    return err
}

// Use private key for signing operations
signature, err := crypto.SignData(data, privateKey)
```

## Security Considerations

1. **Encryption at Rest**: Private keys are never stored in plaintext
2. **Strong Key Derivation**: PBKDF2 with 100,000 iterations prevents brute force
3. **Authenticated Encryption**: GCM mode provides authenticity guarantees
4. **Unique Nonces**: Each encryption uses a unique nonce
5. **Password Requirements**: Minimum 16 character password enforced

## Next Steps (Task #2)

The org key management is now complete. The next task will:
1. Use these org keys to sign site licenses
2. Replace placeholder signatures with real ECDSA signatures
3. Implement signature verification chain
4. Add fingerprint matching logic

## Testing

Tested successfully:
- ✅ Org key generation works
- ✅ Encryption/decryption works
- ✅ Duplicate prevention works
- ✅ Database storage works
- ✅ Public key display works
- ✅ No lint errors

## Files Created/Modified

**Created:**
- `backend/pkg/crypto/encryption.go` - Encryption functions
- `backend/internal/repository/org_keys_repository.go` - Repository methods

**Modified:**
- `backend/internal/config/config.go` - Added encryption password config
- `backend/cmd/genkeys/main.go` - Added org key generation

**Dependencies:**
- `golang.org/x/crypto/pbkdf2` - Added via `go get`

## Status

✅ **Task #1 Complete**
- All verification criteria met
- Org keys can be created, encrypted, stored, and retrieved
- Encryption uses AES-256-GCM
- Private keys never stored in plaintext
- Key rotation possible (via delete + regenerate)
- Ready for use in Task #2

