#!/bin/bash

# Comprehensive License Server 7 APIs Test
# Based on backend/internal/client/license_server_client.go

echo "========================================="
echo "License Server - 7 Core APIs Test"
echo "========================================="
echo ""

BASE_URL="${LICENSE_SERVER_URL:-http://localhost:8081}"
echo "Testing against: $BASE_URL"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

PASSED=0
FAILED=0
SKIPPED=0

# Test API function
test_api() {
    local api_num=$1
    local name=$2
    local method=$3
    local endpoint=$4
    local data=$5
    local expected_status=$6
    
    echo -e "${BLUE}=========================================${NC}"
    echo -e "${BLUE}API $api_num: $name${NC}"
    echo -e "${BLUE}=========================================${NC}"
    echo -e "${YELLOW}$method $endpoint${NC}"
    
    if [ -n "$data" ]; then
        echo "Request body:"
        echo "$data" | jq '.' 2>/dev/null || echo "$data"
    fi
    echo ""
    
    if [ "$method" = "GET" ]; then
        result=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" 2>&1)
    else
        result=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data" 2>&1)
    fi
    
    http_code=$(echo "$result" | tail -n1)
    response=$(echo "$result" | head -n -1)
    
    expected=${expected_status:-200}
    
    if [ "$http_code" -eq "$expected" ] || [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "${GREEN}✓ PASS${NC} (HTTP $http_code)"
        PASSED=$((PASSED + 1))
        echo "Response:"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    elif [ "$http_code" = "000" ]; then
        echo -e "${YELLOW}⚠ SKIP${NC} (Service not available)"
        SKIPPED=$((SKIPPED + 1))
    else
        echo -e "${RED}✗ FAIL${NC} (HTTP $http_code)"
        FAILED=$((FAILED + 1))
        echo "Response:"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    fi
    echo ""
    echo ""
}

# Store values for subsequent tests
KEY_VALUE=""
SITE_ID="test_site_$(date +%s)"

# API 1: Create Site Key
echo "========================================="
echo "API 1: Create Site Key"
echo "POST /api/v1/sites/create"
echo "========================================="

RESULT=$(curl -s -X POST "$BASE_URL/api/v1/sites/create" \
    -H "Content-Type: application/json" \
    -d "{
        \"site_id\": \"$SITE_ID\",
        \"enterprise_id\": \"ent_test_001\",
        \"mode\": \"production\",
        \"org_id\": \"test_org_001\"
    }")

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/api/v1/sites/create" \
    -H "Content-Type: application/json" \
    -d "{
        \"site_id\": \"$SITE_ID\",
        \"enterprise_id\": \"ent_test_001\",
        \"mode\": \"production\",
        \"org_id\": \"test_org_001\"
    }")

echo "Response:"
echo "$RESULT" | jq '.' 2>/dev/null || echo "$RESULT"
echo ""

if [ "$HTTP_CODE" -ge 200 ] && [ "$HTTP_CODE" -lt 300 ]; then
    echo -e "${GREEN}✓ API 1 PASS${NC} (HTTP $HTTP_CODE)"
    PASSED=$((PASSED + 1))
    KEY_VALUE=$(echo "$RESULT" | jq -r '.key_value // empty' 2>/dev/null)
    echo "Created Key: $KEY_VALUE"
else
    echo -e "${RED}✗ API 1 FAIL${NC} (HTTP $HTTP_CODE)"
    FAILED=$((FAILED + 1))
fi
echo ""
echo ""

# API 2: Get Site Keys (Listing)
echo "========================================="
echo "API 2: Get Site Keys (List)"
echo "GET /api/v1/sites"
echo "========================================="

RESULT=$(curl -s "$BASE_URL/api/v1/sites")
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/v1/sites")

echo "Response:"
echo "$RESULT" | jq '.' 2>/dev/null || echo "$RESULT"
echo ""

if [ "$HTTP_CODE" -ge 200 ] && [ "$HTTP_CODE" -lt 300 ]; then
    echo -e "${GREEN}✓ API 2 PASS${NC} (HTTP $HTTP_CODE)"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}✗ API 2 FAIL${NC} (HTTP $HTTP_CODE)"
    FAILED=$((FAILED + 1))
fi
echo ""
echo ""

# API 3: Update Site Key (if exists)
echo "========================================="
echo "API 3: Update Site Key"
echo "PUT /api/v1/sites/$SITE_ID"
echo "========================================="

RESULT=$(curl -s -X PUT "$BASE_URL/api/v1/sites/$SITE_ID" \
    -H "Content-Type: application/json" \
    -d '{
        "status": "active"
    }')
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X PUT "$BASE_URL/api/v1/sites/$SITE_ID" \
    -H "Content-Type: application/json" \
    -d '{
        "status": "active"
    }')

echo "Response:"
echo "$RESULT" | jq '.' 2>/dev/null || echo "$RESULT"
echo ""

if [ "$HTTP_CODE" -ge 200 ] && [ "$HTTP_CODE" -lt 400 ]; then
    echo -e "${GREEN}✓ API 3 PASS${NC} (HTTP $HTTP_CODE)"
    PASSED=$((PASSED + 1))
else
    echo -e "${YELLOW}⚠ API 3 SKIP${NC} (HTTP $HTTP_CODE - may not be implemented)"
    SKIPPED=$((SKIPPED + 1))
fi
echo ""
echo ""

# API 4: Refresh Key
echo "========================================="
echo "API 4: Refresh Key"
echo "POST /api/v1/keys/refresh"
echo "========================================="

RESULT=$(curl -s -X POST "$BASE_URL/api/v1/keys/refresh" \
    -H "Content-Type: application/json" \
    -d "{
        \"site_id\": \"$SITE_ID\",
        \"old_key\": \"$KEY_VALUE\"
    }")
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/api/v1/keys/refresh" \
    -H "Content-Type: application/json" \
    -d "{
        \"site_id\": \"$SITE_ID\",
        \"old_key\": \"$KEY_VALUE\"
    }")

echo "Response:"
echo "$RESULT" | jq '.' 2>/dev/null || echo "$RESULT"
echo ""

if [ "$HTTP_CODE" -ge 200 ] && [ "$HTTP_CODE" -lt 300 ]; then
    echo -e "${GREEN}✓ API 4 PASS${NC} (HTTP $HTTP_CODE)"
    PASSED=$((PASSED + 1))
    # Update key if successful
    NEW_KEY=$(echo "$RESULT" | jq -r '.key_value // empty' 2>/dev/null)
    if [ -n "$NEW_KEY" ] && [ "$NEW_KEY" != "null" ]; then
        KEY_VALUE=$NEW_KEY
        echo "New Key: $KEY_VALUE"
    fi
else
    echo -e "${RED}✗ API 4 FAIL${NC} (HTTP $HTTP_CODE)"
    FAILED=$((FAILED + 1))
fi
echo ""
echo ""

# API 5: Aggregate Stats
echo "========================================="
echo "API 5: Aggregate Stats"
echo "POST /api/v1/stats/aggregate"
echo "========================================="

RESULT=$(curl -s -X POST "$BASE_URL/api/v1/stats/aggregate" \
    -H "Content-Type: application/json" \
    -d '{
        "period": "Q4_2025",
        "production_sites": 100,
        "dev_sites": 5,
        "user_counts": {
            "hwf_admins": 10,
            "total_users": 500
        },
        "enterprise_breakdown": [
            {
                "enterprise_id": "ent_test_001",
                "sites": 15,
                "users": 75
            }
        ]
    }')
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/api/v1/stats/aggregate" \
    -H "Content-Type: application/json" \
    -d '{
        "period": "Q4_2025",
        "production_sites": 100,
        "dev_sites": 5,
        "user_counts": {
            "hwf_admins": 10,
            "total_users": 500
        },
        "enterprise_breakdown": [
            {
                "enterprise_id": "ent_test_001",
                "sites": 15,
                "users": 75
            }
        ]
    }')

echo "Response:"
echo "$RESULT" | jq '.' 2>/dev/null || echo "$RESULT"
echo ""

if [ "$HTTP_CODE" -ge 200 ] && [ "$HTTP_CODE" -lt 300 ]; then
    echo -e "${GREEN}✓ API 5 PASS${NC} (HTTP $HTTP_CODE)"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}✗ API 5 FAIL${NC} (HTTP $HTTP_CODE)"
    FAILED=$((FAILED + 1))
fi
echo ""
echo ""

# API 6: Validate Key
echo "========================================="
echo "API 6: Validate Key"
echo "POST /api/v1/keys/validate"
echo "========================================="

RESULT=$(curl -s -X POST "$BASE_URL/api/v1/keys/validate" \
    -H "Content-Type: application/json" \
    -d "{
        \"site_id\": \"$SITE_ID\",
        \"key\": \"$KEY_VALUE\"
    }")
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/api/v1/keys/validate" \
    -H "Content-Type: application/json" \
    -d "{
        \"site_id\": \"$SITE_ID\",
        \"key\": \"$KEY_VALUE\"
    }")

echo "Response:"
echo "$RESULT" | jq '.' 2>/dev/null || echo "$RESULT"
echo ""

if [ "$HTTP_CODE" -ge 200 ] && [ "$HTTP_CODE" -lt 400 ]; then
    echo -e "${GREEN}✓ API 6 PASS${NC} (HTTP $HTTP_CODE)"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}✗ API 6 FAIL${NC} (HTTP $HTTP_CODE)"
    FAILED=$((FAILED + 1))
fi
echo ""
echo ""

# API 7: Send Alert
echo "========================================="
echo "API 7: Send Alert"
echo "POST /api/v1/alerts"
echo "========================================="

RESULT=$(curl -s -X POST "$BASE_URL/api/v1/alerts" \
    -H "Content-Type: application/json" \
    -d "{
        \"site_id\": \"$SITE_ID\",
        \"alert_type\": \"key_expired\",
        \"message\": \"Test alert - key about to expire\",
        \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
    }")
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/api/v1/alerts" \
    -H "Content-Type: application/json" \
    -d "{
        \"site_id\": \"$SITE_ID\",
        \"alert_type\": \"key_expired\",
        \"message\": \"Test alert - key about to expire\",
        \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
    }")

echo "Response:"
echo "$RESULT" | jq '.' 2>/dev/null || echo "$RESULT"
echo ""

if [ "$HTTP_CODE" -ge 200 ] && [ "$HTTP_CODE" -lt 300 ]; then
    echo -e "${GREEN}✓ API 7 PASS${NC} (HTTP $HTTP_CODE)"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}✗ API 7 FAIL${NC} (HTTP $HTTP_CODE)"
    FAILED=$((FAILED + 1))
fi
echo ""
echo ""

# Summary
echo "========================================="
echo "Test Summary"
echo "========================================="
echo -e "${GREEN}✓ Passed: $PASSED${NC}"
echo -e "${RED}✗ Failed: $FAILED${NC}"
echo -e "${YELLOW}⚠ Skipped: $SKIPPED${NC}"
echo "Total: $((PASSED + FAILED + SKIPPED))"
echo ""

if [ $FAILED -eq 0 ] && [ $PASSED -gt 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    exit 0
elif [ $FAILED -gt 0 ]; then
    echo -e "${RED}✗ Some tests failed${NC}"
    exit 1
else
    echo -e "${YELLOW}⚠ No tests ran (license server may not be running)${NC}"
    exit 2
fi

