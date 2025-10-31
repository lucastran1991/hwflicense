# Key Management Service (KMS-lite)

A standalone Golang backend that acts as an independent key management server. It securely creates, stores, validates, refreshes, and revokes symmetric/asymmetric keys for other internal systems.

## Features

- **Symmetric Key Management**: Generate and validate AES-256 symmetric keys
- **Asymmetric Key Management**: Generate and validate Ed25519 key pairs
- **License File Generation**: Generate and validate `.lic` license files with digital signatures
- **Envelope Encryption**: All private keys encrypted at rest using AES-256-GCM
- **Key Expiry**: Support for TTL and manual key revocation
- **Security Hardening**: Zero memory wiping, secure key handling, rate limiting
- **RESTful API**: HTTP/JSON API for all key operations

## Requirements

- Go 1.23 or higher
- A 32-byte (256-bit) master encryption key (base64 encoded)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd kms
```

2. Install dependencies:
```bash
go mod download
```

3. Build the service:
```bash
go build -o kms-server ./cmd/server
```

## Configuration

The service requires the following environment variables:

- `KMS_MASTER_KEY` (required): Base64-encoded 32-byte master encryption key
- `KMS_DB_PATH` (optional): Path to BoltDB database file (default: `./kms.db`)
- `KMS_PORT` (optional): HTTP server port (default: `:8080`)

### Generating Master Key

Generate a secure master key:

```bash
# Generate 32 random bytes
openssl rand -base64 32

# Or using Go
go run -exec 'echo "package main; import (\"crypto/rand\"; \"encoding/base64\"; \"os\"); func main() { k := make([]byte, 32); rand.Read(k); os.Stdout.WriteString(base64.StdEncoding.EncodeToString(k)) }"'
```

Set the environment variable:

```bash
export KMS_MASTER_KEY=<your-base64-encoded-key>
export KMS_DB_PATH=./kms.db
export KMS_PORT=8080
```

## Running

### Using Management Scripts (Recommended)

The service includes management scripts for easy operation:

**Generate Master Key:**
```bash
# Generate a 32-byte master key (base64 encoded)
openssl rand -base64 32

# Set environment variable
export KMS_MASTER_KEY="<your-generated-key>"
export KMS_DB_PATH="./kms.db"
export KMS_PORT="8080"
```

**Start Service:**
```bash
./start.sh
```

The script will:
- Check if service is already running
- Build the service if needed
- Start the service in background
- Create a PID file (`kms.pid`)
- Write logs to `kms.log`
- Perform a health check

**Stop Service:**
```bash
./stop.sh
```

The script will:
- Send SIGTERM for graceful shutdown
- Wait up to 10 seconds for shutdown
- Force kill if necessary
- Clean up PID file

**Restart Service:**
```bash
./restart.sh
```

**Check Status:**
```bash
./status.sh
```

This shows:
- Process status and PID
- Health check result
- Configuration
- Recent log entries

**View Logs:**
```bash
tail -f kms.log
```

### Manual Running

Start the server manually:

```bash
./kms-server
```

Or directly with Go:

```bash
go run ./cmd/server/main.go
```

The server will start on the configured port (default: `:8080`).

## API Documentation

### Health Check

```
GET /health
```

Returns the health status of the service.

**Response:**
```json
{
  "status": "ok"
}
```

### Register Key

```
POST /keys
```

Register a new key or generate one automatically.

**Request Body:**
```json
{
  "key_type": "symmetric|asymmetric",
  "expires_in_seconds": 31536000,
  "key_material": "base64-encoded-key" // Optional, if not provided, key will be generated
}
```

**Response:**
```json
{
  "key_id": "uuid",
  "key_type": "symmetric|asymmetric",
  "public_key": "base64-encoded-public-key", // Only for asymmetric keys
  "expires_at": "2025-01-01T00:00:00Z",
  "created_at": "2024-01-01T00:00:00Z"
}
```

**Example - Generate Symmetric Key:**
```bash
curl -X POST http://localhost:8080/keys \
  -H "Content-Type: application/json" \
  -d '{
    "key_type": "symmetric",
    "expires_in_seconds": 31536000
  }'
```

**Example - Generate Asymmetric Key Pair:**
```bash
curl -X POST http://localhost:8080/keys \
  -H "Content-Type: application/json" \
  -d '{
    "key_type": "asymmetric",
    "expires_in_seconds": 31536000
  }'
```

### Validate Key

```
POST /keys/validate
```

Validate a symmetric key or verify an Ed25519 signature.

**Request Body (Symmetric):**
```json
{
  "key_id": "uuid",
  "key_material": "base64-encoded-key"
}
```

**Request Body (Asymmetric):**
```json
{
  "key_id": "uuid",
  "message": "message-to-verify",
  "signature": "base64-encoded-signature"
}
```

**Response:**
```json
{
  "valid": true,
  "expired": false,
  "revoked": false
}
```

**Example - Validate Symmetric Key:**
```bash
curl -X POST http://localhost:8080/keys/validate \
  -H "Content-Type: application/json" \
  -d '{
    "key_id": "uuid",
    "key_material": "base64-encoded-key"
  }'
```

**Example - Validate Asymmetric Signature:**
```bash
curl -X POST http://localhost:8080/keys/validate \
  -H "Content-Type: application/json" \
  -d '{
    "key_id": "uuid",
    "message": "Hello, World!",
    "signature": "base64-encoded-signature"
  }'
```

### Refresh Key Expiry

```
POST /keys/:id/refresh
```

Extend the expiry time of a key.

**Request Body:**
```json
{
  "expires_in_seconds": 31536000
}
```

**Response:**
```json
{
  "key_id": "uuid",
  "new_expires_at": "2025-01-01T00:00:00Z"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/keys/{key-id}/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "expires_in_seconds": 31536000
  }'
```

### Remove Key

```
DELETE /keys/:id
```

Revoke a key by setting its status to revoked.

**Response:**
```json
{
  "success": true,
  "key_id": "uuid"
}
```

**Example:**
```bash
curl -X DELETE http://localhost:8080/keys/{key-id}
```

### Generate License File

```
POST /licenses/generate
```

Generate a license file (`.lic`) containing key information and metadata for distribution to clients.

**Request Body:**
```json
{
  "key_id": "uuid",
  "license_type": "enterprise|site|trial|etc",
  "metadata": {
    "customer_id": "CUST001",
    "site_name": "Main Office",
    "max_users": "100"
  }
}
```

**Response:**
```json
{
  "license_file": "base64-encoded-license-content",
  "filename": "enterprise.lic",
  "license_id": "uuid"
}
```

**Example - Generate License for Asymmetric Key:**
```bash
curl -X POST http://localhost:8080/licenses/generate \
  -H "Content-Type: application/json" \
  -d '{
    "key_id": "7ff272a4-1427-4d83-8b72-fb9e1852bf08",
    "license_type": "enterprise",
    "metadata": {
      "customer_id": "CUST001",
      "site_name": "Main Office",
      "max_users": "100"
    }
  }' | jq -r '.license_file' | base64 -d > enterprise.lic
```

**Example - Generate License for Symmetric Key:**
```bash
curl -X POST http://localhost:8080/licenses/generate \
  -H "Content-Type: application/json" \
  -d '{
    "key_id": "uuid-of-symmetric-key",
    "license_type": "site"
  }' | jq -r '.license_file' | base64 -d > site.lic
```

### Validate License File

```
POST /licenses/validate
```

Validate a license file by verifying its signature and checking key status.

**Request (Multipart File Upload):**
```bash
curl -X POST http://localhost:8080/licenses/validate \
  -F "file=@enterprise.lic"
```

**Request (JSON Body with Base64):**
```json
{
  "license_content": "base64-encoded-license-file"
}
```

**Response (Valid License):**
```json
{
  "valid": true,
  "license_id": "uuid",
  "license_type": "enterprise",
  "key_id": "uuid",
  "expires_at": "2026-10-30T17:10:08.206353Z",
  "expired": false,
  "revoked": false,
  "metadata": {
    "customer_id": "CUST001",
    "site_name": "Main Office",
    "max_users": "100"
  }
}
```

**Response (Invalid License):**
```json
{
  "valid": false,
  "expired": false,
  "revoked": false,
  "error": "invalid license signature"
}
```

**Response (Expired License):**
```json
{
  "valid": false,
  "expired": true,
  "revoked": false,
  "license_id": "uuid",
  "key_id": "uuid"
}
```

**Response (Revoked Key):**
```json
{
  "valid": false,
  "expired": false,
  "revoked": true,
  "license_id": "uuid",
  "key_id": "uuid"
}
```

**Example - Validate License File (File Upload):**
```bash
curl -X POST http://localhost:8080/licenses/validate \
  -F "file=@enterprise.lic" | jq .
```

**Example - Validate License File (JSON Body):**
```bash
LICENSE_CONTENT=$(cat enterprise.lic | base64 | tr -d '\n')
curl -X POST http://localhost:8080/licenses/validate \
  -H "Content-Type: application/json" \
  -d "{\"license_content\": \"$LICENSE_CONTENT\"}" | jq .
```

## License File Format

License files (`.lic`) are JSON files containing key information, metadata, and a digital signature for integrity verification.

### License File Structure

```json
{
  "license_id": "550e8400-e29b-41d4-a716-446655440000",
  "license_type": "enterprise",
  "key_id": "7ff272a4-1427-4d83-8b72-fb9e1852bf08",
  "key_type": "asymmetric",
  "public_key": "NsUi47Lx2/Z2vt069PO+Xd+DEXx4OlU2bcyo6wmtiUU=",
  "issued_at": "2025-10-30T17:10:23.215091Z",
  "expires_at": "2026-10-30T17:10:08.206353Z",
  "metadata": {
    "customer_id": "CUST001",
    "site_name": "Main Office",
    "max_users": "100"
  },
  "signature": "HMAC-SHA256-signature-base64-encoded"
}
```

### License File Fields

- **license_id**: Unique identifier for the license file
- **license_type**: Type of license (enterprise, site, trial, etc.)
- **key_id**: Reference to the key stored in KMS database
- **key_type**: Type of key (symmetric or asymmetric)
- **public_key**: Base64-encoded public key (only for asymmetric keys)
- **issued_at**: Timestamp when license was issued
- **expires_at**: Timestamp when license expires
- **metadata**: Custom metadata fields (optional, key-value pairs)
- **signature**: HMAC-SHA256 signature of the license file (excluding signature field)

### Security Features

- **Digital Signature**: Each license file is signed with HMAC-SHA256 using the master key
- **Integrity Verification**: Signature verification ensures license file hasn't been tampered with
- **Expiry Checking**: License validation checks if the license has expired
- **Revocation Support**: License validation verifies if the underlying key has been revoked
- **Metadata Protection**: Metadata is included in the signature to prevent tampering

### Usage Flow

1. **Generate License**: Create a key in KMS, then generate a license file using `/licenses/generate`
2. **Distribute**: Send the `.lic` file to the client
3. **Validate**: Client can validate the license by uploading the file to `/licenses/validate`
4. **Revoke**: If needed, revoke the key in KMS, and all related licenses will become invalid

### Example Workflow

```bash
# 1. Create an asymmetric key
KEY_RESPONSE=$(curl -s -X POST http://localhost:8080/keys \
  -H "Content-Type: application/json" \
  -d '{
    "key_type": "asymmetric",
    "expires_in_seconds": 31536000
  }')

KEY_ID=$(echo $KEY_RESPONSE | jq -r '.key_id')
echo "Created key: $KEY_ID"

# 2. Generate license file
curl -X POST http://localhost:8080/licenses/generate \
  -H "Content-Type: application/json" \
  -d "{
    \"key_id\": \"$KEY_ID\",
    \"license_type\": \"enterprise\",
    \"metadata\": {
      \"customer_id\": \"CUST001\",
      \"site_name\": \"Main Office\"
    }
  }" | jq -r '.license_file' | base64 -d > enterprise.lic

# 3. Validate license file
curl -X POST http://localhost:8080/licenses/validate \
  -F "file=@enterprise.lic" | jq .

# 4. Later, revoke key (all licenses become invalid)
curl -X DELETE http://localhost:8080/keys/$KEY_ID

# 5. Validate again (will show revoked)
curl -X POST http://localhost:8080/licenses/validate \
  -F "file=@enterprise.lic" | jq .
```

## Security Considerations

### Key Material Protection

- **Never logged**: Key material and private keys are never logged, even in debug mode
- **Encrypted at rest**: All private key material is encrypted using AES-256-GCM envelope encryption
- **Memory security**: Keys are zeroed out in memory after use
- **Constant-time comparison**: Key comparisons use constant-time operations to prevent timing attacks

### Rate Limiting

The service implements rate limiting (100 requests per minute per IP) to prevent brute-force attacks.

### Master Key Security

- **Environment variable**: Master key must be loaded from `KMS_MASTER_KEY` environment variable
- **Size validation**: Master key must be exactly 32 bytes (256 bits) after base64 decoding
- **Never exposed**: Master key is never logged or exposed in API responses
- **License signing**: Master key is used to sign license files (HMAC-SHA256)

### License File Security

- **Digital signatures**: All license files are signed with HMAC-SHA256
- **Signature verification**: License validation includes signature verification using constant-time comparison
- **No private keys**: License files never contain private keys, only public keys for asymmetric keys
- **Tamper detection**: Any modification to license file content will invalidate the signature
- **Metadata protection**: Metadata is included in signature calculation to prevent tampering

### Production Recommendations

1. Use a Hardware Security Module (HSM) or cloud KMS for root key operations
2. Implement mTLS or JWT authentication for API access
3. Add audit logging for compliance
4. Monitor for suspicious activity (many validation failures)
5. Regularly rotate the master key
6. Backup encrypted key database regularly

## Development

### Running Tests

```bash
go test ./...
```

### Code Structure

```
core/kms/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── api/                 # HTTP handlers, router, middleware
│   ├── crypto/              # Cryptographic operations
│   ├── licenses/            # License file generation and validation
│   ├── storage/             # BoltDB storage layer
│   └── config/              # Configuration loading
├── pkg/
│   └── errors/              # Error definitions
└── tests/                   # Unit tests
```

## Future Enhancements

- [ ] Integrate with AWS KMS, GCP KMS, or Azure Key Vault
- [ ] Add mTLS/JWT authentication and RBAC
- [ ] Multi-tenant support with key scoping
- [ ] Metrics and alerting for suspicious activity
- [ ] Audit logging for compliance
- [ ] Key rotation and versioning support
- [ ] Client SDKs (Go, Python)

## License

[Add your license here]

