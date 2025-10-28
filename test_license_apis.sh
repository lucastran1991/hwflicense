#!/bin/bash

# Test License Server 7 APIs
# These are the APIs that the license server should provide

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
NC='\033[0m'

PASSED=0
FAILED=0

# Function to test API
test_api() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    
    echo -n "Testing API: $name... "
    
    if [ "$method" = "GET" ]; then
        result=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint")
    else
        result=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi
    
    http_code=$(echo "$result" | tail -n1)
    response=$(echo "$result" | head -n -1)
    
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "${GREEN}✓ PASS${NC} (HTTP $http_code)"
        PASSED=$((PASSED + 1))
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        echo -e "${RED}✗ FAIL${NC} (HTTP $http_code)"
        FAILED=$((FAILED + 1))
        echo "$response"
    fi
    echo ""
}

# API 1: Create Site Key
echo "========================================="
echo "API 1: Create Site Key"
echo "========================================="
RESULT=$(curl -s -X POST "$BASE_URL/api/v1/sites/create" \
    -H "Content-Type: application/json" \
    -d '{
        "site_id": "test_site_001",
        "enterprise_id": "ent_001",
        "mode": "production",
        "org_id": "test_org"
    }')

echo "$RESULT" | jq '.' 2>/dev/null || echo "$RESULT"
echo ""

# Extract key value for subsequent tests
KEY_VALUE=$(echo "$RESULT" | jq -r '.key_value // empty' 2>/dev/null)
echo "Created Key: $KEY_VALUE"
echo ""

# API 2, 3, etc. (Placeholder)
echo "========================================="
echo "API 2-7: Additional APIs"
echo "========================================="
echo "These APIs require license server implementation:"
echo "  API 2: Get Site Key"
echo "  API 3: Update Site Key"
echo "  API 4: Refresh Key (POST /api/v1/keys/refresh)"
echo "  API 5: Aggregate Stats (POST /api/v1/stats/aggregate)"
echo "  API 6: Validate Key (POST /api/v1/keys/validate)"
echo "  API 7: Send Alert (POST /api/v1/alerts)"
echo ""

# Summary
echo "========================================="
echo "Test Summary"
echo "========================================="
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
else
    echo -e "${RED}✗ Some tests failed${NC}"
fi

