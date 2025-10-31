#!/bin/bash

# KMS Service Stop Script
# This script stops the KMS service gracefully

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Configuration
PID_FILE="$SCRIPT_DIR/kms.pid"

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

# Try to find PID from file first
PID=""
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    # Verify the PID file contains a valid PID
    if ! ps -p "$PID" > /dev/null 2>&1; then
        print_warning "PID file exists but process $PID is not running. Removing stale PID file."
        rm -f "$PID_FILE"
        PID=""
    fi
fi

# If no valid PID from file, try to find process by name or port
if [ -z "$PID" ]; then
    # Try offspring process by name
    FOUND_PID=$(pgrep -f "kms-server" | head -1)
    if [ -n "$FOUND_PID" ]; then
        PID="$FOUND_PID"
        print_info "Found KMS service process by name (PID: $PID)"
    else
        # Try to find process using port 8080 (default port)
        DEFAULT_PORT="8080"
        if command -v lsof > /dev/null 2>&1; then
            PORT_PID=$(lsof -ti :$DEFAULT_PORT 2>/dev/null | head -1)
            if [ -n "$PORT_PID" ] && ps -p "$PORT_PID" > /dev/null 2>&1; then
                # Check if it's actually kms-server
                if ps -p "$PORT_PID" -o comm= | grep -q "kms-server"; then
                    PID="$PORT_PID"
                    print_info "Found KMS service process on port $DEFAULT_PORT (PID: $PID)"
                fi
            fi
        fi
    fi
fi

# If still no PID found, exit
if [ -z "$PID" ]; then
    print_warning "KMS service is not running (no process found)."
    exit 0
fi

print_info "Stopping KMS service (PID: $PID)..."

# Try graceful shutdown with SIGTERM
kill -TERM "$PID" 2>/dev/null || {
    print_error "Failed to send SIGTERM to process $PID"
    exit 1
}

# Wait for graceful shutdown (max 10 seconds)
MAX_WAIT=10
WAITED=0
while [ $WAITED -lt $MAX_WAIT ]; do
    if ! ps -p "$PID" > /dev/null 2>&1; then
        break
    fi
    sleep 1
    WAITED=$((WAITED + 1))
done

# If process is still running, force kill
if ps -p "$PID" > /dev/null 2>&1; then
    print_warning "Process did not terminate gracefully. Force killing..."
    kill -KILL "$PID" 2>/dev/null || true
    sleep 1
fi

# Verify process is stopped
if ps -p "$PID" > /dev/null 2>&1; then
    print_error "Failed to stop process $PID"
    exit 1
fi

# Remove PID file if it exists and matches stopped PID
if [ -f "$PID_FILE" ]; then
    FILE_PID=$(cat "$PID_FILE")
    if [ "$FILE_PID" = "$PID" ]; then
        rm -f "$PID_FILE"
    fi
fi

print_info "KMS service stopped successfully"

