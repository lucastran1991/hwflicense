#!/bin/bash

# Comprehensive API Testing Script
# Tests both License Server and Hub APIs

set -e

echo "========================================="
echo "License Server & Hub API Testing"
echo "========================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Test counters
PASSED=0
FAILED=0

test_api() {
    local name=$1
    local endpoint=$2
    local method=$3
    local data=$4
    
    echo -n "Testing $name... "
    
    if [ "$method" = "GET" ]; then
        result=$(curl -s -w "\n%{http_code}" -X $method "http://localhost:8081$endpoint")
    elif [ "$method" = "POST" ]; then
        result=$(curl -s -w "\n%{http_code}" -X $method "http://localhost:8081$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    elif [ "$method" = "PUT" ]; then
        result=$(curl -s -w "\n%{http_code}" -X $method "http://localhost:8081$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    elif [ "$method" = "DELETE" ]; then
        result=$(curl -s -w "\n%{http_code}" -X $method "http://localhost:8081$endpoint")
    fi
    
    # Extract HTTP code
    http_code=$(echo "$result" | tail -n1)
    
    # Check if HTTP code is 2xx
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "${GREEN}✓ PASS${NC} (HTTP $http_code)"
        PASSED=$((PASSED + 1))
        echo "$result" | head -n -1 | jq '.' 2>/dev/null || echo "$result" | head -n -1
    else
        echo -e "${RED}✗ FAIL${NC} (HTTP $http_code)"
        FAILED=$((FAILED + 1))
        echo "$result" | head -n -1
    fi
    echo ""
}

echo "========================================="
echo "1. Testing License Server APIs"
echo "========================================="
echo ""

# Test 1: Health Check
test_api "Health Check" "/health" "GET" ""

# Test 2: Create Site Key (API 1)
echo "--- Creating test site key ---"
RESULT=$(curl -s -X POST "http://localhost:8081/api/v1/sites/create" \
    -H "Content-Type: application/json" \
    -d '{"site_id":"test_api_001","enterprise_id":"ent_test_001","mode":"production","org_id":"test_org"}')
echo "$RESULT" | jq '.'

KEY_VALUE=$(echo "$RESULT" | jq -r '.key_value')
echo "Created Key: $KEY_VALUE"
echo ""

if [ "$KEY_VALUE" != "null" ] && [ -n "$KEY_VALUE" ]; then
    PASSED=$((PASSED + 1))
    echo -e "${GREEN}✓ Site key created successfully${NC}"
else
    FAILED=$((FAILED + 1))
    echo -e "${RED}✗ Site key creation failed${NC}"
fi
echo ""

# Test 3: Get Aggregate Stats (API 5)
echo "--- Testing Stats Aggregate ---"
RESULT=$(curl -s -X POST "http://localhost:8081/api/v1/stats/aggregate" \
    -H "Content-Type: application/json" \
    -d '{"period":"Q4_2025","production_sites":100,"dev_sites":5,"user_counts":{"hwf_admins":[]},"enterprise_breakdown":[]}')
echo "$RESULT" | jq '.'

if echo "$RESULT" | jq -e '.period' > /dev/null 2>&1; then
    PASSED=$((PASSED + 1))
    echo -e "${GREEN}✓ Stats saved successfully${NC}"
else
    FAILED=$((FAILED + 1))
    echo -e "${RED}✗ Stats save failed${NC}"
fi
echo ""

# Test 4: Validate Key (API 6)
echo "--- Testing Key Validation ---"
RESULT=$(curl -s -X POST "http://localhost:8081/api/v1/keys/validate" \
    -H "Content-Type: application/json" \
    -d "{\"site_id\":\"test_api_001\",\"key\":\"test_invalid_key\"}")
echo "$RESULT" | jq '.'

if echo "$RESULT" | jq -e '.valid' > /dev/null 2>&1; then
    PASSED=$((PASSED + 1))
    echo -e "${GREEN}✓ Validation working (expected: invalid key)${NC}"
else
    FAILED=$((FAILED + 1))
    echo -e "${RED}✗ Validation failed${NC}"
fi
echo ""

# Test 5: Refresh Key (API 4)
echo "--- Testing Key Refresh ---"
if [ "$KEY_VALUE" != "null" ] && [ -n "$KEY_VALUE" ]; then
    RESULT=$(curl -s -X POST "http://localhost:8081/api/v1/keys/refresh" \
        -H "Content-Type: application/json" \
        -d "{\"site_id\":\"test_api_001\",\"old_key\":\"$KEY_VALUE\"}")
    echo "$RESULT" | jq '.'

    if echo "$RESULT" | jq -e '.key_value' > /dev/null 2>&1; then
        PASSED=$((PASSED + 1))
        echo -e "${GREEN}✓ Key refresh successful${NC}"
        NEW_KEY=$(echo "$RESULT" | jq -r '.key_value')
        echo "New Key: $NEW_KEY"
    else
        FAILED=$((FAILED + 1))
        echo -e "${RED}✗ Key refresh failed${NC}"
    fi
else
    echo -e "${YELLOW}⚠ Skipping refresh test (no valid key)${NC}"
fi
echo ""

# Test 6: Alert (API 7)
echo "--- Testing Alert ---"
RESULT=$(curl -s -X POST "http://localhost:8081/api/v1/alerts" \
    -H "Content-Type: application/json" \
    -d '{"site_id":"test_api_001","alert_type":"key_expired","message":"Test alert","timestamp":"2025-12-28T12:00:00Z"}')
echo "$RESULT" | jq '.'

if echo "$RESULT" | jq -e '.site_id' > /dev/null 2>&1; then
    PASSED=$((PASSED + 1))
    echo -e "${GREEN}✓ Alert received successfully${NC}"
else
    FAILED=$((FAILED + 1))
    echo -e "${RED}✗ Alert failed${NC}"
fi
echo ""

# Test 7: Hub Health Check
echo "========================================="
echo "2. Testing Hub APIs"
echo "========================================="
echo ""

RESULT=$(curl -s http://localhost:8080/api/health)
if echo "$RESULT" | jq -e '.status' > /dev/null 2>&1; then
    PASSED=$((PASSED + 1))
    echo -e "${GREEN}✓ Hub health check passed${NC}"
    echo "$RESULT"
else
    FAILED=$((FAILED + 1))
    echo -e "${RED}✗ Hub health check failed${NC}"
fi
echo ""

# Summary
echo "========================================="
echo "Test Summary"
echo "========================================="
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo "Total: $((PASSED + FAILED))"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some tests failed${NC}"
    exit 1
fi

