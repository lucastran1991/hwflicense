# License Server - 7 Core APIs

## Overview

The License Server provides 7 core APIs for managing site keys, validating licenses, handling stats, and sending alerts.

## API Endpoints

### API 1: Create Site Key
**Endpoint:** `POST /api/v1/sites/create`

**Description:** Creates a new site key for a site ID.

**Request:**
```json
{
  "site_id": "site_abc123",
  "enterprise_id": "ent_001",
  "mode": "production",  // or "dev"
  "org_id": "org_001"
}
```

**Response:**
```json
{
  "id": "key_uuid",
  "site_id": "site_abc123",
  "enterprise_id": "ent_001",
  "key_type": "production",
  "key_value": "license-key-string",
  "issued_at": "2025-01-15T10:00:00Z",
  "expires_at": "2026-01-15T10:00:00Z",
  "status": "active"
}
```

---

### API 2: Get Site Keys
**Endpoint:** `GET /api/v1/sites`

**Description:** Lists all site keys (or filter by enterprise).

**Response:**
```json
{
  "keys": [
    {
      "id": "key_uuid",
      "site_id": "site_abc123",
      "enterprise_id": "ent_001",
      "key_type": "production",
      "key_value": "license-key-string",
      "issued_at": "2025-01-15T10:00:00Z",
      "expires_at": "2026-01-15T10:00:00Z",
      "status": "active"
    }
  ],
  "total": 1
}
```

---

### API 3: Update Site Key
**Endpoint:** `PUT /api/v1/sites/:id`

**Description:** Updates an existing site key (e.g., change status).

**Request:**
```json
{
  "status": "active"  // or "revoked"
}
```

**Response:**
```json
{
  "id": "key_uuid",
  "status": "active",
  "updated_at": "2025-01-15T10:00:00Z"
}
```

---

### API 4: Refresh Key
**Endpoint:** `POST /api/v1/keys/refresh`

**Description:** Refreshes a site key, generating a new key value.

**Request:**
```json
{
  "site_id": "site_abc123",
  "old_key": "old-license-key-string"
}
```

**Response:**
```json
{
  "id": "key_uuid",
  "site_id": "site_abc123",
  "enterprise_id": "ent_001",
  "key_type": "production",
  "key_value": "new-license-key-string",
  "issued_at": "2025-01-15T11:00:00Z",
  "expires_at": "2026-01-15T11:00:00Z",
  "status": "active"
}
```

---

### API 5: Aggregate Stats
**Endpoint:** `POST /api/v1/stats/aggregate`

**Description:** Sends quarterly usage statistics to the license server.

**Request:**
```json
{
  "period": "Q4_2025",
  "production_sites": 100,
  "dev_sites": 5,
  "user_counts": {
    "hwf_admins": 10,
    "total_users": 500
  },
  "enterprise_breakdown": [
    {
      "enterprise_id": "ent_001",
      "sites": 15,
      "users": 75
    }
  ]
}
```

**Response:**
```json
{
  "status": "saved",
  "period": "Q4_2025",
  "timestamp": "2025-01-15T10:00:00Z"
}
```

---

### API 6: Validate Key
**Endpoint:** `POST /api/v1/keys/validate`

**Description:** Validates a site key and returns JWT token if valid.

**Request:**
```json
{
  "site_id": "site_abc123",
  "key": "license-key-string"
}
```

**Response:**
```json
{
  "valid": true,
  "token": "jwt-token-string",
  "expires_in": 3600,
  "message": "License valid"
}
```

**Invalid Response:**
```json
{
  "valid": false,
  "message": "Key not found or expired"
}
```

---

### API 7: Send Alert
**Endpoint:** `POST /api/v1/alerts`

**Description:** Sends an alert to the license server (e.g., key expired).

**Request:**
```json
{
  "site_id": "site_abc123",
  "alert_type": "key_expired",  // or "key_invalid"
  "message": "Key expired on 2025-01-15",
  "timestamp": "2025-01-15T10:00:00Z"
}
```

**Response:**
```json
{
  "status": "received",
  "site_id": "site_abc123",
  "alert_type": "key_expired",
  "timestamp": "2025-01-15T10:00:00Z"
}
```

---

## Testing

### Quick Test
```bash
# Run the test script
./test_license_server_apis.sh

# Or test with custom URL
LICENSE_SERVER_URL=http://your-server:8081 ./test_license_server_apis.sh
```

### Manual Testing

```bash
# 1. Create site key
curl -X POST http://localhost:8081/api/v1/sites/create \
  -H "Content-Type: application/json" \
  -d '{
    "site_id": "test_site",
    "enterprise_id": "ent_001",
    "mode": "production",
    "org_id": "org_001"
  }'

# 2. List site keys
curl http://localhost:8081/api/v1/sites

# 3. Refresh key
curl -X POST http://localhost:8081/api/v1/keys/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "site_id": "test_site",
    "old_key": "old-key-value"
  }'

# 4. Validate key
curl -X POST http://localhost:8081/api/v1/keys/validate \
  -H "Content-Type: application/json" \
  -d '{
    "site_id": "test_site",
    "key": "license-key-value"
  }'

# 5. Send stats
curl -X POST http://localhost:8081/api/v1/stats/aggregate \
  -H "Content-Type: application/json" \
  -d '{
    "period": "Q4_2025",
    "production_sites": 100,
    "dev_sites": 5
  }'

# 6. Send alert
curl -X POST http://localhost:8081/api/v1/alerts \
  -H "Content-Type: application/json" \
  -d '{
    "site_id": "test_site",
    "alert_type": "key_expired",
    "message": "Key expired"
  }'
```

## Implementation Status

**Note:** The License Server is not yet fully implemented. The code structure exists in:
- `backend/internal/client/license_server_client.go` - Client library for these APIs
- Comments in the code indicate these 7 APIs should exist

To implement these APIs, you would need to create a separate license server microservice.

## Current Architecture

Currently, the Hub (backend) includes a client for calling these APIs, but the actual license server implementation is commented out in:
- `ecosystem.config.js` - License server is commented out
- `license-server/` - Directory exists but is empty/incomplete

The system works without the license server by handling license generation directly in the Hub.

