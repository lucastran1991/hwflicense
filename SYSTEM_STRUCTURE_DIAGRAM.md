# System Structure & Architecture

## 🏗️ High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                        A-Stack License Server                       │
│                     (TO BE IMPLEMENTED - Q&A Meeting)                │
│  - Generic reusable component (not HWF-specific)                    │
│  - 7 Core API endpoints                                             │
│  - Monthly key refresh enforcement                                   │
│  - Quarterly stats collection                                       │
└─────────────────────────┬───────────────────────────────────────────┘
                          │
                          │ 1. Create/Update/Delete Site
                          │ 2. Refresh Key (monthly)
                          │ 3. Get Aggregate Stats (quarterly)
                          │ 4. Check Validity
                          │ 5. Send Alerts
                          ↓
┌─────────────────────────────────────────────────────────────────────┐
│                        Hub (This System) ✅                          │
│                   TaskMaster License Management                      │
│                                                                       │
│  Backend (Go):           Frontend (Next.js):                        │
│  ├─ CML Management       ├─ Dashboard                               │
│  ├─ Site License Minting ├─ Site Management                         │
│  ├─ Usage Ledger         ├─ Manifest Viewing                       │
│  ├─ Org Key Mgmt         ├─ License Download                        │
│  └─ Manifest Generation  └─ Stats Reports                          │
└─────────────────────────┬───────────────────────────────────────────┘
                          │
                          │ Site.lic (sub-license)
                          ↓
┌─────────────────────────────────────────────────────────────────────┐
│                           Site Nodes                                 │
│                    (HWF Sites / Boost Sites)                        │
│                                                                       │
│  - Production Keys (HWF)  │  - Dev Keys (Boost)                     │
│  - Monthly refresh        │  - Dev → Production transition          │
│  - Quarterly stats        │  - Watermark for dev keys               │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 📦 Current System Structure (Implemented)

```
hwflicense/
│
├── backend/                          ✅ Complete (100%)
│   ├── cmd/
│   │   ├── server/main.go           # Main API server
│   │   ├── genkeys/main.go          # Key generation utility
│   │   └── astack-mock/main.go     # Mock A-Stack server
│   │
│   ├── internal/
│   │   ├── api/                     # API handlers
│   │   │   ├── auth_handler.go
│   │   │   ├── cml_handler.go
│   │   │   ├── site_handler.go
│   │   │   ├── ledger_handler.go
│   │   │   └── manifest_handler.go
│   │   │
│   │   ├── config/
│   │   │   └── config.go            # Environment config
│   │   │
│   │   ├── database/
│   │   │   └── database.go          # DB connection & migrations
│   │   │
│   │   ├── models/
│   │   │   └── models.go            # Data structures
│   │   │
│   │   ├── repository/              # Data access layer
│   │   │   ├── repository.go        # Interface
│   │   │   ├── cml_repository.go
│   │   │   ├── site_repository.go
│   │   │   ├── ledger_repository.go
│   │   │   ├── manifest_repository.go
│   │   │   └── org_keys_repository.go
│   │   │
│   │   ├── middleware/
│   │   │   └── auth.go              # JWT middleware
│   │   │
│   │   └── service/                 # Business logic
│   │       ├── cml_service.go
│   │       ├── site_service.go
│   │       └── manifest_service.go
│   │
│   ├── pkg/
│   │   └── crypto/
│   │       ├── crypto.go            # ECDSA signing
│   │       └── encryption.go        # AES-256-GCM encryption
│   │
│   ├── migrations/
│   │   └── 001_initial_schema.sql  # Database schema
│   │
│   ├── keys/                        # Key storage
│   │   ├── root_private.pem
│   │   └── root_public.pem
│   │
│   ├── data/
│   │   └── taskmaster_license.db   # SQLite database
│   │
│   ├── go.mod
│   ├── go.sum
│   └── server                       # Compiled binary
│
├── frontend/                         ✅ Core Complete (70%)
│   ├── app/
│   │   ├── dashboard/
│   │   │   ├── layout.tsx
│   │   │   ├── page.tsx            # Home dashboard
│   │   │   ├── sites/
│   │   │   │   ├── page.tsx        # Site list
│   │   │   │   └── [id]/page.tsx   # Site details
│   │   │   └── manifests/
│   │   │       └── page.tsx        # Manifest list
│   │   │
│   │   ├── login/page.tsx
│   │   ├── layout.tsx
│   │   ├── page.tsx
│   │   └── globals.css
│   │
│   ├── lib/
│   │   ├── api-client.ts           # API client
│   │   └── auth-context.tsx        # Auth context
│   │
│   ├── components/
│   │   └── ui/                      # Reusable components
│   │
│   ├── package.json
│   ├── tailwind.config.ts
│   └── tsconfig.json
│
├── scripts/                          ✅ Management Scripts
│   ├── manage.sh                   # Master script
│   ├── backend.sh                  # Backend management
│   ├── frontend.sh                 # Frontend management
│   └── deploy.sh                   # Production deployment
│
└── Documentation/
    ├── Q&A_INDEXED_KNOWLEDGE_BASE.md
    ├── MAIN_TOPICS_FROM_MEETING.md
    ├── KEY_MANAGEMENT_SUMMARY.md
    ├── SYSTEM_COMPLETE.md
    └── projectPRD.md
```

---

## 🔐 Key Management Hierarchy

### Hierarchical Chain of Trust

```
┌─────────────────────┐
│   Root Keys          │  ← A-Stack
│   (A-Stack)          │     Sign CML (Customer Master License)
└──────────┬──────────┘
           │ Signs CML
           ↓
┌─────────────────────┐
│   Org Keys           │  ← Hub (This System)
│   (Hub)              │     Mint site.lic (site licenses)
└──────────┬──────────┘
           │ Signs site.lic
           ↓
┌─────────────────────┐
│   Enterprise Keys    │  ← NEW from Q&A meeting
│   (Veolia)           │     Enterprise-level licensing
└──────────┬──────────┘
           │
           ↓
┌─────────────────────┐
│   Site Keys          │  ← Sites (HWF/Boost)
│   (Production/Dev)   │     Actual usage
└─────────────────────┘
```

### Key Storage

**Root Keys:**
- Location: `backend/keys/`
- Files: `root_private.pem`, `root_public.pem`
- Format: PEM, ECDSA P-256

**Org Keys:**
- Location: Database (`org_keys` table)
- Encrypted: AES-256-GCM
- Types: `dev`, `prod`
- Generated: `go run cmd/genkeys/main.go org <org-id> <dev|prod>`

---

## 🎯 Proposed License Server Structure (From Q&A)

```
A-Stack License Server (TO BE BUILT)
│
├── API Layer (7 Endpoints)
│   ├── 1. POST /sites/create         Create new site license
│   ├── 2. PUT /sites/:id             Update site status
│   ├── 3. DELETE /sites/:id           Remove site
│   ├── 4. POST /keys/refresh         Monthly key refresh
│   ├── 5. GET /stats/aggregate        Quarterly reporting
│   ├── 6. POST /keys/validate         Check key validity
│   └── 7. POST /alerts                Send invalid key alerts
│
├── Key Management
│   ├── Key Generation
│   ├── Key Storage (AWS Secrets Manager / Database)
│   ├── Key Validation
│   └── Key Refresh Logic
│
├── Stats Collection
│   ├── Site Counts (production/dev)
│   ├── User Counts (HWF admins/others)
│   ├── Enterprise Breakdowns
│   └── Quarterly Aggregation
│
└── Security
    ├── Token-based validation
    ├── Monthly refresh enforcement
    ├── AWS Secrets Manager integration
    └── Alert system
```

---

## 🔄 Data Flow

### 1. Site Creation Flow

```
HWF/Boost Application
         │
         │ Request: Create Site (with mode)
         ↓
License Server API
         │
         │ Generate Key (dev/prod)
         ├─→ Store in Database
         ├─→ Store in AWS Secrets Manager (HWF side)
         ├─→ Create audit log
         ↓
Response: Key + Metadata
         │
         ↓
HWF/Boost Application
         │
         └─→ Site activated with key
```

### 2. Monthly Key Refresh Flow

```
Site/Application
         │
         │ Key approaching expiration
         ↓
Call: POST /keys/refresh
         │
         ├─→ Validate old key
         ├─→ Generate new key
         ├─→ Invalidate old key
         ↓
Response: New key (valid 1 month)
         │
         ↓
Store in AWS Secrets Manager
         │
         └─→ Continue operations
```

### 3. Quarterly Stats Flow

```
HWF Application
         │
         │ Quarterly aggregation
         ├─→ Count production sites
         ├─→ Count dev sites
         ├─→ Count users by role
         ├─→ Aggregate by enterprise
         ↓
Call: POST /stats/aggregate
         │
         ├─→ Validate data format
         ├─→ Store JSON manifest
         ├─→ Create report
         ↓
A-Stack License Server
         │
         └─→ Billing calculation
```

### 4. Validation Flow

```
Site/Application
         │
         │ Periodic validation check
         ├─→ Read token from cache
         ├─→ Check expiration
         ↓
If expired → Call: POST /keys/validate
         │
         ├─→ Validate key signature
         ├─→ Check expiration
         ├─→ Check revocation
         ↓
Response: Valid/Invalid
         │
         ├─→ If valid: Continue operations
         └─→ If invalid: Block access + Alert
```

---

## 📊 Database Schema

### Tables

```sql
-- Root-level licensing
cmls
  ├─ id, org_id, license_data, signature
  └─ issue_date, expires_at

-- Organization keys (NEW)
org_keys
  ├─ id, org_id, key_type (dev/prod)
  ├─ private_key_encrypted (AES-256-GCM)
  └─ public_key, created_at

-- Site-level licensing
site_licenses
  ├─ id, site_id, cml_id
  ├─ license_data, signature
  ├─ status (active/revoked)
  └─ created_at, last_seen

-- Usage tracking
usage_ledger
  ├─ id, site_id, timestamp
  ├─ usage_data (JSONB)
  └─ recorded_at

-- Quarterly manifests
usage_manifests
  ├─ id, report_period, manifest_data
  ├─ signature, sent_to_astack
  └─ created_at
```

---

## 🛡️ Security Layers

### Layer 1: Application Security
```
- JWT Authentication
- API Key Validation
- Role-based Access
```

### Layer 2: Key Security
```
- ECDSA P-256 Signing
- AES-256-GCM Encryption
- PBKDF2 Key Derivation (100K iterations)
```

### Layer 3: Storage Security
```
- AWS Secrets Manager (HWF)
- Database encryption (optional)
- Key rotation (monthly)
```

### Layer 4: Network Security
```
- HTTPS only
- Token-based validation
- Cached validation results
```

---

## 🔑 Key Types Summary

| Key Type | Usage | Storage | Refresh | Billing |
|----------|-------|---------|---------|---------|
| **Root Key** | Sign CML | File system | Manual | N/A |
| **Org Key** | Sign site.lic | Database (encrypted) | Manual | N/A |
| **Site Key (Prod)** | Site operations | AWS Secrets Manager | Monthly | ✅ Yes |
| **Site Key (Dev)** | Configuration/testing | AWS Secrets Manager | Monthly | ❌ No |

---

## 📈 Current vs Proposed Architecture

### Current (Implemented)
```
✅ Hub with CML storage
✅ Site license minting
✅ Usage ledger tracking
✅ Manifest generation
❌ License server (TO BE BUILT)
❌ Monthly key refresh
❌ Quarterly stats collection
```

### Proposed (From Q&A Meeting)
```
✅ Generic reusable component
✅ 7 API endpoints
✅ Monthly key refresh enforcement
✅ Quarterly stats collection
✅ Dev vs production key distinction
✅ Token-based validation
✅ Alert system
```

---

## 🎯 Component Responsibilities

### A-Stack License Server (NEW)
- Key generation & refresh
- License provisioning
- Stats collection
- Validation & enforcement
- Alert management

### Hub (This System - Current)
- CML storage
- Site license minting
- Usage tracking
- Manifest generation
- Org key management

### Applications (HWF/Boost)
- Site creation requests
- Key storage (AWS Secrets Manager)
- Monthly refresh calls
- Quarterly stats submission
- Access control enforcement

---

**Status**: Current system is 85% complete. The new license server architecture from Q&A meeting needs to be implemented as a separate, generic reusable component.

