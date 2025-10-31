#!/bin/bash

# Master Status Script
# This script checks the status of both KMS backend and Next.js frontend services

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Paths
KMS_DIR="$PROJECT_ROOT/kms"
INTERFACE_DIR="$PROJECT_ROOT/interface"

# Load ports from environment.json if available
ENV_CONFIG_FILE="$PROJECT_ROOT/config/environment.json"
DEFAULT_BACKEND_PORT=8080
DEFAULT_FRONTEND_PORT=3000
DEFAULT_BACKEND_API_URL="http://localhost:8080"

# Load backend port from environment.json
if [ -f "$ENV_CONFIG_FILE" ]; then
    # Use jq if available (most reliable)
    if command -v jq > /dev/null 2>&1; then
        BACKEND_PORT=$(jq -r '.backend.port // 8080' "$ENV_CONFIG_FILE" 2>/dev/null || echo "$DEFAULT_BACKEND_PORT")
        FRONTEND_PORT=$(jq -r '.frontend.port // 3000' "$ENV_CONFIG_FILE" 2>/dev/null || echo "$DEFAULT_FRONTEND_PORT")
        BACKEND_API_URL=$(jq -r '.frontend.api_url // empty' "$ENV_CONFIG_FILE" 2>/dev/null || echo "")
        BACKEND_HOST=$(jq -r '.backend.host // "localhost"' "$ENV_CONFIG_FILE" 2>/dev/null || echo "localhost")
    elif command -v node > /dev/null 2>&1; then
        # Use node to parse JSON
        BACKEND_PORT=$(node -e "try { const c=require('$ENV_CONFIG_FILE'); console.log(c.backend?.port || $DEFAULT_BACKEND_PORT); } catch(e) { console.log($DEFAULT_BACKEND_PORT); }" 2>/dev/null || echo "$DEFAULT_BACKEND_PORT")
        FRONTEND_PORT=$(node -e "try { const c=require('$ENV_CONFIG_FILE'); console.log(c.frontend?.port || $DEFAULT_FRONTEND_PORT); } catch(e) { console.log($DEFAULT_FRONTEND_PORT); }" 2>/dev/null || echo "$DEFAULT_FRONTEND_PORT")
        BACKEND_API_URL=$(node -e "try { const c=require('$ENV_CONFIG_FILE'); console.log(c.frontend?.api_url || ''); } catch(e) { console.log(''); }" 2>/dev/null || echo "")
        BACKEND_HOST=$(node -e "try { const c=require('$ENV_CONFIG_FILE'); console.log(c.backend?.host || 'localhost'); } catch(e) { console.log('localhost'); }" 2>/dev/null || echo "localhost")
    else
        # Fallback: use grep/sed
        BACKEND_PORT=$(grep -A 5 '"backend"' "$ENV_CONFIG_FILE" | grep -o '"port"[[:space:]]*:[[:space:]]*[0-9]*' | grep -o '[0-9]*' || echo "$DEFAULT_BACKEND_PORT")
        FRONTEND_PORT=$(grep -A 5 '"frontend"' "$ENV_CONFIG_FILE" | grep -o '"port"[[:space:]]*:[[:space:]]*[0-9]*' | grep -o '[0-9]*' || echo "$DEFAULT_FRONTEND_PORT")
        BACKEND_API_URL=$(grep -A 5 '"frontend"' "$ENV_CONFIG_FILE" | grep -o '"api_url"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"' | tr -d '"' || echo "")
        BACKEND_HOST=$(grep -A 5 '"backend"' "$ENV_CONFIG_FILE" | grep -o '"host"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"' | tr -d '"' || echo "localhost")
        if [ -z "$BACKEND_PORT" ] || [ "$BACKEND_PORT" = "" ]; then
            BACKEND_PORT=$DEFAULT_BACKEND_PORT
        fi
        if [ -z "$FRONTEND_PORT" ] || [ "$FRONTEND_PORT" = "" ]; then
            FRONTEND_PORT=$DEFAULT_FRONTEND_PORT
        fi
    fi
    
    # Construct API URL if not found
    if [ -z "$BACKEND_API_URL" ] || [ "$BACKEND_API_URL" = "" ]; then
        BACKEND_API_URL="http://${BACKEND_HOST}:${BACKEND_PORT}"
    fi
else
    BACKEND_PORT=$DEFAULT_BACKEND_PORT
    FRONTEND_PORT=$DEFAULT_FRONTEND_PORT
    BACKEND_API_URL=$DEFAULT_BACKEND_API_URL
    BACKEND_HOST="localhost"
fi

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
            # Environment variable KMS_PORT takes precedence
            KMS_PORT="${KMS_PORT:-:$BACKEND_PORT}"
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
            echo "  KMS_PORT: ${KMS_PORT:-:$BACKEND_PORT} (from environment.json: $BACKEND_PORT)"
            
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
    # Environment variable PORT takes precedence
    INTERFACE_PORT="${PORT:-$FRONTEND_PORT}"
    
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
            echo "  Port: $INTERFACE_PORT (from environment.json: $FRONTEND_PORT)"
            echo "  API URL: ${NEXT_PUBLIC_API_URL:-$BACKEND_API_URL}"
            
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

