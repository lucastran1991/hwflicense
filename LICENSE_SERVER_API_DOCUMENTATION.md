# License Server API Documentation

## Overview

License Server is a microservice implementing 7 core APIs for hierarchical license management based on Q&A meeting requirements (Oct 27, 2025).

**Base URL:** `http://localhost:8081`  
**Server:** Running on port 8081

---

## API Endpoints

### 1. Create Site Key

**Endpoint:** `POST /api/v1/sites/create`

**Description:** Creates a new site license key with dev/production distinction.

**Request Body:**
```json
{
  "site_id": "site_abc123",
  "enterprise_id": "ent_001",
  "mode": "production",
  "org_id": "veolia_hub"
}
```

**Parameters:**
- `site_id` (required): Unique site identifier
- `enterprise_id` (required): Enterprise identifier  
- `mode` (required): `"production"` or `"dev"` - Key type
- `org_id` (required): Organization identifier

**Response (201 Created):**
```json
{
  "id": "key_uuid",
  "site_id": "site_abc123",
  "enterprise_id": "ent_001",
  "key_type": "production",
  "key_value": "base64_encoded_key",
  "issued_at": "2025-10-27T15:00:00Z",
  "expires_at": "2025-11-26T15:00:00Z",
  "status": "active"
}
```

**Key Type Rules:**
- HWF sites: Always `"production"` (automatic)
- Boost sites: Configurable `"dev"` or `"production"`
- Dev keys: Not billed, for configuration/testing
- Production keys: Billed, for active sites

---

### 2. Update Site Key

**Endpoint:** `PUT /api/v1/sites/:site_id`

**Description:** Updates site key or transitions dev ↔ production.

**Request Body:**
```json
{
  "key_type": "production",
  "reason": "site_ready_for_production"
}
```

**Parameters:**
- `key_type` (required): `"production"` or `"dev"` - New key type
- `reason` (optional): Reason for transition

**Response (200 OK):**
```json
{
  "site_id": "site_abc123",
  "key_type": "production",
  "key_value": "new_base64_key",
  "issued_at": "2025-10-27T15:00:00Z",
  "expires_at": "2025-11-26T15:00:00Z",
  "status": "active"
}
```

**Behavior:**
- If key type changes: Revokes old key, generates new key with new type
- If same type: Returns current key
- Logs transition in `key_refresh_log`

---

### 3. Delete Site Key

**Endpoint:** `DELETE /api/v1/sites/:site_id`

**Description:** Revokes a site key immediately.

**Response (200 OK):**
```json
{
  "message": "Site key revoked successfully"
}
```

**Behavior:**
- Sets status to `"revoked"`
- Cannot be recovered
- Archived in database

---

### 4. Refresh Key (Monthly)

**Endpoint:** `POST /api/v1/keys/refresh`

**Description:** Monthly key refresh for security (mandatory every 30 days).

**Request Body:**
```json
{
  "site_id": "site_abc123",
  "old_key": "current_key_token"
}
```

**Parameters:**
- `site_id` (required): Site identifier
- `old_key` (required): Current key to be refreshed

**Response (200 OK):**
```json
{
  "site_id": "site_abc123",
  "key_value": "new_refreshed_key",
  "key_type": "production",
  "issued_at": "2025-11-27T15:00:00Z",
  "expires_at": "2025-12-27T15:00:00Z",
  "status": "active",
  "message": "Key refreshed successfully"
}
```

**Behavior:**
- Validates old key matches
- Generates new key (same type)
- Invalidates old key immediately
- Sets new expiration = now + 30 days
- Logs refresh in `key_refresh_log`

**Security:**
- Keys expire after 30 days
- Monthly refresh mandatory (Alex's requirement)
- Old keys cannot be reused

---

### 5. Get Aggregate Stats (Quarterly)

**Endpoint:** `POST /api/v1/stats/aggregate`

**Description:** Receives quarterly stats from HWF for billing.

**Request Body:**
```json
{
  "period": "Q4_2025",
  "production_sites": 100,
  "dev_sites": 5,
  "user_counts": {
    "hwf_admins": [
      {"name": "John Doe", "email": "john@example.com"}
    ],
    "enterprise_admins": 25,
    "plant_users": 150,
    "demo_users": 10
  },
  "enterprise_breakdown": [
    {
      "enterprise_id": "ent_001",
      "production_sites": 50,
      "dev_sites": 2
    }
  ]
}
```

**Parameters:**
- `period` (required): Quarter identifier (e.g., `"Q4_2025"`)
- `production_sites` (required): Count of production sites
- `dev_sites` (required): Count of dev sites
- `user_counts` (required): User statistics by role
- `enterprise_breakdown` (required): Breakdown by enterprise

**Response (200 OK):**
```json
{
  "period": "Q4_2025",
  "production_sites": 100,
  "dev_sites": 5,
  "message": "Stats saved successfully"
}
```

**Privacy Rules:**
- HWF admins: Full info (names + emails)
- Other roles: Counts only
- Enterprise names: NOT included
- Site names: NOT included

---

### 6. Check Validity

**Endpoint:** `POST /api/v1/keys/validate`

**Description:** Validates site key and returns JWT token for caching.

**Request Body:**
```json
{
  "site_id": "site_abc123",
  "key": "current_key_token"
}
```

**Parameters:**
- `site_id` (required): Site identifier
- `key` (required): Key to validate

**Response (200 OK - Valid):**
```json
{
  "valid": true,
  "token": "jwt_token_here",
  "expires_in": 2592000,
  "message": "Key is valid"
}
```

**Response (401 Unauthorized - Invalid):**
```json
{
  "valid": false,
  "message": "key expired"
}
```

**Validation Checks:**
1. Key signature valid (ECDSA)
2. Not expired (< 30 days old)
3. Not revoked
4. Site exists and active

**Token Caching:**
- Token valid for 30 days (2592000 seconds)
- Stored in `validation_cache` table
- Client should cache and reuse
- Reduces call home frequency

---

### 7. Send Alerts

**Endpoint:** `POST /api/v1/alerts`

**Description:** Receives alerts from HWF when keys are invalid.

**Request Body:**
```json
{
  "site_id": "site_abc123",
  "alert_type": "key_expired",
  "message": "Key expired on 2025-11-27",
  "timestamp": "2025-11-27T15:00:00Z"
}
```

**Parameters:**
- `site_id` (required): Site identifier
- `alert_type` (required): `"key_expired"` or `"key_invalid"`
- `message` (required): Alert message
- `timestamp` (required): Alert timestamp

**Response (200 OK):**
```json
{
  "id": "alert_uuid",
  "site_id": "site_abc123",
  "alert_type": "key_expired",
  "message": "Key expired on 2025-11-27",
  "alert_timestamp": "2025-11-27T15:00:00Z",
  "sent_to_astack": false
}
```

**Behavior:**
- Stores in `alerts` table
- Logs for monitoring
- Can trigger email notification (optional)

---

## Health Check

**Endpoint:** `GET /health`

**Description:** Simple health check endpoint.

**Response (200 OK):**
```json
{
  "status": "ok",
  "service": "license-server"
}
```

---

## Error Handling

### Error Response Format
```json
{
  "error": "Error message here"
}
```

### Common Errors

**400 Bad Request:**
- Missing required fields
- Invalid JSON
- Invalid parameter values

**401 Unauthorized:**
- Invalid key
- Expired key
- Revoked key

**404 Not Found:**
- Site not found
- Enterprise not found

**500 Internal Server Error:**
- Database errors
- Unexpected errors

---

## Authentication

Currently, the License Server APIs do not require authentication tokens. This should be added for production use:

```bash
# Recommended: Add API key or JWT authentication
Authorization: Bearer <token>
```

---

## Database Schema

### enterprises
- `id`, `name`, `org_id`, `enterprise_key`, `created_at`

### site_keys
- `id`, `site_id`, `enterprise_id`, `key_type`, `key_value`, `issued_at`, `expires_at`, `status`, `last_validated`, `created_at`

### key_refresh_log
- `id`, `site_id`, `old_key`, `new_key`, `refreshed_at`, `reason`

### quarterly_stats
- `id`, `period`, `production_sites`, `dev_sites`, `user_counts`, `enterprise_breakdown`, `created_at`

### validation_cache
- `id`, `site_id`, `token`, `expires_at`, `created_at`

### alerts
- `id`, `site_id`, `alert_type`, `message`, `alert_timestamp`, `sent_to_astack`, `created_at`

---

## Key Lifecycle

1. **Creation:** Key generated when site created (expires in 30 days)
2. **Validation:** Periodic check (cached token, refreshed as needed)
3. **Refresh:** Monthly mandatory refresh (before expiration)
4. **Expiration:** Key expires after 30 days
5. **Invalidation:** Key revoked if expired or deleted

---

## Running License Server

### Development
```bash
cd license-server
go run cmd/license-server/main.go
```

### Production
```bash
./scripts/license-server.sh start
```

### Management
```bash
# Start
./scripts/license-server.sh start

# Stop
./scripts/license-server.sh stop

# Restart
./scripts/license-server.sh restart

# Status
./scripts/license-server.sh status

# Logs
./scripts/license-server.sh logs

# Build
./scripts/license-server.sh build
```

---

## Environment Variables

Create `.env` in `license-server/`:

```env
PORT=8081
DATABASE_PATH=./data/license_server.db
JWT_SECRET=license-server-secret-key-change-in-production
ENVIRONMENT=development
```

---

## Next Steps

1. Add JWT authentication
2. Add rate limiting
3. Add metrics/monitoring
4. Add request validation middleware
5. Add comprehensive error handling

