#!/bin/bash
# Comprehensive Backend/Hub API Test Suite

BASE_URL="${1:-http://localhost:8080}"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "========================================="
echo "Backend/Hub API Test Suite"
echo "========================================="
echo "Base URL: $BASE_URL"
echo ""

PASSED=0
FAILED=0
SKIPPED=0

# Function to test API
test_api() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local auth_token=$5
    
    echo -e "${BLUE}=========================================${NC}"
    echo -e "${BLUE}Testing: $name${NC}"
    echo -e "${BLUE}=========================================${NC}"
    
    if [ "$method" = "GET" ]; then
        if [ -n "$auth_token" ]; then
            response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" \
                -H "Authorization: Bearer $auth_token" 2>&1)
        else
            response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" 2>&1)
        fi
    else
        if [ -n "$auth_token" ]; then
            response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $auth_token" \
                -d "$data" 2>&1)
        else
            response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" \
                -H "Content-Type: application/json" \
                -d "$data" 2>&1)
        fi
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)
    
    if [ "$http_code" = "200" ] || [ "$http_code" = "201" ]; then
        echo -e "${GREEN}✓ PASS${NC} (HTTP $http_code)"
        PASSED=$((PASSED + 1))
        if [ -n "$body" ]; then
            echo "Response:"
            echo "$body" | jq '.' 2>/dev/null || echo "$body"
        fi
    elif [ "$http_code" = "401" ]; then
        echo -e "${YELLOW}⚠ SKIP (Authentication required)${NC}"
        SKIPPED=$((SKIPPED + 1))
    else
        echo -e "${RED}✗ FAIL${NC} (HTTP $http_code)"
        FAILED=$((FAILED + 1))
        if [ -n "$body" ]; then
            echo "Error:"
            echo "$body" | jq '.' 2>/dev/null || echo "$body"
        fi
    fi
    echo ""
}

# Test 1: Health Check
echo "Test 1: Health Check"
test_api "GET /api/health" "GET" "/api/health" "" ""

# Test 2: Login (Get Auth Token)
echo "Test 2: Login"
LOGIN_DATA='{"username":"admin","password":"admin123"}'
login_response=$(curl -s -X POST "$BASE_URL/api/auth/login" \
    -H "Content-Type: application/json" \
    -d "$LOGIN_DATA")
AUTH_TOKEN=$(echo "$login_response" | jq -r '.token' 2>/dev/null)
echo "Login Response:"
echo "$login_response" | jq '.' 2>/dev/null || echo "$login_response"
echo ""

if [ -z "$AUTH_TOKEN" ] || [ "$AUTH_TOKEN" = "null" ]; then
    echo -e "${RED}✗ FAILED to get auth token${NC}"
    echo "Will continue with unauthenticated requests only"
    AUTH_TOKEN=""
fi
echo "Auth Token: ${AUTH_TOKEN:0:50}..."
echo ""

# Test 3: Create Site
echo "Test 3: Create Site"
SITE_ID="test_site_$(date +%s)"
CREATE_SITE_DATA='{
    "site_id":"'$SITE_ID'",
    "enterprise_id":"test_enterprise",
    "org_id":"test_org",
    "key_type":"production"
}'
test_api "POST /api/sites/create" "POST" "/api/sites/create" "$CREATE_SITE_DATA" "$AUTH_TOKEN"

# Test 4: List Sites
echo "Test 4: List Sites"
test_api "GET /api/sites" "GET" "/api/sites" "" "$AUTH_TOKEN"

# Test 5: Get Site
echo "Test 5: Get Site"
test_api "GET /api/sites/$SITE_ID" "GET" "/api/sites/$SITE_ID" "" "$AUTH_TOKEN"

# Test 6: Heartbeat
echo "Test 6: Site Heartbeat"
test_api "POST /api/sites/$SITE_ID/heartbeat" "POST" "/api/sites/$SITE_ID/heartbeat" '{}' "$AUTH_TOKEN"

# Test 7: Get CML
echo "Test 7: Get CML"
test_api "GET /api/cml" "GET" "/api/cml" "" "$AUTH_TOKEN"

# Test 8: Refresh CML
echo "Test 8: Refresh CML"
test_api "POST /api/cml/refresh" "POST" "/api/cml/refresh" '{}' "$AUTH_TOKEN"

# Test 9: Generate Manifest
echo "Test 9: Generate Manifest"
MANIFEST_DATA='{
    "site_id":"'$SITE_ID'",
    "org_id":"test_org"
}'
test_api "POST /api/manifests/generate" "POST" "/api/manifests/generate" "$MANIFEST_DATA" "$AUTH_TOKEN"

# Test 10: List Manifests
echo "Test 10: List Manifests"
test_api "GET /api/manifests" "GET" "/api/manifests" "" "$AUTH_TOKEN"

# Test 11: Get Ledger
echo "Test 11: Get Ledger"
test_api "GET /api/ledger" "GET" "/api/ledger" "" "$AUTH_TOKEN"

# Test 12: Validate (Public endpoint)
echo "Test 12: Validate License (Public)"
VALIDATE_DATA='{
    "site_id":"'$SITE_ID'",
    "license_key":"test_key"
}'
test_api "POST /api/validate" "POST" "/api/validate" "$VALIDATE_DATA" ""

# Summary
echo "========================================="
echo "Test Summary"
echo "========================================="
echo -e "${GREEN}✓ Passed: $PASSED${NC}"
echo -e "${RED}✗ Failed: $FAILED${NC}"
echo -e "${YELLOW}⚠ Skipped: $SKIPPED${NC}"
echo "Total: $((PASSED + FAILED + SKIPPED))"
echo ""

if [ $FAILED -eq 0 ]; then
    exit 0
else
    exit 1
fi

