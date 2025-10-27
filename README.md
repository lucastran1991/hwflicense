# TaskMaster License Management System

A complete license management system for the TaskMaster-AI ecosystem with hierarchical license provisioning, cryptographic validation, and usage tracking.

## Features

- **Customer Master License (CML)** management
- **Site License** minting and management
- **Cryptographic Signing** with ECDSA P-256
- **Usage Tracking** and statistics
- **Manifest Generation** for compliance reporting
- **JWT Authentication** with protected routes
- **RESTful API** (18 endpoints)
- **Modern Web UI** with Next.js and Tailwind CSS

## Architecture

```
Frontend (Next.js) → Backend (Golang) → SQLite Database
                       ↓
                Mock A-Stack Server
```

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- npm or yarn

### 1. Generate Root Keys

```bash
cd backend
go run cmd/genkeys/main.go root
```

### 2. Start Backend

```bash
cd backend
go run cmd/server/main.go
```

Backend runs on http://localhost:8080

### 3. Start Mock A-Stack Server (Optional)

```bash
cd backend
go run cmd/astack-mock/main.go
```

Mock A-Stack runs on http://localhost:8081

### 4. Start Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend runs on http://localhost:3000

### 5. Access the Application

1. Open http://localhost:3000
2. Login with `admin` / `admin123`
3. Navigate to Dashboard, Sites, or Manifests

## Project Structure

```
/Users/mac/hwflicense/
├── backend/              # Golang backend
│   ├── cmd/
│   │   ├── server/       # Main Hub server
│   │   └── astack-mock/  # Mock A-Stack server
│   ├── internal/
│   │   ├── api/          # HTTP handlers
│   │   ├── middleware/    # Auth, logging
│   │   ├── service/      # Business logic
│   │   ├── repository/   # Database layer
│   │   ├── models/        # Data models
│   │   ├── config/       # Configuration
│   │   └── database/     # DB connection
│   ├── pkg/
│   │   └── crypto/       # Cryptographic ops
│   └── migrations/        # SQL migrations
│
├── frontend/             # Next.js frontend
│   ├── app/
│   │   ├── dashboard/    # Dashboard pages
│   │   ├── login/        # Login page
│   │   └── layout.tsx    # Root layout
│   └── lib/
│       ├── api-client.ts  # API client
│       └── auth-context.tsx # Auth context
│
└── projectPRD.md         # Product requirements
```

## API Endpoints

### Authentication
- `POST /api/auth/login` - Login and get JWT token

### CML Management
- `POST /api/cml/upload` - Upload CML
- `GET /api/cml` - Get CML info
- `POST /api/cml/refresh` - Refresh CML

### Site License Management
- `POST /api/sites/create` - Create site license
- `GET /api/sites` - List sites
- `GET /api/sites/:site_id` - Get site details
- `DELETE /api/sites/:site_id` - Revoke site
- `POST /api/sites/:site_id/heartbeat` - Update heartbeat

### Manifest Management
- `POST /api/manifests/generate` - Generate manifest
- `GET /api/manifests` - List manifests
- `GET /api/manifests/:id` - Get manifest
- `GET /api/manifests/:id/download` - Download manifest
- `POST /api/manifests/send` - Send to A-Stack

### Other
- `POST /api/validate` - Validate license (public)
- `GET /api/ledger` - Get usage ledger
- `GET /api/health` - Health check

## Configuration

Backend environment variables:
- `DB_PATH` - Database file path (default: data/taskmaster_license.db)
- `JWT_SECRET` - JWT secret key
- `API_PORT` - Server port (default: 8080)
- `ASTACK_MOCK_PORT` - Mock A-Stack port (default: 8081)

Frontend environment variables:
- `NEXT_PUBLIC_API_URL` - Backend API URL (default: http://localhost:8080/api)

## Usage Examples

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### Create Site License
```bash
curl -X POST http://localhost:8080/api/sites/create \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"site_id":"site_001","fingerprint":{"address":"192.168.1.1"}}'
```

### Generate Manifest
```bash
curl -X POST http://localhost:8080/api/manifests/generate \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"period":"2024-01"}'
```

## Development

### Backend
```bash
cd backend
go mod download
go run cmd/server/main.go
```

### Frontend
```bash
cd frontend
npm install
npm run dev
```

### Running Tests
```bash
cd backend
go test ./...
```

## Database Schema

6 tables: cml, site_licenses, usage_ledger, usage_stats, usage_manifests, org_keys

See `backend/migrations/001_initial_schema.sql` for details.

## Security

- JWT authentication
- ECDSA P-256 cryptographic signing
- SQLite with parameterized queries
- HTTPS recommended for production

## Deployment

See deployment guide for AWS EC2 in the project documentation.

## License

Proprietary - TaskMaster License Management System

