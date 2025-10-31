#!/bin/bash

# Next.js Interface Start Script
# This script builds and starts the Next.js interface

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Configuration
PID_FILE="$SCRIPT_DIR/interface.pid"
LOG_FILE="$SCRIPT_DIR/interface.log"

# Load port from environment.json if available
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
ENV_CONFIG_FILE="$PROJECT_ROOT/config/environment.json"
DEFAULT_PORT=3000

# Try to load port and api_url from environment.json
if [ -f "$ENV_CONFIG_FILE" ]; then
    # Use jq if available (most reliable)
    if command -v jq > /dev/null 2>&1; then
        FRONTEND_PORT=$(jq -r '.frontend.port // 3000' "$ENV_CONFIG_FILE" 2>/dev/null || echo "$DEFAULT_PORT")
        FRONTEND_API_URL=$(jq -r '.frontend.api_url // empty' "$ENV_CONFIG_FILE" 2>/dev/null || echo "")
    elif command -v node > /dev/null 2>&1; then
        # Use node to parse JSON (more reliable than grep/sed)
        FRONTEND_PORT=$(node -e "try { const c=require('$ENV_CONFIG_FILE'); console.log(c.frontend?.port || $DEFAULT_PORT); } catch(e) { console.log($DEFAULT_PORT); }" 2>/dev/null || echo "$DEFAULT_PORT")
        FRONTEND_API_URL=$(node -e "try { const c=require('$ENV_CONFIG_FILE'); console.log(c.frontend?.api_url || ''); } catch(e) { console.log(''); }" 2>/dev/null || echo "")
    else
        # Fallback: use grep/sed (less reliable but works without dependencies)
        FRONTEND_PORT=$(grep -A 5 '"frontend"' "$ENV_CONFIG_FILE" | grep -o '"port"[[:space:]]*:[[:space:]]*[0-9]*' | grep -o '[0-9]*' || echo "$DEFAULT_PORT")
        FRONTEND_API_URL=$(grep -A 5 '"frontend"' "$ENV_CONFIG_FILE" | grep -o '"api_url"[[:space:]]*:[[:space:]]*"[^"]*"' | grep -o '"[^"]*"' | tr -d '"' || echo "")
        if [ -z "$FRONTEND_PORT" ] || [ "$FRONTEND_PORT" = "" ]; then
            FRONTEND_PORT=$DEFAULT_PORT
        fi
    fi
    
    # Export API URL if found in environment.json
    if [ -n "$FRONTEND_API_URL" ] && [ "$FRONTEND_API_URL" != "" ]; then
        export NEXT_PUBLIC_API_URL="$FRONTEND_API_URL"
    fi
else
    FRONTEND_PORT=$DEFAULT_PORT
fi

# Environment variable PORT takes precedence
PORT="${PORT:-$FRONTEND_PORT}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
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

# Check if service is already running
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p "$PID" > /dev/null 2>&1; then
        print_warning "Next.js interface is already running (PID: $PID)"
        exit 1
    else
        print_info "Removing stale PID file"
        rm -f "$PID_FILE"
    fi
fi

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    print_info "Installing dependencies..."
    if ! npm install; then
        print_error "Failed to install dependencies"
        exit 1
    fi
    print_info "Dependencies installed successfully"
fi

# Check environment file
if [ ! -f ".env.local" ]; then
    print_warning ".env.local not found. Using default API URL: http://localhost:8080"
    if [ -f ".env.example" ]; then
        print_info "You can copy .env.example to .env.local and configure it"
    fi
fi

# Verify API URL is set before build (critical for Next.js as NEXT_PUBLIC_API_URL is embedded at build time)
if [ -z "$NEXT_PUBLIC_API_URL" ]; then
    print_warning "NEXT_PUBLIC_API_URL is not set. This will be embedded in the build."
    print_warning "API URL will default to: http://localhost:8080"
    if echo "$HOSTNAME" | grep -q "ctxdev\|production\|prod" || [ -n "$DEPLOY_ENV" ]; then
        print_error "WARNING: Detected production environment but API URL may be incorrect!"
        print_error "Please ensure frontend.api_url in environment.json is set to production URL"
    fi
else
    print_info "API URL for build: $NEXT_PUBLIC_API_URL"
    # Warn if API URL contains localhost in non-local environment
    if echo "$NEXT_PUBLIC_API_URL" | grep -q "localhost" && echo "$HOSTNAME" | grep -q "ctxdev\|production\|prod"; then
        print_warning "WARNING: API URL contains 'localhost' but appears to be in production environment!"
        print_warning "This may cause frontend to call wrong backend URL"
    fi
fi

# Build the application
print_info "Building Next.js application..."
print_info "Note: NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL will be embedded in the build"
if ! npm run build; then
    print_error "Failed to build Next.js application"
    exit 1
fi
print_info "Build successful"

print_info "Starting Next.js interface..."
print_info "Configuration:"
print_info "  Port: $PORT"
print_info "  API URL: ${NEXT_PUBLIC_API_URL:-http://localhost:8080}"
print_info "  Log File: $LOG_FILE"

# Export environment variables
export PORT=$PORT
if [ -f ".env.local" ]; then
    export $(grep -v '^#' .env.local | xargs)
fi

# Start the service in background
nohup npm start > "$LOG_FILE" 2>&1 &
NPM_PID=$!

# Wait a moment for Next.js to start
sleep 2

# Find the actual node process (Next.js server process)
PID=$(pgrep -P $NPM_PID 2>/dev/null | head -1)
if [ -z "$PID" ]; then
    # If no child process found, try to find by port
    if command -v lsof > /dev/null 2>&1; then
        PID=$(lsof -ti :$PORT 2>/dev/null | head -1)
    fi
    # If still no PID, use npm PID
    if [ -z "$PID" ]; then
        PID=$NPM_PID
    fi
fi

# Save PID
echo $PID > "$PID_FILE"

# Wait a moment to check if service started successfully
sleep 3

# Verify the process is still running
if ! ps -p "$PID" > /dev/null 2>&1; then
    print_error "Service failed to start. Check logs: $LOG_FILE"
    rm -f "$PID_FILE"
    exit 1
fi

print_info "Next.js interface started successfully (PID: $PID)"
print_info "Service is running on: http://localhost:$PORT"
print_info "Logs are being written to: $LOG_FILE"
print_info "To stop the service, run: ./stop.sh"
print_info "To view logs: tail -f $LOG_FILE"

# Test health endpoint
sleep 2
if curl -s "http://localhost:$PORT" > /dev/null 2>&1; then
    print_info "Health check passed ✓"
else
    print_warning "Health check failed, but service may still be starting..."
fi

