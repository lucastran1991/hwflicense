#!/bin/bash

# KMS Service Restart Script
# This script stops and starts the KMS service
# Works even if service is not currently running (will just start it)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_info "Restarting KMS service..."

# Stop the service
if [ -f "$SCRIPT_DIR/stop.sh" ]; then
    bash "$SCRIPT_DIR/stop.sh"
else
    print_warning "stop.sh not found, trying to stop manually..."
    PID_FILE="$SCRIPT_DIR/kms.pid"
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if ps -p "$PID" > /dev/null 2>&1; then
            kill -TERM "$PID" 2>/dev/null || true
            sleep 2
            if ps -p "$PID" > /dev/null 2>&1; then
                kill -KILL "$PID" 2>/dev/null || true
            fi
            rm -f "$PID_FILE"
        fi
    fi
fi

# Wait a moment before starting
sleep 1

# Start the service
if [ -f "$SCRIPT_DIR/start.sh" ]; then
    bash "$SCRIPT_DIR/start.sh"
else
    print_warning "start.sh not found"
    exit 1
fi

print_info "Restart completed"

