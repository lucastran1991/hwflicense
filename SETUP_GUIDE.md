# TaskMaster License System - Quick Setup Guide

## Current Status

I've implemented the foundational backend infrastructure for the TaskMaster License Management System. This includes:

### ✅ What's Been Implemented

1. **Backend Core**
   - Go module with dependencies (Gin, JWT, SQLite, ECDSA crypto)
   - SQLite database schema (adapted from PostgreSQL)
   - Automatic migration system
   - Configuration management

2. **Cryptographic Operations**
   - ECDSA P-256 key pair generation
   - JSON signing and verification
   - PEM encoding/decoding for keys
   - Helper functions for signature operations

3. **Data Layer**
   - Complete data models (CML, SiteLicense, UsageLedger, UsageStats, UsageManifest)
   - Repository pattern for database operations
   - CML repository (Create, Get, Update, List)
   - Site license repository (Create, Get, List, Update, Delete)
   - JWT authentication middleware

4. **Key Generation Utility**
   - `cmd/genkeys` for generating cryptographic keys

## How to Proceed

### Step 1: Generate Root Keys

```bash
cd backend
go run cmd/genkeys/main.go root
```

This will:
- Create `keys/` directory
- Generate root_private.pem and root_public.pem
- Display the public key for configuration

### Step 2: Set Up Environment

Create `backend/.env` (or use system environment variables):

```env
DB_PATH=data/taskmaster_license.db
JWT_SECRET=your-secret-key-here
API_PORT=8080
API_ENV=development
ROOT_PUBLIC_KEY=<paste from genkeys output>
ASTACK_MOCK_PORT=8081
```

### Step 3: Build and Run

```bash
cd backend
go mod download
go build ./cmd/server
./server
```

The server will:
- Create the database at `data/taskmaster_license.db`
- Run migrations automatically
- Initialize the license management system

## Next Implementation Steps

To complete the full system as specified in the PRD, you need to implement:

### Backend (Continuing from here)
1. **API Handlers** (`internal/api/`)
   - Authentication endpoint (POST /api/auth/login)
   - CML management (GET/POST /api/cml/*)
   - Site license management (GET/POST/DELETE /api/sites/*)
   - License validation (POST /api/validate)
   - Usage ledger API
   - Statistics API
   - Manifest generation and management

2. **Service Layer** (`internal/service/`)
   - CML validation logic
   - Site license minting with chain-of-trust
   - License validation with expiration checks
   - Usage tracking and statistics

3. **Complete Server** (`cmd/server/main.go`)
   - Gin router setup
   - Route handlers
   - Dependency injection
   - Error handling
   - Health check endpoint

4. **Mock A-Stack Server** (`cmd/astack-mock/`)
   - Simple HTTP server for testing
   - CML issuance endpoint
   - Manifest reception and validation

### Frontend (Next Phase)
5. **Next.js Setup**
   - Initialize Next.js 14 with TypeScript
   - Install Shadcn UI components
   - Configure Tailwind CSS
   - Set up API client and React Query

6. **Authentication**
   - Login page
   - Auth context
   - Protected routes

7. **Dashboard**
   - CML status display
   - Site count and statistics
   - Recent activity feed

8. **Site Management UI**
   - List sites with filtering
   - Create site dialog
   - Site detail page
   - Download license file

9. **Statistics & Manifests**
   - Usage statistics dashboard
   - Manifest generation UI
   - Send to A-Stack functionality

## Architecture Overview

```
┌─────────────────────────────────────────┐
│  Frontend (Next.js + Shadcn)            │
│  - Dashboard, Site Management           │
│  - Statistics, Manifests                │
└────────────────┬────────────────────────┘
                 │
                 │ HTTP/REST API
                 ▼
┌─────────────────────────────────────────┐
│  Hub Server (Golang + Gin)              │
│  - CML Storage                           │
│  - Site License Minting                  │
│  - Usage Tracking                        │
│  - Manifest Generation                   │
└────────────────┬────────────────────────┘
                 │
                 │ SQLite
                 ▼
┌─────────────────────────────────────────┐
│  Database (SQLite)                       │
│  - CML, Site Licenses                    │
│  - Usage Ledger, Stats, Manifests       │
└─────────────────────────────────────────┘
```

## File Structure

```
/Users/mac/hwflicense/
├── backend/
│   ├── cmd/
│   │   ├── server/main.go          ✅ Entry point
│   │   └── genkeys/main.go         ✅ Key generation
│   ├── internal/
│   │   ├── api/                     🚧 API handlers (TODO)
│   │   ├── middleware/             ✅ Auth middleware
│   │   ├── models/                  ✅ Data structures
│   │   ├── repository/              ✅ Database layer
│   │   ├── service/                 🚧 Business logic (TODO)
│   │   ├── config/                  ✅ Configuration
│   │   └── database/                ✅ DB connection
│   ├── pkg/
│   │   └── crypto/                  ✅ Cryptographic ops
│   ├── migrations/
│   │   └── 001_initial_schema.sql   ✅ Schema
│   ├── go.mod                       ✅
│   ├── go.sum                       ✅
│   └── README.md                     ✅
├── frontend/                        🚧 Next.js (TODO)
├── projectPRD.md                    ✅ Requirements
└── .gitignore                       ✅
```

## Notes

- All code follows the PRD specifications
- Using SQLite instead of PostgreSQL (as requested)
- Mock A-Stack server will be created for testing
- Root key generation included for development

The foundation is solid and ready for the next development phase.

