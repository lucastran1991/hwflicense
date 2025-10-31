#!/bin/bash

# Complete KMS API Test Script
# This script tests all available KMS API endpoints

BASE_URL="http://localhost:8080"
COLOR_GREEN='\033[0;32m'
COLOR_RED='\033[0;31m'
COLOR_YELLOW='\033[1;33m'
COLOR_BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "\n${COLOR_BLUE}========================================${NC}"
    echo -e "${COLOR_BLUE}$1${NC}"
    echo -e "${COLOR_BLUE}========================================${NC}\n"
}

print_success() {
    echo -e "${COLOR_GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${COLOR_RED}✗ $1${NC}"
}

print_info() {
    echo -e "${COLOR_YELLOW}ℹ $1${NC}"
}

# Check if service is running
print_header "Checking Service Status"

if ! curl -s -f "$BASE_URL/health" > /dev/null 2>&1; then
    print_error "Service is not running or not accessible at $BASE_URL"
    print_info "Please start the service first: ./start.sh"
    exit 1
fi

print_success "Service is running"

# Test Health Check
print_header "1. Health Check"

HEALTH_RESPONSE=$(curl -s -X GET "$BASE_URL/health")
echo "Response: $HEALTH_RESPONSE"
if echo "$HEALTH_RESPONSE" | grep -q "ok"; then
    print_success "Health check passed"
else
    print_error "Health check failed"
fi

# Test Register Symmetric Key
print_header "2. Register Symmetric Key (Auto-generated)"

SYM_RESPONSE=$(curl -s -X POST "$BASE_URL/keys" \
  -H "Content-Type: application/json" \
  -d '{
    "key_type": "symmetric",
    "expires_in_seconds": 31536000
  }')

echo "$SYM_RESPONSE" | jq . 2>/dev/null || echo "$SYM_RESPONSE"

SYM_KEY_ID=$(echo "$SYM_RESPONSE" | jq -r '.key_id' 2>/dev/null)

if [ -n "$SYM_KEY_ID" ] && [ "$SYM_KEY_ID" != "null" ]; then
    print_success "Symmetric key registered: $SYM_KEY_ID"
else
    print_error "Failed to register symmetric key"
fi

# Test Register Asymmetric Key
print_header "3. Register Asymmetric Key (Auto-generated)"

ASYM_RESPONSE=$(curl -s -X POST "$BASE_URL/keys" \
  -H "Content-Type: application/json" \
  -d '{
    "key_type": "asymmetric",
    "expires_in_seconds": 31536000
  }')

echo "$ASYM_RESPONSE" | jq . 2>/dev/null || echo "$ASYM_RESPONSE"

ASYM_KEY_ID=$(echo "$ASYM_RESPONSE" | jq -r '.key_id' 2>/dev/null)
ASYM_PUBLIC_KEY=$(echo "$ASYM_RESPONSE" | jq -r '.public_key' 2>/dev/null)

if [ -n "$ASYM_KEY_ID" ] && [ "$ASYM_KEY_ID" != "null" ]; then
    print_success "Asymmetric key registered: $ASYM_KEY_ID"
    print_info "Public Key (first 50 chars): ${ASYM_PUBLIC_KEY:0:50}..."
else
    print_error "Failed to register asymmetric key"
fi

# Test Register with External Key Material (Symmetric)
print_header "4. Register Symmetric Key (External Key Material)"

# Generate a test key
TEST_KEY=$(openssl rand -base64 32)

EXT_SYM_RESPONSE=$(curl -s -X POST "$BASE_URL/keys" \
  -H "Content-Type: application/json" \
  -d "{
    \"key_type\": \"symmetric\",
    \"expires_in_seconds\": 31536000,
    \"key_material\": \"$TEST_KEY\"
  }")

echo "$EXT_SYM_RESPONSE" | jq . 2>/dev/null || echo "$EXT_SYM_RESPONSE"

EXT_SYM_KEY_ID=$(echo "$EXT_SYM_RESPONSE" | jq -r '.key_id' 2>/dev/null)

if [ -n "$EXT_SYM_KEY_ID" ] && [ "$EXT_SYM_KEY_ID" != "null" ]; then
    print_success "External symmetric key registered: $EXT_SYM_KEY_ID"
    STORED_TEST_KEY="$TEST_KEY"
else
    print_error "Failed to register external symmetric key"
fi

# Test Validate Symmetric Key (if we have a key)
if [ -n "$EXT_SYM_KEY_ID" ] && [ -n "$STORED_TEST_KEY" ]; then
    print_header "5. Validate Symmetric Key"
    
    VALIDATE_SYM_RESPONSE=$(curl -s -X POST "$BASE_URL/keys/validate" \
      -H "Content-Type: application/json" \
      -d "{
        \"key_id\": \"$EXT_SYM_KEY_ID\",
        \"key_material\": \"$STORED_TEST_KEY\"
      }")
    
    echo "$VALIDATE_SYM_RESPONSE" | jq . 2>/dev/null || echo "$VALIDATE_SYM_RESPONSE"
    
    VALID=$(echo "$VALIDATE_SYM_RESPONSE" | jq -r '.valid' 2>/dev/null)
    
    if [ "$VALID" = "true" ]; then
        print_success "Symmetric key validation passed"
    else
        print_error "Symmetric key validation failed"
    fi
fi

# Test Validate with Wrong Key
if [ -n "$EXT_SYM_KEY_ID" ]; then
    print_header "6. Validate Symmetric Key (Wrong Key - Should Fail)"
    
    WRONG_KEY=$(openssl rand -base64 32)
    
    VALIDATE_WRONG_RESPONSE=$(curl -s -X POST "$BASE_URL/keys/validate" \
      -H "Content-Type: application/json" \
      -d "{
        \"key_id\": \"$EXT_SYM_KEY_ID\",
        \"key_material\": \"$WRONG_KEY\"
      }")
    
    echo "$VALIDATE_WRONG_RESPONSE" | jq . 2>/dev/null || echo "$VALIDATE_WRONG_RESPONSE"
    
    VALID=$(echo "$VALIDATE_WRONG_RESPONSE" | jq -r '.valid' 2>/dev/null)
    
    if [ "$VALID" = "false" ]; then
        print_success "Correctly rejected wrong key"
    else
        print_error "Should have rejected wrong key"
    fi
fi

# Test Refresh Key
if [ -n "$SYM_KEY_ID" ] && [ "$SYM_KEY_ID" != "null" ]; then
    print_header "7. Refresh Key Expiry"
    
    REFRESH_RESPONSE=$(curl -s -X POST "$BASE_URL/keys/$SYM_KEY_ID/refresh" \
      -H "Content-Type: application/json" \
      -d '{
        "expires_in_seconds": 63072000
      }')
    
    echo "$REFRESH_RESPONSE" | jq . 2>/dev/null || echo "$REFRESH_RESPONSE"
    
    NEW_EXPIRY=$(echo "$REFRESH_RESPONSE" | jq -r '.new_expires_at' 2>/dev/null)
    
    if [ -n "$NEW_EXPIRY" ] && [ "$NEW_EXPIRY" != "null" ]; then
        print_success "Key expiry refreshed: $NEW_EXPIRY"
    else
        print_error "Failed to refresh key expiry"
    fi
fi

# Test Remove/Revoke Key
if [ -n "$SYM_KEY_ID" ] && [ "$SYM_KEY_ID" != "null" ]; then
    print_header "8. Remove/Revoke Key"
    
    REMOVE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/keys/$SYM_KEY_ID")
    
    echo "$REMOVE_RESPONSE" | jq . 2>/dev/null || echo "$REMOVE_RESPONSE"
    
    SUCCESS=$(echo "$REMOVE_RESPONSE" | jq -r '.success' 2>/dev/null)
    
    if [ "$SUCCESS" = "true" ]; then
        print_success "Key revoked successfully"
    else
        print_error "Failed to revoke key"
    fi
    
    # Test that revoked key is no longer valid
    print_header "9. Validate Revoked Key (Should Show Revoked)"
    
    VALIDATE_REVOKED_RESPONSE=$(curl -s -X POST "$BASE_URL/keys/validate" \
      -H "Content-Type: application/json" \
      -d "{
        \"key_id\": \"$SYM_KEY_ID\",
        \"key_material\": \"dGVzdA==\"
      }")
    
    echo "$VALIDATE_REVOKED_RESPONSE" | jq . 2>/dev/null || echo "$VALIDATE_REVOKED_RESPONSE"
    
    REVOKED=$(echo "$VALIDATE_REVOKED_RESPONSE" | jq -r '.revoked' 2>/dev/null)
    
    if [ "$REVOKED" = "true" ]; then
        print_success "Correctly detected revoked key"
    else
        print_info "Key validation response: $VALIDATE_REVOKED_RESPONSE"
    fi
fi

# Test Error Cases
print_header "10. Error Cases"

# Invalid key type
print_info "Testing invalid key type..."
INVALID_TYPE_RESPONSE=$(curl -s -X POST "$BASE_URL/keys" \
  -H "Content-Type: application/json" \
  -d '{
    "key_type": "invalid",
    "expires_in_seconds": 31536000
  }')

echo "$INVALID_TYPE_RESPONSE" | jq . 2>/dev/null || echo "$INVALID_TYPE_RESPONSE"

if echo "$INVALID_TYPE_RESPONSE" | grep -q "error"; then
    print_success "Correctly rejected invalid key type"
else
    print_error "Should have rejected invalid key type"
fi

# Missing required fields
print_info "Testing missing required fields..."
MISSING_FIELD_RESPONSE=$(curl -s -X POST "$BASE_URL/keys" \
  -H "Content-Type: application/json" \
  -d '{
    "expires_in_seconds": 31536000
  }')

echo "$MISSING_FIELD_RESPONSE" | jq . 2>/dev/null || echo "$MISSING_FIELD_RESPONSE"

if echo "$MISSING_FIELD_RESPONSE" | grep -q "error"; then
    print_success "Correctly rejected request with missing fields"
else
    print_error "Should have rejected request with missing fields"
fi

# Non-existent key
print_info "Testing non-existent key..."
NOT_FOUND_RESPONSE=$(curl -s -X POST "$BASE_URL/keys/validate" \
  -H "Content-Type: application/json" \
  -d '{
    "key_id": "non-existent-key-id",
    "key_material": "dGVzdA=="
  }')

echo "$NOT_FOUND_RESPONSE" | jq . 2>/dev/null || echo "$NOT_FOUND_RESPONSE"

if echo "$NOT_FOUND_RESPONSE" | grep -q "not found\|404"; then
    print_success "Correctly handled non-existent key"
else
    print_info "Response: $NOT_FOUND_RESPONSE"
fi

# Summary
print_header "Test Summary"

print_info "All API endpoints have been tested"
print_info "Check the results above for any errors"
print_info ""
print_info "Service URL: $BASE_URL"
print_info "To view logs: tail -f kms.log"
print_info "To stop service: ./stop.sh"

echo ""

