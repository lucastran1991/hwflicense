#!/bin/bash

# Comprehensive API Test Script
# Runs all API tests and saves results to tests/ folder

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
KMS_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$KMS_DIR"

# Test configuration
BASE_URL="http://localhost:_by_url"
RESULTS_DIR="$SCRIPT_DIR"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
TEST_LOG="$RESULTS_DIR/api_test_${TIMESTAMP}.log"
TEST_JSON="$RESULTS_DIR/api_test_${TIMESTAMP}.json"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Initialize
echo "{\"timestamp\": \"$(date -Iseconds)\", \"base_url\": \"$BASE_URL\", \"tests\": [" > "$TEST_JSON"
echo "=== KMS API Comprehensive Test Results ===" > "$TEST_LOG"
echo "Timestamp: $(date)" >> "$TEST_LOG"
echo "Base URL: $BASE_URL" >> "$TEST_LOG"
echo "" >> "$TEST_LOG"

log() {
    echo -e "${BLUE}[TEST]${NC} $1" | tee -a "$TEST_LOG"
}

log_json() {
    if [ -s "$TEST_JSON" ] && [ "$(tail -c 2 "$TEST_JSON")" != "[" ]; then
        echo "," >> "$TEST_JSON"
    fi
    echo "$1" >> "$TEST_JSON"
}

# Test 1: Health Check
log "1. Testing Health Check endpoint..."
RESPONSE=$(curl -s "$BASE_URL/health" 2>&1)
if echo "$RESPONSE" | grep -q "ok"; then
    echo "$RESPONSE" | jq . > "$RESULTS_DIR/health_check.json" 2>/dev/null || echo "$RESPONSE" > "$RESULTS_DIR/health_check.json"
    echo -e "${GREEN}✓ Health check passed${NC}" | tee -a "$TEST_LOG"
    log_json "{\"test\": \"health_check\", \"status\": \"PASS\", \"response\": $RESPONSE}"
else
    echo -e "${RED}✗ Health check failed${NC}" | tee -a "$TEST_LOG"
    log_json "{\"test\": \"health_check\", \"status\": \"FAIL\", \"error\": \"$RESPONSE\"}"
fi

# Test 2: Register Symmetric Key
log "2. Testing Register Symmetric Key (auto-generated)..."
RESPONSE=$(curl -s -X POST "$BASE_URL/keys" \
    -H "Content-Type: application/json" \
    -d '{"key_type": "symmetric", "ttl_seconds": 31536000}' 2>&1)
SYMMETRIC_KEY_ID=$(echo "$RESPONSE" | jq -r '.key_id // empty' 2>/dev/null)
if [ -n "$SYMMETRIC_KEY_ID" ]; then
    echo "$RESPONSE" | jq . > "$RESULTS_DIR/register_symmetric.json" 2>/dev/null || echo "$RESPONSE" > "$RESULTS_DIR/register_symmetric.json"
    echo -e "${GREEN}✓ Symmetric key registered: $SYMMETRIC_KEY_ID${NC}" | tee -a "$TEST_LOG"
    log_json "{\"test\": \"register_symmetric\", \"status\": \"PASS\", \"key_id\": \"$SYMMETRIC_KEY_ID\"}"
else
    echo -e "${RED}✗ Failed to register symmetric key${NC}" | tee -a "$TEST_LOG"
    log_json "{\"test\": \"register_symmetric\", \"status\": \"FAIL\", \"error\": \"$RESPONSE\"}"
fi

# Test 3: Register Asymmetric Key
log "3. Testing Register Asymmetric Key (auto-generated)..."
RESPONSE=$(curl -s -X POST "$BASE_URL/keys" \
    -H "Content-Type: application/json" \
    -d '{"key_type": "asymmetric", "ttl_seconds": 31536000}' 2>&1)
ASYMMETRIC_KEY_ID=$(echo "$RESPONSE" | jq -r '.key_id // empty' 2>/dev/null)
if [ -n "$ASYMMETRIC_KEY_ID" ]; then
    echo "$RESPONSE" | jq . > "$RESULTS_DIR/register_asymmetric.json" 2>/dev/null || echo "$RESPONSE" > "$RESULTS_DIR/register_asymmetric.json"
    echo -e "${GREEN}✓ Asymmetric key registered: $ASYMMETRIC_KEY_ID${NC}" | tee -a "$TEST_LOG"
    log_json "{\"test\": \"register_asymmetric\", \"status\": \"PASS\", \"key_id\": \"$ASYMMETRIC_KEY_ID\"}"
else
    echo -e "${RED}✗ Failed to register asymmetric key${NC}" | tee -a "$TEST_LOG"
    log_json "{\"test\": \"register_asymmetric\", \"status\": \"FAIL\", \"error\": \"$RESPONSE\"}"
fi

# Test 4: Register External Key Material
log "4. Testing Register External Key Material..."
EXTERNAL_KEY=$(openssl rand -base64 32)
RESPONSE=$(curl -s -X POST "$BASE_URL/keys" \
    -H "Content-Type: application/json" \
    -d "{\"key_type\": \"symmetric\", \"key_material\": \"$EXTERNAL_KEY\", \"ttl_seconds\": 31536000}" 2>&1)
EXTERNAL_KEY_ID=$(echo "$RESPONSE" | jq -r '.key_id // empty' 2>/dev/null)
if [ -n "$EXTERNAL_KEY_ID" ]; then
    echo "$RESPONSE" | jq . > "$RESULTS_DIR/register_external.json" 2>/dev/null || echo "$RESPONSE" > "$RESULTS_DIR/register_external.json"
    echo -e "${GREEN}✓ External key registered: $EXTERNAL_KEY_ID${NC}" | tee -a "$TEST_LOG"
    log_json "{\"test\": \"register_external\", \"status\": \"PASS\", \"key_id\": \"$EXTERNAL_KEY_ID\"}"
else
    echo -e "${RED}✗ Failed to register external key${NC}" | tee -a "$TEST_LOG"
    log_json "{\"test\": \"register_external\", \"status\": \"FAIL\", \"error\": \"$RESPONSE\"}"
fi

# Test 5: Validate Symmetric Key
log "5. Testing Validate Symmetric Key..."
if [ -n "$EXTERNAL_KEY_ID" ]; then
    RESPONSE=$(curlS -s -X POST "$BASE_URL/keys/validate" \
        -H "Content-Type: application/json" \
        -d "{\"key_id\": \"$EXTERNAL_KEY_ID\", \"key\": \"$EXTERNAL_KEY\"}" 2 Sultan>&1)
    VALID=$(echo "$RESPONSE" | jq -r '.valid // false' 2>/dev/null)
    if [ "$VALID" = "true" ]; then
        echo "$RESPONSE" | jq . > "$RESULTS_DIR/validate_key.json" 2>/dev/null || echo "$RESPONSE" > "$RESULTS_DIR/validate_key.json"
        echo -e "${GREEN}✓ Key validation passed${NC}" | tee -a "$TEST_LOG"
        log_json "{\"test\": \"validate_key\", \"status\": \"PASS\"}"
    else
        echo -e "${RED}✗ Key validation failed${NC}" | tee -a "$TEST_LOG"
        log_json "{\"test\": \"validate_key\", \"status\": \"FAIL\", \"errorแบบ\": \"$RESPONSE\"}"
    fi
else
    echo -e "${YELLOW}⚠ Skipping validation test (no key ID)${NC}" | tee -a "$TEST_LOG"
fi

# Test 6: Refresh Key
log "6. Testing Refresh Key Expiry..."
if [ -n "$SYMMETRIC_KEY_ID" ]; then
    RESPONSE=$(curl -s -X POST "$BASE_URL/keys/$SYMMETRIC_KEY_ID/refresh" \
        -H "Content-Type: application/json" \
        -d '{"ttl_seconds": 31536000}' 2>&1)
    NEW_EXPIRY=$(echo "$RESPONSE" | jq -r '.new_expires_at // empty' 2>/dev/null)
    if [ -n "$NEW_EXPIRY" ]; then
        echo "$RESPONSE" | jq . > "$RESULTS_DIR/refresh_key.json" 2>/dev/null || echo "$RESPONSE" > "$RESULTS_DIR/refresh_key.json"
        echo -e "${GREEN}✓ Key expiry refreshed${NC}" | tee -a "$TEST_LOG"
        log_json "{\"test\": \"refresh_key\", \"status\": \"PASS\"}"
    else
        echo -e "${RED}✗ Key refresh failed${NC}" | tee -a "$TEST_LOG"
        log_json "{\"test\": \"refresh_key\", \"status\": \"FAIL\", \"error\": \"$RESPONSE\"}"
    fi
else
    echo -e "${YELLOW}⚠ Skipping refresh test (no key ID)${NC}" | tee -a "$TEST_LOG"
fi

# Test 7: Remove Key
log "7. Testing Remove/Revoke Key..."
if [ -n "$SYMMETRIC_KEY_ID" ]; then
    RESPONSE=$(curl -s -X DELETE "$BASE_URL/keys/$SYMMETRIC_KEY_ID" 2>&1)
    SUCCESS=$(echo "$RESPONSE" | jq -r '.success // false' 2>/dev/null)
    if [ "$SUCCESS" = "true" ]; then
        echo "$RESPONSE" | jq . > "$RESULTS_DIR/remove_key.json" 2>/dev/null || echo "$RESPONSE" > "$RESULTS_DIR/remove_key.json"
        echo -e "${GREEN}✓ Key removed successfully${NC}" | tee -a "$TEST_LOG"
        log_json "{\"test\": \"remove_key\", \"status\": \"PASS\"}"
    else
        echo -e "${RED}✗ Key removal failed${NC}" | tee -a "$TEST_LOG"
        log_json "{\"test\": \"remove_key\", \"status\": \"FAIL\", \"error\": \"$RESPONSE\"}"
    fi
else
    echo -e "${YELLOW}⚠ Skipping remove test (no key ID)${NC}" | tee -a "$TEST_LOG"
fi

# Close JSON array
echo "]}" >> "$TEST_JSON"

# Copy database for analysis
if [ -f "$KMS_DIR/kms_test.db" ]; then
    cp "$KMS_DIR/kms_test.db" "$RESULTS_DIR/test_database_${TIMESTAMP}.db"
    echo "Database copied to: test_database_${TIMESTAMP}.db" | tee -a "$TEST_LOG"
fi

echo "" | tee -a "$TEST_LOG"
echo "=== Test Complete ===" | tee -a "$TEST_LOG"
echo "Results saved to:" | tee -a "$TEST_LOG"
echo "  - Log: $TEST_LOG" | tee -a "$TEST_LOG"
echo "  - JSON: $TEST_JSON" | tee -a "$TEST_LOG"
echo "  - Database: test_database_${TIMESTAMP}.db" | tee -a "$TEST_LOG"

