#!/bin/bash
# Complete License Server API Test Suite

set -e

BASE_URL="${1:-http://localhost:8081}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "========================================="
echo "License Server API Test Suite"
echo "========================================="
echo "Base URL: $BASE_URL"
echo ""

# Test variables
SITE_ID="test_site_$(date +%s)"
ENTERPRISE_ID="test_enterprise_001"
ORG_ID="test_org_001"
MODE="production"
CREATED_KEY=""

# Function to print test result
print_result() {
    local status=$1
    local test_name=$2
    if [ $status -eq 0 ]; then
        echo -e "${GREEN}✓ $test_name${NC}"
    else
        echo -e "${RED}✗ $test_name${NC}"
    fi
}

# Test 1: Health Check
echo -e "${BLUE}Test 1: Health Check${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/health")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)
echo "Response: $body"
if [ "$http_code" = "200" ]; then
    print_result 0 "Health check passed"
else
    print_result 1 "Health check failed (HTTP $http_code)"
fi
echo ""

# Test 2: Create Site Key
echo -e "${BLUE}Test 2: Create Site Key${NC}"
payload="{\"site_id\":\"$SITE_ID\",\"enterprise_id\":\"$ENTERPRISE_ID\",\"mode\":\"$MODE\",\"org_id\":\"$ORG_ID\"}"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/v1/sites/create" \
    -H "Content-Type: application/json" \
    -d "$payload")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)
echo "Response: $body"

if [ "$http_code" = "201" ] || [ "$http_code" = "200" ]; then
    # Extract key value if present
    CREATED_KEY=$(echo "$body" | grep -o '"key_value":"[^"]*"' | cut -d'"' -f4 || echo "")
    print_result 0 "Create site key passed"
    echo "Site ID: $SITE_ID"
    [ -n "$CREATED_KEY" ] && echo "Created Key: $CREATED_KEY"
else
    print_result 1 "Create site key failed (HTTP $http_code)"
fi
echo ""

# Test 3: List Site Keys
echo -e "${BLUE}Test 3: List Site Keys${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/v1/sites")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)
echo "Response: $body"

if [ "$http_code" = "200" ]; then
    print_result 0 "List site keys passed"
    key_count=$(echo "$body" | grep -o '"site_id"' | wc -l | tr -d ' ')
    echo "Total keys found: $key_count"
else
    print_result 1 "List site keys failed (HTTP $http_code)"
fi
echo ""

# Test 4: Update Site Key Status
echo -e "${BLUE}Test 4: Update Site Key Status${NC}"
payload="{\"status\":\"active\"}"
response=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_URL/api/v1/sites/$SITE_ID" \
    -H "Content-Type: application/json" \
    -d "$payload")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)
echo "Response: $body"

if [ "$http_code" = "200" ]; then
    print_result 0 "Update site key passed"
else
    print_result 1 "Update site key failed (HTTP $http_code)"
fi
echo ""

# Test 5: Validate Key (if we have a created key)
if [ -n "$CREATED_KEY" ]; then
    echo -e "${BLUE}Test 5: Validate Key${NC}"
    payload="{\"site_id\":\"$SITE_ID\",\"key\":\"$CREATED_KEY\"}"
    response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/v1/keys/validate" \
        -H "Content-Type: application/json" \
        -d "$payload")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)
    echo "Response: $body"

    if [ "$http_code" = "200" ]; then
        print_result 0 "Validate key passed"
        is_valid=$(echo "$body" | grep -o '"valid":true' || echo "")
        if [ -n "$is_valid" ]; then
            echo "Key is valid ✓"
        else
            echo "Key validation returned invalid"
        fi
    else
        print_result 1 "Validate key failed (HTTP $http_code)"
    fi
else
    echo -e "${YELLOW}Test 5: Validate Key - Skipped (no key created)${NC}"
fi
echo ""

# Test 6: Aggregate Stats
echo -e "${BLUE}Test 6: Aggregate Stats${NC}"
payload="{\"period\":\"Q1_2025\",\"production_sites\":100,\"dev_sites\":5,\"user_counts\":{\"total\":150},\"enterprise_breakdown\":[{\"enterprise_id\":\"ent001\",\"sites\":10}]}"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/v1/stats/aggregate" \
    -H "Content-Type: application/json" \
    -d "$payload")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)
echo "Response: $body"

if [ "$http_code" = "200" ]; then
    print_result 0 "Aggregate stats passed"
else
    print_result 1 "Aggregate stats failed (HTTP $http_code)"
fi
echo ""

# Test 7: Send Alert
echo -e "${BLUE}Test 7: Send Alert${NC}"
current_time=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
payload="{\"site_id\":\"$SITE_ID\",\"alert_type\":\"key_expired\",\"message\":\"Key expired at $current_time\",\"timestamp\":\"$current_time\",\"sent_to_astack\":0}"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/v1/alerts" \
    -H "Content-Type: application/json" \
    -d "$payload")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)
echo "Response: $body"

if [ "$http_code" = "200" ]; then
    print_result 0 "Send alert passed"
else
    print_result 1 "Send alert failed (HTTP $http_code)"
fi
echo ""

echo "========================================="
echo "Test Suite Complete"
echo "========================================="
echo "Test Site ID: $SITE_ID"
echo "Test Enterprise ID: $ENTERPRISE_ID"
echo ""

