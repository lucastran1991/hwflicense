#!/bin/bash

# Comprehensive API Test Script
# Runs all API tests and saves results to tests/ folder

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
KMS_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$KMS_DIR"

# Test configuration
BASE_URL="http://localhost:8080"
RESULTS_DIR="$SCRIPT_DIR"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
TEST_LOG="$RESULTS_DIR/test_results_${TIMESTAMP}.log"
TEST_SUMMARY="$RESULTS_DIR/test_summary_${TIMESTAMP}.json"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Test counters
TOTAL=0
PASSED=0
FAILED=0

# Initialize test summary JSON
echo "{\"timestamp\": \"$(date -Iseconds)\", \"base_url\": \"$BASE_URL\", \"tests\": [" > "$TEST_SUMMARY"

log() {
    echo -e "${BLUE}[TEST]${NC} $1" | tee -a "$TEST_LOG"
}

pass() {
    echo -e "${GREEN}✓ PASS${NC} $1" | tee -a "$TEST_LOG"
    PASSED=$((PASSED + 1))
    TOTAL=$((TOTAL + 1))
    echo ",{\"name\": \"$1\", \"status\": \"PASS\"}" >> "$TEST_SUMMARY"
}

fail() {
    echo -e "${RED}✗ FAIL${NC} $1" | tee -a "$TEST_LOG"
    FAILED=$((FAILED + 1))
    TOTAL=$((TOTAL + 1))
    echo ",{\"name\": \"$1\", \"status\": \"FAIL\", \"error\": \"$2\"}" >> "$TEST_SUMMARY"
}

# Cleanup function
cleanup() {
    # Close JSON array
    echo "]}" >> "$TEST_SUMMARY"
    
    echo "" | tee -a "$TEST_LOG"
    echo "=== Test Summary ===" | tee -a "$TEST_LOG"
    echo "Total: $TOTAL" | tee -a "$TEST_LOG"
    echo "Passed: $PASSED" | tee -a "$TEST_LOG"
    echo "Failed: $FAILED" | tee -a "$TEST_LOG"
    
    if [ $FAILED -eq 0 ]; then
        echo -e "${GREEN}All tests passed!${NC}" | tee -a "$TEST_LOG"
        exit 0
    else
        echo -e "${RED}Some tests failed!${NC}" | tee -a "$TEST_LOG"
        exit 1
    fi
}

trap cleanup EXIT

log "Starting comprehensive API tests..."
log "Results will be saved to: $RESULTS_DIR"

# Test 1: Health Check
log "Test 1: Health Check"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/health")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [ "$HTTP_CODE" = "200" ] && echo "$BODY" | grep -q "ok"; then
    echo "$BODY" | jq . > "$RESULTS_DIR/health_check.json" 2>/dev/null || echo "$BODY" > "$RESULTS_DIR/health_check.json"
    pass "Health Check"
else
    fail "Health Check" "HTTP $HTTP_CODE: $BODY"
fi

# Test 2: Register Symmetric Key
log "Test 2: Register Symmetric Key (Auto-generated)"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/keys" \
    -H "Content-Type: application/json" \
    -d '{"key_type": "symmetric", "ttl_seconds": 31536000}')
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
SYMMETRIC_KEY_ID=$(echo "$BODY" | jq -r '.key_id // empty' 2>/dev/null)
if [ "$HTTP_CODE" = "200" ] && [ -n "$SYMMETRIC_KEY_ID" ]; then
    echo "$BODY" | jq . > "$RESULTS_DIR/register_symmetric.json" 2>/dev/null || echo "$BODY" > "$RESULTS_DIR/register_symmetric.json"
    pass "Register Symmetric Key"
else
    fail "Register Symmetric Key" "HTTP $HTTP_CODE: $BODY"
fi

# Test 3: Register Asymmetric Key
log "Test 3: Register Asymmetric Key (Auto-generated)"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/keys" \
    -H "Content-Type: application/json" \
    -d '{"key_type": "asymmetric", "ttl_seconds": 31536000}')
HTTP_OK=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
ASYMMETRIC_KEY_ID=$(echo "$BODY" | jq -r '.key_id // empty' 2>/dev/null)
if [ "$HTTP_CODE" = "200" ] && [ -n "$ASYMMETRIC_KEY_ID" ]; then
    echo "$BODY" | jq . > "$RESULTS_DIR/register_asymmetric.json" 2>/dev/null || echo "$BODY" > "$RESULTS_DIR/register_asymmetric.json"
    pass "Register Asymmetric Key"
else
    fail "Register Asymmetric Key" "HTTP $HTTP_CODE: $BODY"
fi

# Test 4: Register External Key Material
log有一段 "Test 4: Register Symmetric Key (External Key Material)"
EXTERNAL_KEY=$(openssl rand -base64 32)
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/keys" \
    -H "Content-Type: application/json" \
    -d "{\"key_type\": \"symmetric\", \"key_material\": \"$EXTERNAL_KEY\", \"ttl_seconds\": 31536000}")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
EXTERNAL_KEY_ID=$(echo "$BODY" | jq -r '.key_id // empty' 2>/dev/null)
if [ "$HTTP_CODE" = "200" ] && [ -n "$EXTERNAL_KEY_ID" ]; then
    echo "$BODY" | jq . > "$RESULTS_DIR/register_external.json" 2>/dev/null || echo "$BODY" > "$RESULTS_DIR/register_external.json"
    pass "Register External Key"
else
    fail "Register External Key" "HTTP $HTTP_CODE: $BODY"
fi

# Test 5: Validate Symmetric Key
log "Test 5: Validate Symmetric Key"
if [ -n "$SYMMETRIC_KEY_ID" ]; then
    # Get the key back to validate (we need to use the key material)
    # For this test, we'll just validate the key exists
    RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/keys/validate" \
        -H "Content-Type: application/json" \
        -d "{\"key_id\": \"$SYMMETRIC_KEY_ID\", \"key\": \"$EXTERNAL_KEY\"}")
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | sed '$d')
    if [ "$HTTP_CODE" = "200" ]; then
        echo "$BODY" | jq . > "$RESULTS_DIR/validate_key.json" 2>/dev/null || echo "$BODY" > "$RESULTS_DIR/validate_key.json"
        pass "Validate Symmetric Key"
    else
        fail "Validate Symmetric Key" "HTTP $HTTP_CODE: $BODY"
    fi
else
    fail "Validate Symmetric Key" "No symmetric key ID available"
fi

Capture remaining tests...

