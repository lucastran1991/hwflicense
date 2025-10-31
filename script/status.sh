#!/bin/bash

# Master Status Script
# This script checks the status of both KMS backend and Next.js frontend services

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Paths
KMS_DIR="$PROJECT_ROOT/kms"
INTERFACE_DIR="$PROJECT_ROOT/interface"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored messages
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_service() {
    echo -e "${BLUE}[SERVICE]${NC} $1"
}

echo ""
echo "=========================================="
echo "  KMS System Status"
echo "=========================================="
echo ""

# Check KMS Backend Status
print_service "=== KMS Backend ==="

if [ ! -d "$KMS_DIR" ]; then
    print_error "KMS directory not found: $KMS_DIR"
else
    KMS_PID_FILE="$KMS_DIR/kms.pid"
    KMS_LOG_FILE="$KMS_DIR/kms.log"
    
    if [ ! -f "$KMS_PID_FILE" ]; then
        print_warning "Backend is not running (PID file not found)"
    else
        KMS_PID=$(cat "$KMS_PID_FILE" 2>/dev/null || echo "")
        if [ -z "$KMS_PID" ]; then
            print_warning "Backend PID file is empty"
        elif ! ps -p "$KMS_PID" > /dev/null 2>&1; then
            print_warning "Backend is not running (process $KMS_PID not found)"
        else
            print_info "Backend is running (PID: $KMS_PID)"
            
            # Process info
            if command -v ps > /dev/null 2>&1; then
                echo "  Process Info:"
                ps -p "$KMS_PID" -o pid,ppid,user,%cpu,%mem,etime,cmd | tail -n +2 | sed 's/^/    /'
            fi
            
            # Health check
            export KMS_PORT="${KMS_PORT:-:8080}"
            HEALTH_PORT="$KMS_PORT"
            if [[ ! "$HEALTH_PORT" =~ ^: ]]; then
                HEALTH_PORT=":$HEALTH_PORT"
            fi
            HEALTH_URL="http://localhost${HEALTH_PORT#:}/health"
            
            echo ""
            print_info "Health Check:"
            if curl -s -f "$HEALTH_URL" > /dev/null 2>&1; then
                HEALTH_RESPONSE=$(curl -s "$HEALTH_URL")
                print_info "  ✓ Service is healthy"
                echo "  Response: $HEALTH_RESPONSE"
                echo "  URL: $HEALTH_URL"
            else
                print_warning "  ✗ Health check failed"
                print_warning "  URL: $HEALTH_URL"
            fi
            
            # Configuration
            echo ""
            print_info "Configuration:"
            echo "  KMS_DB_PATH: ${KMS_DB_PATH:-./kms.db}"
            echo "  KMS_PORT: ${KMS_PORT:-:8080}"
            
            # Log file
            if [ -f "$KMS_LOG_FILE" ]; then
                echo ""
                print_info "Log File: $KMS_LOG_FILE"
                LOG_SIZE=$(wc -l < "$KMS_LOG_FILE" 2>/dev/null || echo "0")
                echo "  Lines: $LOG_SIZE"
                echo ""
                echo "  Last 3 lines of log:"
                tail -n 3 "$KMS_LOG_FILE" 2>/dev/null | sed 's/^/    /' || echo "    (log file empty or not readable)"
            fi
        fi
    fi
fi

echo ""
echo "----------------------------------------"
echo ""

# Check Next.js Frontend Status
print_service "=== Next.js Frontend ==="

if [ ! -d "$INTERFACE_DIR" ]; then
    print_error "Interface directory not found: $INTERFACE_DIR"
else
    INTERFACE_PID_FILE="$INTERFACE_DIR/interface.pid"
    INTERFACE_LOG_FILE="$INTERFACE_DIR/interface.log"
    INTERFACE_PORT="${PORT:-3000}"
    
    if [ ! -f "$INTERFACE_PID_FILE" ]; then
        print_warning "Frontend is not running (PID file not found)"
    else
        INTERFACE_PID=$(cat "$INTERFACE_PID_FILE" 2>/dev/null || echo "")
        if [ -z "$INTERFACE_PID" ]; then
            print_warning "Frontend PID file is empty"
        elif ! ps -p "$INTERFACE_PID" > /dev/null 2>&1; then
            print_warning "Frontend is not running (process $INTERFACE_PID not found)"
        else
            print_info "Frontend is running (PID: $INTERFACE_PID)"
            
            # Process info
            if command -v ps > /dev/null 2>&1; then
                echo "  Process Info:"
                ps -p "$INTERFACE_PID" -o pid,ppid,user,%cpu,%mem,etime,cmd | tail -n +2 | sed 's/^/    /'
            fi
            
            # Health check
            HEALTH_URL="http://localhost:$INTERFACE_PORT"
            
            echo ""
            print_info "Health Check:"
            if curl -s -f "$HEALTH_URL" > /dev/null 2>&1; then
                print_info "  ✓ Service is healthy"
                echo "  URL: $HEALTH_URL"
            else
                print_warning "  ✗ Health check failed"
                print_warning "  URL: $HEALTH_URL"
            fi
            
            # Configuration
            echo ""
            print_info "Configuration:"
            echo "  Port: $INTERFACE_PORT"
            echo "  API URL: ${NEXT_PUBLIC_API_URL:-http://localhost:8080}"
            
            # Log file
            if [ -f "$INTERFACE_LOG_FILE" ]; then
                echo ""
                print_info "Log File: $INTERFACE_LOG_FILE"
                LOG_SIZE=$(wc -l < "$INTERFACE_LOG_FILE" 2>/dev/null || echo "0")
                echo "  Lines: $LOG_SIZE"
                echo ""
                echo "  Last 3 lines of log:"
                tail -n 3 "$INTERFACE_LOG_FILE" 2>/dev/null | sed 's/^/    /' || echo "    (log file empty or not readable)"
            fi
        fi
    fi
fi

echo ""
echo "=========================================="
echo "=== End Status ==="
echo ""

