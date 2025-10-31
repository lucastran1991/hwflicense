#!/bin/bash

# Master Start Script
# This script starts both KMS backend and Next.js frontend services

set -e

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
echo "Starting KMS Backend and Frontend Services"
echo "=========================================="
echo ""

# Check if directories exist
if [ ! -d "$KMS_DIR" ]; then
    print_error "KMS directory not found: $KMS_DIR"
    exit 1
fi

if [ ! -d "$INTERFACE_DIR" ]; then
    print_error "Interface directory not found: $INTERFACE_DIR"
    exit 1
fi

# Start KMS Backend first
print_service "Starting KMS Backend..."
if [ ! -f "$KMS_DIR/start.sh" ]; then
    print_error "KMS start.sh not found: $KMS_DIR/start.sh"
    exit 1
fi

cd "$KMS_DIR"
if ! bash "$KMS_DIR/start.sh"; then
    print_error "Failed to start KMS backend service"
    exit 1
fi

# Load ports from environment.json if available
ENV_CONFIG_FILE="$PROJECT_ROOT/config/environment.json"
DEFAULT_BACKEND_PORT=8080
DEFAULT_FRONTEND_PORT=3000

# Load backend port from environment.json
if [ -f "$ENV_CONFIG_FILE" ]; then
    # Use jq if available (most reliable)
    if command -v jq > /dev/null 2>&1; then
        BACKEND_PORT=$(jq -r '.backend.port // 8080' "$ENV_CONFIG_FILE" 2>/dev/null || echo "$DEFAULT_BACKEND_PORT")
        FRONTEND_PORT=$(jq -r '.frontend.port // 3000' "$ENV_CONFIG_FILE" 2>/dev/null || echo "$DEFAULT_FRONTEND_PORT")
    elif command -v node > /dev/null 2>&1; then
        # Use node to parse JSON
        BACKEND_PORT=$(node -e "try { const c=require('$ENV_CONFIG_FILE'); console.log(c.backend?.port || $DEFAULT_BACKEND_PORT); } catch(e) { console.log($DEFAULT_BACKEND_PORT); }" 2>/dev/null || echo "$DEFAULT_BACKEND_PORT")
        FRONTEND_PORT=$(node -e "try { const c=require('$ENV_CONFIG_FILE'); console.log(c.frontend?.port || $DEFAULT_FRONTEND_PORT); } catch(e) { console.log($DEFAULT_FRONTEND_PORT); }" 2>/dev/null || echo "$DEFAULT_FRONTEND_PORT")
    else
        # Fallback: use grep/sed
        BACKEND_PORT=$(grep -A 5 '"backend"' "$ENV_CONFIG_FILE" | grep -o '"port"[[:space:]]*:[[:space:]]*[0-9]*' | grep -o '[0-9]*' || echo "$DEFAULT_BACKEND_PORT")
        FRONTEND_PORT=$(grep -A 5 '"frontend"' "$ENV_CONFIG_FILE" | grep -o '"port"[[:space:]]*:[[:space:]]*[0-9]*' | grep -o '[0-9]*' || echo "$DEFAULT_FRONTEND_PORT")
        if [ -z "$BACKEND_PORT" ] || [ "$BACKEND_PORT" = "" ]; then
            BACKEND_PORT=$DEFAULT_BACKEND_PORT
        fi
        if [ -z "$FRONTEND_PORT" ] || [ "$FRONTEND_PORT" = "" ]; then
            FRONTEND_PORT=$DEFAULT_FRONTEND_PORT
        fi
    fi
else
    BACKEND_PORT=$DEFAULT_BACKEND_PORT
    FRONTEND_PORT=$DEFAULT_FRONTEND_PORT
fi

# Wait for backend health check
print_info "Waiting for KMS backend to be healthy..."
# Environment variable KMS_PORT takes precedence
KMS_PORT="${KMS_PORT:-:$BACKEND_PORT}"
HEALTH_PORT="$KMS_PORT"
if [[ ! "$HEALTH_PORT" =~ ^: ]]; then
    HEALTH_PORT=":$HEALTH_PORT"
fi

MAX_WAIT=30
WAITED=0
BACKEND_HEALTHY=0

while [ $WAITED -lt $MAX_WAIT ]; do
    if curl -s "http://localhost${HEALTH_PORT}/health" > /dev/null 2>&1; then
        BACKEND_HEALTHY=1
        print_info "KMS backend is healthy ✓"
        break
    fi
    sleep 1
    WAITED=$((WAITED + 1))
done

if [ $BACKEND_HEALTHY -eq 0 ]; then
    print_warning "KMS backend health check timeout, but continuing..."
fi

# Wait a moment before starting frontend
sleep 2

# Start Next.js Frontend
print_service "Starting Next.js Frontend..."
if [ ! -f "$INTERFACE_DIR/start.sh" ]; then
    print_error "Interface start.sh not found: $INTERFACE_DIR/start.sh"
    exit 1
fi

cd "$INTERFACE_DIR"
if ! bash "$INTERFACE_DIR/start.sh"; then
    print_error "Failed to start Next.js frontend service"
    print_warning "Backend is still running"
    exit 1
fi

# Summary
echo ""
echo "=========================================="
print_info "All services started successfully!"
echo "=========================================="
echo ""
print_info "Services Status:"
echo "  • KMS Backend:     http://localhost${HEALTH_PORT#:}"
echo "  • Next.js Frontend: http://localhost:${FRONTEND_PORT}"
echo ""
print_info "To check status:  cd script && ./status.sh"
print_info "To stop services: cd script && ./stop.sh"
print_info "To view logs:"
echo "  • KMS:     tail -f $KMS_DIR/kms.log"
echo "  • Frontend: tail -f $INTERFACE_DIR/interface.log"
echo ""

