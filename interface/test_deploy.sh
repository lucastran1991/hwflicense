#!/bin/bash

# Deployment Test Script for Next.js Interface
# This script validates the deployment setup without actually starting the service

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

print_header() {
    echo ""
    echo "=========================================="
    echo "$1"
    echo "=========================================="
}

TEST_PASSED=0
TEST_FAILED=0

# Test function
run_test() {
    local test_name="$1"
    local test_command="$2"
    
    if eval "$test_command" > /dev/null 2>&1; then
        print_info "$test_name"
        ((TEST_PASSED++))
        return 0
    else
        print_error "$test_name"
        ((TEST_FAILED++))
        return 1
    fi
}

print_header "Next.js Interface Deployment Test"

# Test 1: Check if required files exist
echo ""
echo "Checking required files..."
run_test "package.json exists" "test -f package.json"
run_test "next.config.js exists" "test -f next.config.js"
run_test "tsconfig.json exists" "test -f tsconfig.json"
run_test "start.sh exists" "test -f start.sh"
run_test "stop.sh exists" "test -f stop.sh"

# Test 2: Check file permissions
echo ""
echo "Checking file permissions..."
run_test "start.sh is executable" "test -x start.sh"
run_test "stop.sh is executable" "test -x stop.sh"

# Test 3: Check script syntax
echo ""
echo "Validating script syntax..."
run_test "start.sh syntax is valid" "bash -n start.sh"
run_test "stop.sh syntax is valid" "bash -n stop.sh"

# Test 4: Check Node.js and npm
echo ""
echo "Checking Node.js environment..."
if command -v node > /dev/null 2>&1; then
    NODE_VERSION=$(node --version | cut -d'v' -f2 | cut -d'.' -f1)
    if [ "$NODE_VERSION" -ge 16 ]; then
        print_info "Node.js version compatible: $(node --version)"
        ((TEST_PASSED++))
    else
        print_error "Node.js version too old: $(node --version) (requires >= 16)"
        ((TEST_FAILED++))
    fi
else
    print_error "Node.js not found"
    ((TEST_FAILED++))
fi

if command -v npm > /dev/null 2>&1; then
    NPM_VERSION=$(npm --version | cut -d'.' -f1)
    if [ "$NPM_VERSION" -ge 8 ]; then
        print_info "npm version compatible: $(npm --version)"
        ((TEST_PASSED++))
    else
        print_error "npm version too old: $(npm --version) (requires >= 8)"
        ((TEST_FAILED++))
    fi
else
    print_error "npm not found"
    ((TEST_FAILED++))
fi

# Test 5: Check package.json scripts
echo ""
echo "Checking package.json scripts..."
if grep -q '"start"' package.json; then
    print_info "npm start script found"
    ((TEST_PASSED++))
else
    print_error "npm start script not found"
    ((TEST_FAILED++))
fi

if grep -q '"build"' package.json; then
    print_info "npm build script found"
    ((TEST_PASSED++))
else
    print_error "npm build script not found"
    ((TEST_FAILED++))
fi

# Test 6: Check dependencies (if node_modules exists)
echo ""
echo "Checking dependencies..."
if [ -d "node_modules" ]; then
    print_info "node_modules directory exists"
    ((TEST_PASSED++))
    
    # Check for critical dependencies
    if [ -d "node_modules/next" ]; then
        print_info "Next.js dependency installed"
        ((TEST_PASSED++))
    else
        print_warning "Next.js dependency not found in node_modules"
    fi
    
    if [ -d "node_modules/@chakra-ui/react" ]; then
        print_info "ChakraUI dependency installed"
        ((TEST_PASSED++))
    else
        print_warning "ChakraUI dependency not found in node_modules"
    fi
else
    print_warning "node_modules not found (will be installed on first run)"
fi

# Test 7: Check environment configuration
echo ""
echo "Checking environment configuration..."
if [ -f ".env.local" ]; then
    print_info ".env.local exists"
    if grep -q "NEXT_PUBLIC_API_URL" .env.local; then
        print_info "NEXT_PUBLIC_API_URL is configured"
        ((TEST_PASSED++))
    else
        print_warning "NEXT_PUBLIC_API_URL not found in .env.local"
    fi
else
    print_warning ".env.local not found (will use defaults)"
    if [ -f ".env.example" ]; then
        print_info ".env.example exists (can be copied to .env.local)"
    fi
fi

# Test 8: Check for build artifacts (optional)
echo ""
echo "Checking build artifacts..."
if [ -d ".next" ]; then
    print_info ".next build directory exists"
    ((TEST_PASSED++))
    
    if [ -f ".next/BUILD_ID" ]; then
        print_info "Build artifacts found"
        ((TEST_PASSED++))
    else
        print_warning "Build artifacts incomplete (run npm run build)"
    fi
else
    print_warning ".next build directory not found (will be created on build)"
fi

# Test 9: Check port availability (optional)
echo ""
echo "Checking port availability..."
PORT="${PORT:-3000}"
if command -v lsof > /dev/null 2>&1; then
    if lsof -ti :$PORT > /dev/null 2>&1; then
        print_warning "Port $PORT is already in use"
        print_warning "  This may prevent the service from starting"
        print_warning "  Run 'lsof -ti :$PORT | xargs kill' to free the port"
    else
        print_info "Port $PORT is available"
        ((TEST_PASSED++))
    fi
else
    print_warning "lsof not available, cannot check port"
fi

# Test 10: Check PID file status
echo ""
echo "Checking service status..."
PID_FILE="$SCRIPT_DIR/interface.pid"
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE" 2>/dev/null || echo "")
    if [ -n "$PID" ] && ps -p "$PID" > /dev/null 2>&1; then
        print_warning "Service appears to be running (PID: $PID)"
        print_warning "  Run './stop.sh' before starting again"
    else
        print_info "Stale PID file found (will be cleaned up)"
    fi
else
    print_info "No existing service running"
    ((TEST_PASSED++))
fi

# Summary
echo ""
print_header "Test Summary"
echo "Passed: $TEST_PASSED"
echo "Failed: $TEST_FAILED"
echo ""

if [ $TEST_FAILED -eq 0 ]; then
    print_info "All tests passed! Deployment script is ready."
    echo ""
    echo "To start the service, run:"
    echo "  ./start.sh"
    echo ""
    echo "To stop the service, run:"
    echo "  ./stop.sh"
    exit 0
else
    print_error "Some tests failed. Please fix the issues above before deploying."
    exit 1
fi

