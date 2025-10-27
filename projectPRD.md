# Product Requirements Document (PRD)
# Task-Master-AI License Management System

## 1. Executive Summary

Task-Master-AI License Management System is a secure, hierarchical license provisioning and compliance tracking platform for the HWF ecosystem. It enables A-Stack to issue Customer Master Licenses (CML) to Hub operators, who can then mint site-specific sub-licenses for deployment across site nodes.

Purpose: centralized license management with offline-first operations, cryptographic chain-of-trust validation, and monthly compliance reporting.

Stakeholders: A-Stack License Server (home/issuer), Hub operators, Site nodes.

Timeline: 12–16 weeks. Priority: High.

## 2. Product Overview

### 2.1 Business Objectives
- Provide centralized license issuance and management
- Enable Hub-based sub-licensing without always-on connectivity
- Track usage and compliance via monthly manifests
- Support Dev and Prod environments with separate keys
- Enforce licensing constraints (max sites, feature packs, expiration)

### 2.2 Success Metrics
- 100% of issued licenses cryptographically validated
- Zero license bypass incidents
- Monthly manifest accuracy: >99.9%
- License validation latency: <100ms
- 99.9% API uptime on AWS EC2

## 3. System Architecture

### 3.1 High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│                 A-Stack License Server                  │
│  (Home/Issuer) - Generates CML, Validates Manifests     │
└────────────────────────┬────────────────────────────────┘
                         │
                         │ Issues CML
                         │ (Signed)
                         ↓
┌─────────────────────────────────────────────────────────┐
│                      Hub (Operator)                      │
│  - Stores CML                                            │
│  - Mint site.lic (sub-licenses)                          │
│  - Maintain usage ledger                                 │
│  - Generate monthly usage manifests                      │
└────────────────────────┬────────────────────────────────┘
                         │
                         │ site.lic per site
                         ↓
┌─────────────────────────────────────────────────────────┐
│                    Site Nodes (HWF)                      │
│  - Load site.lic                                         │
│  - Validate chain of trust                               │
│  - Verify constraints                                    │
│  - Send heartbeat to Hub (optional)                      │
└─────────────────────────────────────────────────────────┘
```

### 3.2 Component Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Backend API** | Golang (v1.21+) | RESTful APIs for license operations, crypto operations |
| **Frontend** | Next.js 14+ with Shadcn UI | Admin dashboard for Hub operators |
| **Database** | PostgreSQL 15+ | Store CML, site licenses, usage ledger, manifests |
| **Crypto Library** | Golang `crypto/ecdsa` | ECDSA P-256 signatures |
| **Deployment** | AWS EC2 (t3.medium+) | Backend + frontend hosting |
| **Authentication** | JWT tokens | API authentication |
| **File Storage** | AWS S3 (optional) | Store usage manifest exports |

## 4. Functional Requirements

### 4.1 FR-1: License Provisioning

- FR-1.1: A-Stack License Server issues a signed CML
- FR-1.2: Hub validates CML with A-Stack public key
- FR-1.3: Hub stores CML constraints (org_id, max_sites, validity, feature_packs)
- FR-1.4: Support Dev Key and Prod Key
- FR-1.5: Update CML with refreshed keys

Acceptance criteria:
- CML signed with ECDSA P-256
- Validation fails if signature is invalid
- Hub dashboard shows active CML status

### 4.2 FR-2: Site License Creation

- FR-2.1: Hub UI/CLI allows entering site_id and optional fingerprint fields
- FR-2.2: Hub generates site.lic and signs with org key
- FR-2.3: site.lic chains back to root (CML → Root)
- FR-2.4: Fingerprints: address, DNS suffix, deployment tag
- FR-2.5: Enforce max_sites
- FR-2.6: Download site.lic as JSON

Acceptance criteria:
- Each site.lic traces back to A-Stack root
- Cannot exceed max_sites
- Fingerprints are optional but stored when provided

### 4.3 FR-3: Runtime License Validation

- FR-3.1: Site node verifies signature chain (Root → CML → site)
- FR-3.2: Validate expiration
- FR-3.3: Validate fingerprint matches
- FR-3.4: Enforce feature_packs
- FR-3.5: Log validation results

Acceptance criteria:
- Access blocked if validation fails
- Expired licenses trigger a 30-day grace period
- Site fingerprint mismatch logs an audit event

### 4.4 FR-4: Usage Statistics & Ledger

- FR-4.1: User stats: Veolia Admin Users
- FR-4.2: Site stats: Boost Sites (Active/Basic), HWF Sites (Active/Basic)
- FR-4.3: Track site activity (issued_at, last_seen)
- FR-4.4: Store signed ledger of active sites
- FR-4.5: Update on site creation, deletion, heartbeat

Acceptance criteria:
- Real-time counts update
- Ledger entries are immutable and signed
- Accuracy matches dashboard totals

### 4.5 FR-5: Usage Manifest Generation

- FR-5.1: Generate monthly Usage Manifest (JSON + signature)
- FR-5.2: Include active site count and details
- FR-5.3: Include user and site statistics
- FR-5.4: Sign manifest with Hub org key
- FR-5.5: Export manifest for email or upload
- FR-5.6: View manifest history

Acceptance criteria:
- Monthly manifest reflects current period
- JSON validates and can be verified
- Exports include required fields

### 4.6 FR-6: Call Home / Manifest Upload

- FR-6.1: Hub sends monthly manifest to A-Stack License Server
- FR-6.2: A-Stack validates manifest signature
- FR-6.3: A-Stack records compliance
- FR-6.4: Alert Hub operator for expired CML

Acceptance criteria:
- Status response records delivery
- Valid manifests trigger no alerts
- Invalid signatures are rejected

## 5. Data Models

### 5.1 Database Schema

```sql
-- Customer Master License
CREATE TABLE cml (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id VARCHAR(100) UNIQUE NOT NULL,
    max_sites INTEGER NOT NULL,
    validity TIMESTAMP NOT NULL,
    feature_packs JSONB,
    dev_key_public TEXT NOT NULL,
    prod_key_public TEXT NOT NULL,
    cml_data JSONB NOT NULL,
    signature TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Site Licenses
CREATE TABLE site_licenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id VARCHAR(100) UNIQUE NOT NULL,
    org_id VARCHAR(100) REFERENCES cml(org_id),
    fingerprint JSONB, -- {address, dns_suffix, deployment_tag}
    license_data JSONB NOT NULL,
    signature TEXT NOT NULL,
    issued_at TIMESTAMP DEFAULT NOW(),
    last_seen TIMESTAMP,
    status VARCHAR(20) DEFAULT 'active', -- active, expired, revoked
    created_at TIMESTAMP DEFAULT NOW()
);

-- Usage Ledger
CREATE TABLE usage_ledger (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id VARCHAR(100) NOT NULL,
    entry_type VARCHAR(50), -- site_created, site_deleted, heartbeat
    site_id VARCHAR(100),
    data JSONB, -- flexible event data
    signature TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Usage Statistics
CREATE TABLE usage_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id VARCHAR(100) NOT NULL,
    period VARCHAR(10) NOT NULL, -- YYYY-MM format
    user_stats JSONB, -- {admin_users: count}
    site_stats JSONB, -- {boost: {active, basic}, hwf: {active, basic}}
    total_active_sites INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(org_id, period)
);

-- Usage Manifests
CREATE TABLE usage_manifests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id VARCHAR(100) NOT NULL,
    period VARCHAR(10) NOT NULL,
    manifest_data JSONB NOT NULL,
    signature TEXT NOT NULL,
    sent_to_astack BOOLEAN DEFAULT FALSE,
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Key Storage (Encrypted)
CREATE TABLE org_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id VARCHAR(100) UNIQUE NOT NULL,
    key_type VARCHAR(20) NOT NULL, -- dev, prod
    private_key_encrypted TEXT NOT NULL, -- AES-256 encrypted
    public_key TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_site_licenses_org_id ON site_licenses(org_id);
CREATE INDEX idx_site_licenses_status ON site_licenses(status);
CREATE INDEX idx_usage_ledger_org_id ON usage_ledger(org_id);
CREATE INDEX idx_usage_stats_org_id ON usage_stats(org_id, period);
CREATE INDEX idx_usage_manifests_org_id ON usage_manifests(org_id, period);
```

### 5.2 JSON Data Structures

```json
// CML Data (cml_data field)
{
  "type": "customer_master_license",
  "org_id": "org_12345",
  "max_sites": 100,
  "validity": "2025-12-31T23:59:59Z",
  "feature_packs": ["basic", "advanced", "analytics"],
  "key_type": "prod", // or "dev"
  "issued_by": "astack_root",
  "issuer_public_key": "-----BEGIN PUBLIC KEY-----\n...",
  "issued_at": "2024-01-15T10:00:00Z"
}

// Site License Data (license_data field)
{
  "type": "site_license",
  "site_id": "site_abc123",
  "parent_cml": "org_12345",
  "parent_cml_sig": "...",
  "fingerprint": {
    "address": "192.168.1.100",
    "dns_suffix": "hwf.local",
    "deployment_tag": "production-cluster-1"
  },
  "issued_at": "2024-01-16T08:30:00Z",
  "expires_at": "2025-12-31T23:59:59Z",
  "features": ["basic", "advanced"]
}

// Usage Manifest (manifest_data field)
{
  "type": "usage_manifest",
  "org_id": "org_12345",
  "period": "2024-01",
  "generated_at": "2024-02-01T00:00:00Z",
  "stats": {
    "users": {
      "admin_users": 5,
      "total_users": 120
    },
    "sites": {
      "boost": {
        "active": 25,
        "basic": 8
      },
      "hwf": {
        "active": 50,
        "basic": 15
      }
    },
    "total_active_sites": 98
  },
  "active_sites": [
    {
      "site_id": "site_abc123",
      "issued_at": "2024-01-16T08:30:00Z",
      "last_seen": "2024-01-30T14:22:00Z",
      "status": "active"
    }
    // ... more sites
  ]
}
```

## 6. API Specifications

### 6.1 Authentication

```
GET /api/auth/login
Headers: Content-Type: application/json
Body: { "username": "...", "password": "..." }
Response: { "token": "jwt...", "expires_in": 3600 }
```

### 6.2 CML Management

```
GET /api/cml
Headers: Authorization: Bearer {token}
Response: { "cml": {...}, "status": "active" }

POST /api/cml/upload
Headers: Authorization: Bearer {token}
Body: { "cml_data": {...}, "signature": "..." }
Response: { "status": "valid", "org_id": "..." }

POST /api/cml/refresh
Headers: Authorization: Bearer {token}
Body: { "refreshed_cml_data": {...}, "signature": "..." }
Response: { "status": "refreshed" }
```

### 6.3 Site License Management

```
GET /api/sites
Headers: Authorization: Bearer {token}
Query: ?status=active&limit=50&offset=0
Response: { "sites": [...], "total": 100, "limit": 50 }

POST /api/sites/create
Headers: Authorization: Bearer {token}
Body: {
  "site_id": "site_abc123",
  "fingerprint": {
    "address": "192.168.1.100",
    "dns_suffix": "hwf.local",
    "deployment_tag": "prod-cluster-1"
  }
}
Response: { "license": {...}, "site.lic": "..." }

GET /api/sites/{site_id}
Headers: Authorization: Bearer {token}
Response: { "license": {...} }

DELETE /api/sites/{site_id}
Headers: Authorization: Bearer {token}
Response: { "status": "revoked" }

POST /api/sites/{site_id}/heartbeat
Headers: Authorization: Bearer {token}
Body: { "timestamp": "2024-01-30T14:22:00Z" }
Response: { "status": "updated" }
```

### 6.4 License Validation (External API for Site Nodes)

```
POST /api/validate
Body: { "license": {...}, "fingerprint": {...} }
Response: {
  "valid": true,
  "message": "License valid",
  "features": ["basic", "advanced"],
  "expires_at": "2025-12-31T23:59:59Z"
}
```

### 6.5 Usage Statistics

```
GET /api/stats
Headers: Authorization: Bearer {token}
Query: ?period=2024-01
Response: { "stats": {...} }

GET /api/ledger
Headers: Authorization: Bearer {token}
Query: ?limit=100&offset=0
Response: { "entries": [...] }
```

### 6.6 Manifest Management

```
GET /api/manifests
Headers: Authorization: Bearer {token}
Query: ?period=2024-01
Response: { "manifests": [...] }

POST /api/manifests/generate
Headers: Authorization: Bearer {token}
Body: { "period": "2024-01" }
Response: { "manifest": {...}, "signature": "..." }

POST /api/manifests/send
Headers: Authorization: Bearer {token}
Body: { "manifest_id": "uuid", "astack_endpoint": "https://..." }
Response: { "status": "sent", "timestamp": "..." }

GET /api/manifests/{manifest_id}/download
Headers: Authorization: Bearer {token}
Response: File download (JSON)
```

## 7. Security Requirements

### 7.1 Cryptographic Requirements
- ECDSA P-256 with SHA-256 for signatures
- Key derivation via PBKDF2-SHA256; strong random seed
- Private keys encrypted at rest (AES-256-GCM)
- TLS 1.3 for all HTTPS
- Signatures verified before acceptance

### 7.2 Access Control
- JWT with 1-hour expiration
- Multi-factor auth for admin actions
- Role-based access control:
  - Admin: CML management, manifest generation
  - Operator: site license operations
  - Viewer: read-only stats

### 7.3 Data Security
- Secrets stored in environment variables (not code)
- DB encrypted at rest (AES-256)
- Audit trail for CML changes and site licenses
- Rate limiting on sensitive APIs
- SQL injection prevention via parameterized queries

### 7.4 Validation Requirements
- Signature chain verification Root → CML → Site
- Expiration checks with grace period
- Fingerprint matching on optional fields
- Max site limits enforced

## 8. User Interface Requirements

### 8.1 Dashboard (Hub Operator)
- CML status
- Current site count (active/max)
- License expiration warnings
- Quick create site license
- Recent activity

### 8.2 Site License Management
- List sites (filter, search, pagination)
- Create site (form: site_id, fingerprint)
- View details (license data, chain, timestamps)
- Revoke, download license file
- Bulk import

### 8.3 Statistics & Reporting
- User stats (Veolia admin users)
- Site stats (Boost Active/Basic, HWF Active/Basic)
- Date range filters
- Export CSV/PDF
- Usage trends

### 8.4 Manifest Management
- Generate for selected month
- Preview JSON
- Send to A-Stack
- Download manifest file
- View send history and status

## 9. Deployment & Infrastructure

### 9.1 AWS EC2 Configuration
- Instance: t3.medium or larger (2 vCPU, 4 GB RAM minimum)
- OS: Ubuntu 22.04 LTS
- Security groups:
  - 80/tcp, 443/tcp, 22/tcp, 5432/tcp (Postgres)
- EBS: 20 GB gp3

### 9.2 PostgreSQL Configuration
- Version: 15+
- Connection string: `postgresql://tql:tql12345@localhost:5432/taskmaster_license`
- Database: `taskmaster_license`
- Pooling: pgBouncer
- Backups: daily snapshots

### 9.3 Application Deployment
```
/var/www/taskmaster-license/
├── backend/          # Golang binary + configs
├── frontend/         # Next.js build
├── nginx.conf        # Reverse proxy
├── systemd/          # Service files
└── scripts/          # Deployment scripts
```

### 9.4 Environment Variables
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=tql
DB_PASSWORD=tql12345
DB_NAME=taskmaster_license

# JWT
JWT_SECRET=your-secret-key-here

# API
API_PORT=8080
API_ENV=production

# Frontend
NEXT_PUBLIC_API_URL=https://your-domain.com/api

# AWS (if using S3)
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=...
AWS_SECRET_ACCESS_KEY=...

# Root Public Key (for CML validation)
ROOT_PUBLIC_KEY=...
```

## 10. Technical Stack Details

### 10.1 Backend (Golang)
- Version: Go 1.21+
- Libraries:
  - `github.com/gin-gonic/gin` — web framework
  - `github.com/golang-jwt/jwt` — JWT auth
  - `github.com/lib/pq` — PostgreSQL
  - `crypto/ecdsa` — signatures
  - `github.com/spf13/viper` — config

### 10.2 Frontend (Next.js + Shadcn)
- Version: Next.js 14+
- Libraries:
  - React 18
  - Shadcn UI
  - Tailwind CSS
  - Axios
  - React Query
  - React Hook Form + Zod
  - date-fns

### 10.3 Database
- PostgreSQL 15+
- Connection string: `postgresql://tql:tql12345@localhost:5432/taskmaster_license`

## 11. Development Timeline

### Phase 1: Core Infrastructure (Weeks 1-4)
- Week 1: DB schema; CML data models; crypto helpers
- Week 2: CML upload/validation; org key generation; signature chain verification
- Week 3: Site license creation; Hub org signing; constraints; fingerprint validation
- Week 4: License validation API; heartbeat; index CML/sites

Deliverables: backend for license operations; DB with seed

### Phase 2: Frontend & UI (Weeks 5-8)
- Week 5: Next.js; Shadcn; auth flow; layouts
- Week 6: CML dashboard; site list, create, actions
- Week 7: Statistics; manifest generation and management
- Week 8: Responsive UI; error handling

Deliverables: functional UI; Hub workflows

### Phase 3: Manifests & Reporting (Weeks 9-12)
- Week 9: Usage ledger
- Week 10: Monthly manifest generation and export
- Week 11: A-Stack send/status tracking
- Week 12: Stats: User (Veolia admin), Sites (Boost/HWF Active/Basic)

Deliverables: report generation; A-Stack integration

### Phase 4: Testing & Deployment (Weeks 13-16)
- Week 13: Unit tests
- Week 14: Integration and E2E
- Week 15: AWS EC2 provisioning and deployment
- Week 16: Production smoke tests; monitoring

Deliverables: production deployment; documentation

## 12. Testing Strategy

### 12.1 Unit Tests
- Crypto: signature generation/verification
- Validation: constraints, expiration, fingerprints
- Data: model serialization

### 12.2 Integration Tests
- CML upload → store → validation
- Site creation → mint license → verify chain
- Manifest generation → signature → send

### 12.3 End-to-End Tests
- CML upload → create 10 sites → generate manifest → send to A-Stack
- Fingerprint enforcement
- Max sites limit
- Expiration with grace period

## 13. Monitoring & Observability

### 13.1 Metrics
- API response times and errors
- License validation success/failures
- Active sites (active/expired)
- Manifest generation timestamps
- DB connection pool

### 13.2 Logging
- Structured JSON logs (Go logrus)
- Info: successful operations
- Warning: expiring licenses
- Error: validation failures
- Audit: CML changes, site license actions

### 13.3 Alerts
- License expiration within 7 days
- Max sites threshold (>90%)
- Manifest not sent for 45 days
- API errors exceed threshold

## 14. Documentation Requirements

### 14.1 Technical
- API docs (OpenAPI/Swagger)
- DB schema
- Deployment guide
- Architecture diagrams

### 14.2 User
- Hub operator guide
- Site node integration
- Troubleshooting

## 15. Acceptance Criteria

### 15.1 Functional
- Upload and validate CML
- Create up to 1000 sites
- Verify signature chain (Root → CML → Site)
- Generate monthly manifests
- Send manifests to A-Stack
- Display user and site stats
- Enforce max sites limit

### 15.2 Non-Functional
- API response <100ms
- UI load <2s
- 99.9% uptime
- Secrets stored securely
- Signatures verified
- GDPR-ready

## 16. Risk Assessment & Mitigation

| Risk | Impact | Probability | Mitigation |
|------|--------|------------|------------|
| License bypass attempts | High | Low | Strong signature verification, fingerprint enforcement |
| CML expiration without refresh | High | Medium | Automatic alerts 7 days before, dashboard warnings |
| Database corruption | High | Low | Automated backups, point-in-time recovery |
| AWS EC2 instance failure | Medium | Low | Use EBS snapshots, documented recovery procedures |
| Key compromise | Critical | Very Low | Key rotation, HSM for production keys |
| Manifest delivery failure | Medium | Medium | Retry logic, email fallback, delivery tracking |

## 17. Future Enhancements

### 17.1 Planned Features (Post-V1)
- Multi-org Hub support
- Advanced fingerprint matching (geo-location, hardware serials)
- License usage analytics
- API for A-Stack to query Hub status
- CLI for automated manifest generation

### 17.2 Considerations
- License pooling
- Micro-services architecture
- Kubernetes deployment
- Mobile app for operators

## 18. Appendix

### 18.1 Glossary
- CML: Customer Master License
- site.lic: Site license sub-license
- Hub: License operator component
- A-Stack License Server: Root license authority
- Fingerprint: Physical/digital site identifiers
- Usage Manifest: Monthly usage report
- Chain of Trust: Root → CML → Site

### 18.2 References
- ECDSA: https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.186-5.pdf
- JWT: https://jwt.io/
- PostgreSQL: https://www.postgresql.org/docs/

---

Document Version: 1.0  
Last Updated: January 2024  
Owner: A-Stack Development Team  
Status: Draft