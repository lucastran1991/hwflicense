#!/bin/bash

# KMS Service Status Script
# This script checks the status of the KMS service

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Configuration
PID_FILE="$SCRIPT_DIR/kms.pid"
LOG_FILE="$SCRIPT_DIR/kms.log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

echo "=== KMS Service Status ==="
echo ""

# Check PID file
if [ ! -f "$PID_FILE" ]; then
    print_error "Service is not running (PID file not found)"
    exit 1
fi

PID=$(cat "$PID_FILE")

# Check if process is running
if ! ps -p "$PID" > /dev/null 2>&1; then
    print_error "Service is not running (process with PID $PID not found)"
    rm -f "$PID_FILE"
    exit 1
fi

print_info "Service is running"
echo "  PID: $PID"

# Get process details
if command -v ps > /dev/null 2>&1; then
    echo "  Process Info:"
    ps -p "$PID" -o pid,ppid,user,%cpu,%mem,etime,cmd | tail -n +2 | sed 's/^/    /'
fi

# Check health endpoint
export KMS_PORT="${KMS_PORT:-:8080}"
HEALTH_URL="http://localhost${KMS_PORT#:}/health"

echo ""
print_info "Health Check:"
if curl -s -f "$HEALTH_URL" > /dev/null 2>&1; then
    HEALTH_RESPONSE=$(curl -s "$HEALTH_URL")
    print_info "  ✓ Service is healthy"
    echo "  Response: $HEALTH_RESPONSE"
else
    print_warning "  ✗ Health check failed"
    print_warning "  URL: $HEALTH_URL"
fi

# Show configuration
echo ""
print_info "Configuration:"
echo "  KMS_DB_PATH: ${KMS_DB_PATH:-./kms.db}"
echo "  KMS_PORT: ${KMS_PORT:-:8080}"

# Show log file info
if [ -f "$LOG_FILE" ]; then
    echo ""
    print_info "Log File: $LOG_FILE"
    LOG_SIZE=$(wc -l < "$LOG_FILE" 2>/dev/null || echo "0")
    echo "  Lines: $LOG_SIZE"
    
    echo ""
    echo "Last 5 lines of log:"
    tail -n 5 "$LOG_FILE" 2>/dev/null | sed 's/^/  /' || echo "  (log file empty or not readable)"
fi

echo ""
echo "=== End Status ==="

