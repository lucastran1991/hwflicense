# License Server - License Creation Flow

## Overview

This document explains the complete license creation flow in the License Server microservice.

## Flow Diagram

```
Client Request
    ↓
API Handler (handlers.go)
    ↓
Input Validation
    ↓
Service Layer (site_service.go)
    ↓
Business Logic (Key Generation, Expiration)
    ↓
Repository Layer (repository.go)
    ↓
Database (SQLite)
    ↓
Return Created License
```

## Step-by-Step Flow

### Step 1: API Request Received

**Endpoint:** `POST /api/v1/sites/create`

**Request Body:**
```json
{
  "site_id": "site_001",
  "enterprise_id": "ent_001",
  "mode": "production",
  "org_id": "org_001"
}
```

### Step 2: API Handler Validation

**File:** `license-server/internal/api/handlers.go`

**Handler:** `CreateSiteKey()`

```go
func (h *Handlers) CreateSiteKey(c *gin.Context) {
    var req struct {
        SiteID       string `json:"site_id" binding:"required"`
        EnterpriseID string `json:"enterprise_id" binding:"required"`
        Mode         string `json:"mode" binding:"required"`
        OrgID        string `json:"org_id" binding:"required"`
    }
    
    // 1. Parse and validate JSON
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // 2. Validate mode (must be "production" or "dev")
    if req.Mode != "production" && req.Mode != "dev" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "mode must be 'production' or 'dev'"})
        return
    }
    
    // 3. Call service layer
    siteKey, err := h.siteService.CreateSiteKey(req.SiteID, req.EnterpriseID, req.Mode, req.OrgID)
    
    // 4. Return response
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, siteKey)
}
```

**Validation:**
- ✅ All required fields present
- ✅ Mode must be "production" or "dev"
- ✅ Returns 400 Bad Request for invalid input

### Step 3: Service Layer Processing

**File:** `license-server/internal/service/site_service.go`

**Method:** `CreateSiteKey()`

```go
func (s *SiteService) CreateSiteKey(siteID, enterpriseID, mode, orgID string) (*models.SiteKey, error) {
    // 1. Generate unique key value
    keyValue := generateKeyValue()  // Format: "LS-{uuid}"
    // Example: "LS-a1b2c3d4-e5f6-7890"
    
    // 2. Calculate expiration (30 days from now)
    issuedAt := time.Now()
    expiresAt := issuedAt.AddDate(0, 0, 30)  // 30 days
    
    // 3. Create SiteKey model
    siteKey := &models.SiteKey{
        ID:           uuid.New().String(),      // Unique ID
        SiteID:       siteID,                   // From request
        EnterpriseID: enterpriseID,             // From request
        KeyType:      mode,                     // "production" or "dev"
        KeyValue:     keyValue,                 // Generated "LS-{uuid}"
        IssuedAt:     issuedAt,                 // Current time
        ExpiresAt:     expiresAt,                // +30 days
        Status:       "active",                 // Always "active" on creation
        CreatedAt:    time.Now(),               // Current time
    }
    
    // 4. Save to database via repository
    if err := s.repo.CreateSiteKey(siteKey); err != nil {
        return nil, fmt.Errorf("failed to create site key: %w", err)
    }
    
    // 5. Return created key
    return siteKey, nil
}

// Helper function to generate unique key value
func generateKeyValue() string {
    return fmt.Sprintf("LS-%s", uuid.New().String())
}
```

**Key Actions:**
1. ✅ Generate unique key value using UUID
2. ✅ Set expiration to 30 days from now
3. ✅ Set status to "active"
4. ✅ Create SiteKey model with all fields
5. ✅ Pass to repository layer

### Step 4: Repository Layer - Database Insert

**File:** `license-server/internal/repository/repository.go`

**Method:** `CreateSiteKey()`

```go
func (r *Repository) CreateSiteKey(siteKey *models.SiteKey) error {
    query := `INSERT INTO site_keys (
        id, 
        site_id, 
        enterprise_id, 
        key_type, 
        key_value, 
        issued_at, 
        expires_at, 
        status, 
        created_at
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
    
    _, err := r.db.Connection.Exec(query, 
        siteKey.ID,              // Generated UUID
        siteKey.SiteID,           // From request
        siteKey.EnterpriseID,    // From request
        siteKey.KeyType,          // "production" or "dev"
        siteKey.KeyValue,         // "LS-{uuid}"
        timeToString(siteKey.IssuedAt),   // RFC3339 format
        timeToString(siteKey.ExpiresAt),  // RFC3339 format
        siteKey.Status,           // "active"
        timeToString(siteKey.CreatedAt)  // RFC3339 format
    )
    
    return err
}
```

**Database Operation:**
- ✅ INSERT into `site_keys` table
- ✅ All 9 fields inserted
- ✅ Uses prepared statement (SQL injection protection)
- ✅ Returns error if insert fails

### Step 5: Database Schema

**Table:** `site_keys`

**Schema:**
```sql
CREATE TABLE IF NOT EXISTS site_keys (
    id TEXT PRIMARY KEY,                    -- Generated UUID
    site_id TEXT UNIQUE NOT NULL,           -- From request
    enterprise_id TEXT NOT NULL,            -- From request
    key_type TEXT NOT NULL DEFAULT 'production',  -- From request
    key_value TEXT NOT NULL,                -- Generated "LS-{uuid}"
    issued_at TEXT DEFAULT (datetime('now')),    -- Current time
    expires_at TEXT NOT NULL,               -- +30 days
    status TEXT DEFAULT 'active',           -- Always "active"
    last_validated TEXT,                    -- NULL initially
    created_at TEXT DEFAULT (datetime('now'))    -- Current time
);
```

**Indexes:**
- `idx_site_keys_enterprise` - On enterprise_id
- `idx_site_keys_status` - On status
- `idx_site_keys_expires` - On expires_at

### Step 6: Response Returned

**Status Code:** `201 Created`

**Response Body:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "site_id": "site_001",
  "enterprise_id": "ent_001",
  "key_type": "production",
  "key_value": "LS-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "issued_at": "2025-10-28T16:30:00Z",
  "expires_at": "2025-11-27T16:30:00Z",
  "status": "active",
  "created_at": "2025-10-28T16:30:00Z"
}
```

## Key Features

### 1. Key Generation
- **Format:** `LS-{UUID}`
- **Example:** `LS-a1b2c3d4-e5f6-7890-abcd-ef1234567890`
- **Uniqueness:** Guaranteed by UUID
- **Location:** `generateKeyValue()` in service layer

### 2. Expiration Calculation
```go
issuedAt := time.Now()
expiresAt := issuedAt.AddDate(0, 0, 30)  // 30 days
```
- **Period:** Fixed 30 days
- **Flexible:** Can be modified in service layer
- **Automatic:** Calculated on creation

### 3. Key Types
- **production:** For production environments
- **dev:** For development/testing environments
- **Validation:** Enforced at API handler level
- **Storage:** Stored in `key_type` column

### 4. Status Management
- **Initial Status:** Always "active" on creation
- **Possible Values:** "active", "revoked"
- **Updates:** Via `PUT /api/v1/sites/:id`
- **Validation:** Checked before operations

### 5. Database Constraints
- **Unique Site ID:** `site_id TEXT UNIQUE NOT NULL`
- **Foreign Key:** Links to `enterprises` table (enterprise_id)
- **Indexes:** On enterprise_id, status, expires_at

## Error Handling

### Common Errors

**1. Missing Required Field**
```json
// Request
{
  "site_id": "site_001",
  "enterprise_id": "ent_001"
  // Missing "mode" and "org_id"
}

// Response: 400 Bad Request
{
  "error": "Key: 'Mode' Error:Field validation for 'Mode' failed on the 'required' tag"
}
```

**2. Invalid Mode**
```json
// Request
{
  "mode": "invalid"
}

// Response: 400 Bad Request
{
  "error": "mode must be 'production' or 'dev'"
}
```

**3. Database Error**
```json
// Response: 500 Internal Server Error
{
  "error": "UNIQUE constraint failed: site_keys.site_id"
}
```

## Example Request/Response

### Complete Flow

**Request:**
```bash
curl -X POST http://localhost:8082/api/v1/sites/create \
  -H "Content-Type: application/json" \
  -d '{
    "site_id": "site_12345",
    "enterprise_id": "ent_abc",
    "mode": "production",
    "org_id": "org_xyz"
  }'
```

**Response:**
```json
{
  "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "site_id": "site_12345",
  "enterprise_id": "ent_abc",
  "key_type": "production",
  "key_value": "LS-8a4f3e2b-1d5c-4b6a-9f8e-7d6c5b4a3210",
  "issued_at": "2025-10-28T16:30:39Z",
  "expires_at": "2025-11-27T16:30:39Z",
  "status": "active",
  "created_at": "2025-10-28T16:30:39Z"
}
```

## Summary

The license creation flow follows a clean 4-layer architecture:

1. **API Handler** - Validates input, handles HTTP
2. **Service Layer** - Business logic (key generation, expiration)
3. **Repository** - Database operations
4. **Database** - Persistent storage

Key features:
- ✅ Input validation at API level
- ✅ Unique key generation (UUID-based)
- ✅ 30-day expiration (configurable)
- ✅ Mode support (production/dev)
- ✅ Status management (active/revoked)
- ✅ Error handling at each layer
- ✅ SQL injection protection
- ✅ Return full created license object

The flow ensures data integrity, security, and proper error handling throughout the entire process.

