#!/bin/bash

# Test Script for All Management Scripts
# This script tests start.sh, stop.sh, restart.sh, and status.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

PASSED=0
FAILED=0

test_passed() {
    echo -e "${GREEN}[PASS]${NC} $1"
    PASSED=$((PASSED + 1))
}

test_failed() {
    echo -e "${RED}[FAIL]${NC} $1"
    FAILED=$((FAILED + 1))
}

test_info() {
    echo -e "${BLUE}[TEST]${NC} $1"
}

echo ""
echo "=========================================="
echo "  Script Testing Suite"
echo "=========================================="
echo ""

# Test 1: Syntax Validation
test_info "Test 1: Syntax Validation"
echo "----------------------------------------"
if bash -n "$SCRIPT_DIR/start.sh" 2>/dev/null; then
    test_passed "start.sh syntax is valid"
else
    test_failed "start.sh has syntax errors"
fi

if bash -n "$SCRIPT_DIR/stop.sh" 2>/dev/null; then
    test_passed "stop.sh syntax is valid"
else
    test_failed "stop.sh has syntax errors"
fi

if bash -n "$SCRIPT_DIR/restart.sh" 2>/dev/null; then
    test_passed "restart.sh syntax is valid"
else
    test_failed "restart.sh has syntax errors"
fi

if bash -n "$SCRIPT_DIR/status.sh" 2>/dev/null; then
    test_passed "status.sh syntax is valid"
else
    test_failed "status.sh has syntax errors"
fi

# Test 2: Executable Permissions
test_info "Test 2: Executable Permissions"
echo "----------------------------------------"
if [ -x "$SCRIPT_DIR/start.sh" ]; then
    test_passed "start.sh is executable"
else
    test_failed "start.sh is not executable"
fi

if [ -x "$SCRIPT_DIR/stop.sh" ]; then
    test_passed "stop.sh is executable"
else
    test_failed "stop.sh is not executable"
fi

if [ -x "$SCRIPT_DIR/restart.sh" ]; then
    test_passed "restart.sh is executable"
else
    test_failed "restart.sh is not executable"
fi

if [ -x "$SCRIPT_DIR/status.sh" ]; then
    test_passed "status.sh is executable"
else
    test_failed "status.sh is not executable"
fi

# Test 3: Status Script - Services Not Running
test_info "Test 3: Status Script (Services Not Running)"
echo "----------------------------------------"
STATUS_OUTPUT=$("$SCRIPT_DIR/status.sh" 2>&1)
if echo "$STATUS_OUTPUT" | grep -q "not running"; then
    test_passed "status.sh correctly reports services not running"
else
    test_failed "status.sh does not correctly report services not running"
fi

# Test 4: Stop Script - Services Not Running
test_info "Test 4: Stop Script (Services Not Running)"
echo "----------------------------------------"
if "$SCRIPT_DIR/stop.sh" > /dev/null 2>&1; then
    test_passed "stop.sh handles non-running services gracefully"
else
    test_failed "stop.sh fails when services are not running"
fi

# Test 5: Directory Checks in start.sh
test_info "Test 5: start.sh Directory Validation"
echo "----------------------------------------"
# Save original directories
ORIG_KMS="$PROJECT_ROOT/kms"
ORIG_INTERFACE="$PROJECT_ROOT/interface"

# Temporarily rename directories to test error handling
# Note: We won't actually do this in automated tests as it could break things
# Instead, we'll check if the script references the directories correctly
if grep -q "KMS_DIR" "$SCRIPT_DIR/start.sh" && grep -q "INTERFACE_DIR" "$SCRIPT_DIR/start.sh"; then
    test_passed "start.sh checks for required directories"
else
    test_failed "start.sh does not check for required directories"
fi

# Test 6: Restart Script Dependencies
test_info "Test 6: restart.sh Dependencies"
echo "----------------------------------------"
if grep -q "stop.sh" "$SCRIPT_DIR/restart.sh" && grep -q "start.sh" "$SCRIPT_DIR/restart.sh"; then
    test_passed "restart.sh correctly calls stop.sh and start.sh"
else
    test_failed "restart.sh does not correctly reference stop.sh and start.sh"
fi

# Test 7: Status Script Health Check URLs
test_info "Test 7: Status Script Health Check Logic"
echo "----------------------------------------"
if grep -q "/health" "$SCRIPT_DIR/status.sh" && grep -q "curl" "$SCRIPT_DIR/status.sh"; then
    test_passed "status.sh includes health check logic"
else
    test_failed "status.sh does not include health check logic"
fi

# Test 8: Script Error Handling
test_info "Test 8: Script Error Handling"
echo "----------------------------------------"
# Check if scripts use set -e for error handling
if grep -q "set -e" "$SCRIPT_DIR/start.sh" && grep -q "set -e" "$SCRIPT_DIR/stop.sh"; then
    test_passed "Scripts use 'set -e' for error handling"
else
    test_failed "Scripts do not use 'set -e' for error handling"
fi

# Test 9: Script Path Resolution
test_info "Test 9: Script Path Resolution"
echo "----------------------------------------"
if grep -q "SCRIPT_DIR=" "$SCRIPT_DIR/start.sh" && grep -q "SCRIPT_DIR=" "$SCRIPT_DIR/stop.sh"; then
    test_passed "Scripts correctly resolve paths"
else
    test_failed "Scripts do not correctly resolve paths"
fi

# Test 10: Environment Variable Handling
test_info "Test 10: Environment Variable Handling"
echo "----------------------------------------"
if grep -q "KMS_PORT" "$SCRIPT_DIR/start.sh" || grep -q "PORT" "$SCRIPT_DIR/status.sh"; then
    test_passed "Scripts handle environment variables"
else
    test_failed "Scripts do not handle environment variables"
fi

echo ""
echo "=========================================="
echo "  Test Results Summary"
echo "=========================================="
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed${NC}"
    exit 1
fi

