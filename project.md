# Key Management Service (KMS) - Workspace Index

## Project Overview

The Key Management Service (KMS) is a full-stack application consisting of:
- **Backend**: A Golang-based key management server (KMS) that securely creates, stores, validates, refreshes, and revokes symmetric/asymmetric keys
- **Frontend**: A Next.js 13 web interface for managing keys and licenses

### Architecture

The system follows a client-server architecture:
- **Backend Service** (`kms/`): RESTful API service written in Go, using Gin framework
- **Frontend Interface** (`interface/`): Next.js 13 application with React and ChakraUI
- **Configuration**: Environment-based configuration with optional JSON settings file
- **Storage**: BoltDB for persistent key storage with encryption at rest

### Main Features

- **Symmetric Key Management**: Generate and manage AES-256 symmetric keys
- **Asymmetric Key Management**: Generate and manage Ed25519 key pairs
- **License File Generation**: Create digitally signed `.lic` license files
- **License Validation**: Validate license files with signature verification
- **Envelope Encryption**: All private keys encrypted at rest using AES-256-GCM
- **Key Expiry**: Support for TTL and manual key revocation
- **Security Hardening**: Zero memory wiping, secure key handling, rate limiting

---

## Directory Structure

```
core/
├── kms/                          # Go backend service
│   ├── cmd/
│   │   └── server/
│   │       └── main.go          # Application entry point
│   ├── internal/
│   │   ├── api/                 # HTTP handlers, router, middleware
│   │   │   ├── handlers.go      # API request handlers
│   │   │   ├── middleware.go    # CORS, rate limiting, logging, recovery
│   │   │   └── router.go        # Route definitions
│   │   ├── config/              # Configuration loading
│   │   │   └── config.go        # Environment and file-based config
│   │   ├── crypto/              # Cryptographic operations
│   │   │   ├── asymmetric.go   # Ed25519 key pair generation/verification
│   │   │   ├── envelope.go      # AES-GCM envelope encryption
│   │   │   └── symmetric.go     # AES-256 key generation
│   │   ├── licenses/            # License file generation and validation
│   │   │   ├── generator.go    # License file creation
│   │   │   ├── models.go        # License data structures
│   │   │   ├── signer.go        # HMAC-SHA256 signing
│   │   │   └── validator.go     # License validation logic
│   │   └── storage/             # BoltDB storage layer
│   │       ├── bolt.go          # Database operations
│   │       └── models.go        # Key data structures
│   ├── pkg/
│   │   └── errors/              # Error definitions
│   │       └── errors.go
│   ├── config/
│   │   └── setting.json         # Optional configuration file
│   ├── tests/                   # Unit tests
│   │   ├── config_test.go
│   │   ├── crypto_test.go
│   │   ├── storage_test.go
│   │   ├── run_all_tests.sh
│   │   ├── run_comprehensive_tests.sh
│   │   └── save_test_results.sh
│   ├── start.sh                 # Start service script
│   ├── stop.sh                  # Stop service script
│   ├── restart.sh               # Restart service script
│   ├── status.sh                # Status check script
│   ├── test_all_apis.sh         # API testing script
│   ├── go.mod                    # Go dependencies
│   ├── go.sum                    # Go dependency checksums
│   └── README.md                 # KMS documentation
│
├── interface/                    # Next.js frontend
│   ├── src/
│   │   ├── app/                 # Next.js App Router pages
│   │   │   ├── layout.tsx       # Root layout with ChakraUI provider
│   │   │   ├── page.tsx         # Dashboard/home page
│   │   │   ├── providers.tsx    # React context providers
│   │   │   ├── globals.css      # Global styles
│   │   │   ├── keys/            # Key management page
│   │   │   │   └── page.tsx
│   │   │   └── licenses/        # License management page
│   │   │       └── page.tsx
│   │   ├── components/          # React components
│   │   │   ├── Keys/            # Key management components
│   │   │   │   ├── KeyCard.tsx  # Key display card
│   │   │   │   ├── KeyForm.tsx  # Key creation form
│   │   │   │   ├── KeyList.tsx  # Key list display
│   │   │   │   └── KeyUpload.tsx # Key upload component
│   │   │   ├── Licenses/        # License management components
│   │   │   │   ├── LicenseGenerator.tsx # License generation form
│   │   │   │   └── LicenseValidator.tsx # License validation form
│   │   │   └── Layout/          # Layout components
│   │   │       ├── Header.tsx   # Top navigation header
│   │   │       ├── MainLayout.tsx # Main layout wrapper
│   │   │       └── Sidebar.tsx  # Side navigation
│   │   └── lib/                 # Utilities and helpers
│   │       ├── api/             # API client functions
│   │       │   ├── client.ts    # Axios instance and error handling
│   │       │   ├── keys.ts      # Key API functions
│   │       │   └── licenses.ts  # License API functions
│   │       ├── types/           # TypeScript type definitions
│   │       │   ├── keys.ts      # Key-related types
│   │       │   └── licenses.ts  # License-related types
│   │       └── theme.ts         # ChakraUI theme configuration
│   ├── public/                  # Static assets
│   ├── start.sh                 # Start frontend script
│   ├── stop.sh                  # Stop frontend script
│   ├── test_deploy.sh           # Deployment test script
│   ├── package.json             # Node.js dependencies
│   ├── package-lock.json        # Dependency lock file
│   ├── tsconfig.json            # TypeScript configuration
│   ├── next.config.js           # Next.js configuration
│   └── README.md                # Frontend documentation
│
├── script/                      # Shared management scripts
│   ├── start.sh                 # Start both services
│   ├── stop.sh                  # Stop both services
│   ├── restart.sh               # Restart both services
│   └── status.sh                # Check status of both services
│
└── project.md                   # This workspace index document
```

---

## Backend (KMS) Components

### API Layer (`internal/api/`)

#### Router (`router.go`)
- Sets up Gin router with all routes and middleware
- Routes:
  - `GET /health` - Health check endpoint
  - `GET /keys` - List all keys
  - `POST /keys` - Register/create a key
  - `POST /keys/validate` - Validate a key or signature
  - `POST /keys/:id/refresh` - Refresh key expiry
  - `DELETE /keys/:id` - Revoke a key
  - `GET /keys/:id/download` - Download key information
  - `POST /licenses/generate` - Generate a license file
  - `POST /licenses/validate` - Validate a license file

#### Middleware (`middleware.go`)
- **CORS Middleware**: Handles cross-origin requests, allows localhost origins
- **Rate Limiting**: Token bucket rate limiter (100 requests/minute per IP)
- **Logging Middleware**: Request logging (never logs key material)
- **Recovery Middleware**: Panic recovery with JSON error responses

#### Handlers (`handlers.go`)
- `RegisterKey`: Create new symmetric or asymmetric keys
- `ListKeys`: Retrieve all stored keys
- `ValidateKey`: Validate symmetric keys or Ed25519 signatures
- `RefreshKey`: Extend key expiry time
- `RemoveKey`: Revoke a key
- `DownloadKey`: Download key information
- `GenerateLicense`: Create digitally signed license files
- `ValidateLicense`: Validate license files (multipart or JSON)

### Crypto Module (`internal/crypto/`)

#### Symmetric Keys (`symmetric.go`)
- `GenerateSymmetricKey()`: Generate 32-byte (256-bit) AES keys
- Uses `crypto/rand` for secure random generation

#### Asymmetric Keys (`asymmetric.go`)
- `GenerateAsymmetricKeyPair()`: Generate Ed25519 key pairs
- `SignMessage()`: Sign messages with Ed25519 private key
- `VerifySignature()`: Verify Ed25519 signatures

#### Envelope Encryption (`envelope.go`)
- `EncryptKey()`: Encrypt key material using AES-256-GCM with master key
- `DecryptKey()`: Decrypt key material using AES-256-GCM
- Zeroes out plaintext keys in memory after use

### Storage Layer (`internal/storage/`)

#### Models (`models.go`)
- `Key`: Key data structure with:
  - ID, KeyType (symmetric/asymmetric), Status (active/revoked)
  - PublicKey (for asymmetric keys)
  - EncryptedPrivateKey (AES-GCM encrypted)
  - ExpiresAt, CreatedAt, Version
  - Methods: `IsExpired()`, `IsRevoked()`, `IsValid()`

#### BoltDB Storage (`bolt.go`)
- `BoltStore`: Database interface implementation
- Methods:
  - `CreateKey()`: Store a new key
  - `GetKey()`: Retrieve a key by ID
  - `ListKeys()`: List all keys
  - `UpdateKey()`: Update key metadata
  - `DeleteKey()`: Mark key as revoked
  - `Close()`: Close database connection

### License Management (`internal/licenses/`)

#### Models (`models.go`)
- `LicenseFile`: License file structure with metadata and signature
- `GenerateLicenseRequest/Response`: Request/response structures
- `ValidateLicenseRequest/Response`: Validation structures

#### Generator (`generator.go`)
- Creates license files with:
  - License ID, type, key reference
  - Expiration dates, metadata
  - HMAC-SHA256 signature

#### Signer (`signer.go`)
- `SignLicense()`: Generate HMAC-SHA256 signature of license data
- Uses master key for signing

#### Validator (`validator.go`)
- `ValidateLicense()`: Validate license files
- Verifies signature, checks expiry, checks key revocation status

### Configuration (`internal/config/`)

#### Config (`config.go`)
- Loads configuration from:
  1. Environment variables (highest priority)
  2. JSON settings file (`config/setting.json`)
  3. Default values
- Configuration fields:
  - `KMS_MASTER_KEY`: 32-byte base64-encoded master encryption key (required)
  - `KMS_DB_PATH`: BoltDB database file path (default: `./kms.db`)
  - `KMS_PORT`: HTTP server port (default: `:8080`)
  - `KMS_CONFIG_PATH`: Path to settings JSON file (optional)

### Entry Point (`cmd/server/main.go`)
- Loads configuration
- Initializes BoltDB storage
- Creates API handler
- Sets up HTTP server with graceful shutdown
- Handles SIGINT/SIGTERM for clean shutdown

### Error Handling (`pkg/errors/`)

#### Errors (`errors.go`)
- Standard error definitions:
  - `ErrKeyNotFound`, `ErrKeyExpired`, `ErrKeyRevoked`
  - `ErrInvalidKeyMaterial`, `ErrInvalidSignature`
  - `ErrEncryptionFailed`, `ErrDecryptionFailed`
  - `ErrInvalidLicenseFile`, `ErrLicenseSignatureInvalid`
  - `ErrLicenseExpired`, `ErrLicenseRevoked`

---

## Frontend (Interface) Components

### Pages (`src/app/`)

#### Root Layout (`layout.tsx`)
- Wraps application with ChakraUI provider
- Sets up global theme and styling

#### Dashboard (`page.tsx`)
- Health check status display
- API endpoint information
- Quick start guide
- Backend connection status

#### Keys Page (`keys/page.tsx`)
- Key management interface
- Create, list, and manage keys

#### Licenses Page (`licenses/page.tsx`)
- License generation and validation interface

### Components

#### Key Management (`components/Keys/`)
- **KeyForm**: Form for creating new keys (symmetric/asymmetric)
- **KeyList**: Display list of keys with actions
- **KeyCard**: Individual key display card
- **KeyUpload**: Component for uploading key material

#### License Management (`components/Licenses/`)
- **LicenseGenerator**: Form for generating license files
- **LicenseValidator**: Form/upload for validating license files

#### Layout Components (`components/Layout/`)
- **MainLayout**: Main application layout wrapper
- **Header**: Top navigation bar
- **Sidebar**: Side navigation menu

### API Client Library (`lib/api/`)

#### Client (`client.ts`)
- Axios instance configuration
- Base URL from `NEXT_PUBLIC_API_URL` environment variable
- Error handling utilities
- Health check function

#### Keys API (`keys.ts`)
- `createKey()`: Create a new key
- `listKeys()`: List all keys
- `validateKey()`: Validate a key
- `refreshKey()`: Refresh key expiry
- `revokeKey()`: Revoke a key
- `downloadKey()`: Download key information

#### Licenses API (`licenses.ts`)
- `generateLicense()`: Generate a license file
- `validateLicense()`: Validate a license file (file upload or base64)

### Type Definitions (`lib/types/`)

#### Keys Types (`keys.ts`)
- TypeScript interfaces for key-related data structures
- Request/response types for key operations

#### Licenses Types (`licenses.ts`)
- TypeScript interfaces for license-related data structures
- License file structure, validation results

### Theme Configuration (`lib/theme.ts`)
- ChakraUI theme customization
- Color schemes, component styling

---

## API Endpoints

### Health Check

**GET** `/health`
- Returns service health status
- Response: `{ "status": "ok" }`

### Key Management

**GET** `/keys`
- List all stored keys
- Response: Array of key objects

**POST** `/keys`
- Register or generate a new key
- Request body:
  ```json
  {
    "key_type": "symmetric|asymmetric",
    "expires_in_seconds": 31536000,
    "key_material": "base64-encoded-key" // Optional
  }
  ```
- Response: Key metadata with ID, type, public key (if asymmetric), expiry

**POST** `/keys/validate`
- Validate a symmetric key or verify an Ed25519 signature
- Request body (symmetric):
  ```json
  {
    "key_id": "uuid",
    "key_material": "base64-encoded-key"
  }
  ```
- Request body (asymmetric):
  ```json
  {
    "key_id": "uuid",
    "message": "message-to-verify",
    "signature": "base64-encoded-signature"
  }
  ```
- Response: `{ "valid": true, "expired": false, "revoked": false }`

**POST** `/keys/:id/refresh`
- Extend key expiry time
- Request body: `{ "expires_in_seconds": 31536000 }`
- Response: Updated expiry timestamp

**DELETE** `/keys/:id`
- Revoke a key
- Response: `{ "success": true, "key_id": "uuid" }`

**GET** `/keys/:id/download`
- Download key information
- Response: Key metadata JSON

### License Management

**POST** `/licenses/generate`
- Generate a digitally signed license file
- Request body:
  ```json
  {
    "key_id": "uuid",
    "license_type": "enterprise|site|trial|etc",
    "metadata": { "key": "value" } // Optional
  }
  ```
- Response: Base64-encoded license file content and filename

**POST** `/licenses/validate`
- Validate a license file
- Request: Multipart file upload or JSON with `license_content` (base64)
- Response: Validation result with license details or error

---

## Technology Stack

### Backend (KMS)

- **Language**: Go 1.25.1
- **Web Framework**: Gin 1.11.0
- **Database**: BoltDB 1.4.3 (embedded key-value store)
- **Cryptography**: 
  - Standard library `crypto` packages
  - Ed25519 for asymmetric keys
  - AES-256-GCM for encryption
  - HMAC-SHA256 for license signing
- **UUID Generation**: google/uuid v1.6.0
- **Validation**: go-playground/validator/v10

### Frontend (Interface)

- **Framework**: Next.js 13.5.6 (App Router)
- **UI Library**: ChakraUI 1.8.8
- **Language**: TypeScript 4.9.5
- **React**: 18.2.0
- **HTTP Client**: Axios 1.6.0
- **Icons**: react-icons 5.5.0
- **Animation**: framer-motion 6.5.1

### Development Tools

- **Build Tool**: Next.js built-in (Webpack/Turbopack)
- **Type Checking**: TypeScript compiler
- **Linting**: ESLint with Next.js config

---

## Configuration

### Environment Variables

#### Backend (KMS)

**Required:**
- `KMS_MASTER_KEY`: Base64-encoded 32-byte (256-bit) master encryption key
  - Used for encrypting keys at rest and signing licenses
  - Can also be loaded from secure file: `./secrets/master.key`, `../secrets/master.key`, or `/etc/kms/master.key`

**Optional:**
- `KMS_DB_PATH`: Path to BoltDB database file (default: `./kms.db`)
- `KMS_PORT`: HTTP server port (default: `:8080`)
- `KMS_CONFIG_PATH`: Path to JSON settings file (default: `./config/setting.json`)

#### Frontend (Interface)

**Optional:**
- `NEXT_PUBLIC_API_URL`: Backend API base URL (default: `http://localhost:8080`)

### Configuration Files

#### Backend Settings (`kms/config/setting.json`)
Optional JSON file for non-sensitive configuration:
```json
{
  "kms_db_path": "./kms.db",
  "kms_port": "8080"
}
```
Note: Master key must always be set via environment variable for security.

#### Frontend Config (`interface/next.config.js`)
Next.js configuration file for build settings and webpack configuration.

#### TypeScript Config (`interface/tsconfig.json`)
TypeScript compiler configuration with path aliases (`@/*` maps to `./src/*`).

---

## Management Scripts

### Backend Scripts (`kms/`)

#### `start.sh`
- Checks if service is already running
- Builds service if needed
- Starts service in background
- Creates PID file (`kms.pid`)
- Writes logs to `kms.log`
- Performs health check

#### `stop.sh`
- Sends SIGTERM for graceful shutdown
- Waits up to 10 seconds for shutdown
- Force kills if necessary
- Cleans up PID file

#### `restart.sh`
- Stops then starts the service

#### `status.sh`
- Shows process status and PID
- Performs health check
- Displays configuration
- Shows recent log entries

#### `test_all_apis.sh`
- Tests all API endpoints
- Validates API functionality

### Frontend Scripts (`interface/`)

#### `start.sh`
- Starts Next.js development or production server

#### `stop.sh`
- Stops Next.js server

#### `test_deploy.sh`
- Tests deployment configuration

### Shared Scripts (`script/`)

#### `start.sh`
- Starts both backend and frontend services

#### `stop.sh`
- Stops both services

#### `restart.sh`
- Restarts both services

#### `status.sh`
- Checks status of both services

### Test Scripts (`kms/tests/`)

#### `run_all_tests.sh`
- Runs all unit tests

#### `run_comprehensive_tests.sh`
- Runs comprehensive test suite

#### `save_test_results.sh`
- Saves test results to file

---

## Security Features

### Key Protection
- Keys never logged or printed, even in debug mode
- All private keys encrypted at rest using AES-256-GCM envelope encryption
- Keys zeroed out in memory after use
- Constant-time key comparisons to prevent timing attacks

### Rate Limiting
- 100 requests per minute per IP address
- Token bucket rate limiter implementation
- Prevents brute-force attacks

### Master Key Security
- Must be loaded from environment variable or secure file
- Must be exactly 32 bytes (256 bits) after base64 decoding
- Never exposed in API responses or logs
- Used for envelope encryption and license signing

### License Security
- All license files signed with HMAC-SHA256
- Signature verification using constant-time comparison
- License files never contain private keys
- Tamper detection through signature validation

---

## Development Workflow

### Backend Development

1. Set environment variables:
   ```bash
   export KMS_MASTER_KEY=$(openssl rand -base64 32)
   export KMS_DB_PATH=./kms.db
   export KMS_PORT=8080
   ```

2. Run service:
   ```bash
   cd kms
   ./start.sh
   # Or manually: go run ./cmd/server/main.go
   ```

3. Run tests:
   ```bash
   cd kms
   go test ./...
   ```

### Frontend Development

1. Install dependencies:
   ```bash
   cd interface
   npm install
   ```

2. Set environment variables (create `.env.local`):
   ```env
   NEXT_PUBLIC_API_URL=http://localhost:8080
   ```

3. Run development server:
   ```bash
   npm run dev
   ```

4. Build for production:
   ```bash
   npm run build
   npm start
   ```

### Full Stack Development

1. Start both services:
   ```bash
   ./script/start.sh
   ```

2. Check status:
   ```bash
   ./script/status.sh
   ```

3. Stop both services:
   ```bash
   ./script/stop.sh
   ```

---

## License File Format

License files (`.lic`) are JSON files containing:
- `license_id`: Unique identifier
- `license_type`: Type of license (enterprise, site, trial, etc.)
- `key_id`: Reference to the key in KMS
- `key_type`: Symmetric or asymmetric
- `public_key`: Base64-encoded public key (for asymmetric keys)
- `issued_at`: Timestamp when license was issued
- `expires_at`: Timestamp when license expires
- `metadata`: Custom key-value pairs (optional)
- `signature`: HMAC-SHA256 signature (base64-encoded)

The signature is computed over all fields except the signature field itself, ensuring integrity.

---

## Future Enhancements

- Integration with AWS KMS, GCP KMS, or Azure Key Vault
- mTLS or JWT authentication and RBAC
- Multi-tenant support with key scoping
- Metrics and alerting for suspicious activity
- Audit logging for compliance
- Key rotation and versioning support
- Client SDKs (Go, Python)

---

## Notes

- The backend service runs in release mode (debug disabled)
- CORS is configured for localhost origins (development)
- The frontend is compatible with Node.js 16+ (ChakraUI v1.x requirement)
- All cryptographic operations use constant-time algorithms where applicable
- Database file should be backed up regularly in production environments
