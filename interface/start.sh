#!/bin/bash

# Next.js Interface Start Script
# This script builds and starts the Next.js interface

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Configuration
PID_FILE="$SCRIPT_DIR/interface.pid"
LOG_FILE="$SCRIPT_DIR/interface.log"
PORT="${PORT:-3000}"

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

# Build the application
print_info "Building Next.js application..."
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

