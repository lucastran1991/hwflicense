#!/bin/bash

# KMS Service Start Script
# This script builds and starts the KMS service

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

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

# Configuration
PID_FILE="$SCRIPT_DIR/kms.pid"
LOG_FILE="$SCRIPT_DIR/kms.log"
BINARY_NAME="kms-server"
BINARY_PATH="$SCRIPT_DIR/$BINARY_NAME"

# Resolve main.go path - try multiple locations with better error handling
MAIN_PATH="$SCRIPT_DIR/cmd/server/main.go"

# Verify directory structure exists first
if [ ! -d "$SCRIPT_DIR/cmd" ]; then
    print_error "cmd/ directory not found at: $SCRIPT_DIR/cmd"
    print_error "Script directory: $SCRIPT_DIR"
    print_error "Current directory: $(pwd)"
    print_error "Available directories:"
    ls -la "$SCRIPT_DIR" 2>&1 | head -20 || echo "Cannot list directory"
    exit 1
fi

if [ ! -d "$SCRIPT_DIR/cmd/server" ]; then
    print_error "cmd/server/ directory not found at: $SCRIPT_DIR/cmd/server"
    print_error "Available in cmd/:"
    ls -la "$SCRIPT_DIR/cmd" 2>&1 || echo "Cannot list cmd directory"
    exit 1
fi

if [ ! -f "$MAIN_PATH" ]; then
    # Try alternative paths
    if [ -f "cmd/server/main.go" ]; then
        MAIN_PATH="cmd/server/main.go"
    elif [ -f "./cmd/server/main.go" ]; then
        MAIN_PATH="./cmd/server/main.go"
    else
        print_error "Cannot find main.go at: $SCRIPT_DIR/cmd/server/main.go"
        print_error "Current directory: $(pwd)"
        print_error "Script directory: $SCRIPT_DIR"
        print_error "Main path attempted: $MAIN_PATH"
        print_error "Available files in cmd/server/:"
        ls -la "$SCRIPT_DIR/cmd/server" 2>&1 || echo "Cannot list cmd/server directory"
        exit 1
    fi
fi

# Check if service is already running
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p "$PID" > /dev/null 2>&1; then
        print_warning "KMS service is already running (PID: $PID)"
        exit 1
    else
        print_info "Removing stale PID file"
        rm -f "$PID_FILE"
    fi
fi

# Check for master key - try environment variable first, then secure file
MASTER_KEY_FILE="$SCRIPT_DIR/secrets/master.key"
KMS_MASTER_KEY_SET=0

if [ -n "$KMS_MASTER_KEY" ]; then
    # Use environment variable if set
    KMS_MASTER_KEY_SET=1
    print_info "Using master key from environment variable"
elif [ -f "$MASTER_KEY_FILE" ]; then
    # Load from secure file
    export KMS_MASTER_KEY="$(cat "$MASTER_KEY_FILE")"
    KMS_MASTER_KEY_SET=1
    print_info "Loaded master key from secure file: $MASTER_KEY_FILE"
else
    # Generate new master key and save to file
    print_info "Master key not found. Generating new master key..."
    
    # Create secrets directory if it doesn't exist
    mkdir -p "$SCRIPT_DIR/secrets"
    
    # Generate master key
    NEW_MASTER_KEY=$(openssl rand -base64 32)
    
    # Save to file with secure permissions (600 - owner read/write only)
    echo -n "$NEW_MASTER_KEY" > "$MASTER_KEY_FILE"
    chmod 600 "$MASTER_KEY_FILE"
    
    export KMS_MASTER_KEY="$NEW_MASTER_KEY"
    KMS_MASTER_KEY_SET=1
    
    print_info "Generated and saved master key to: $MASTER_KEY_FILE (chmod 600)"
    print_warning "Keep this file secure! Do not commit it to version control."
fi

if [ "$KMS_MASTER_KEY_SET" -eq 0 ]; then
    print_error "Failed to set master key"
    exit 1
fi

# Load default values from setting.json or use defaults
# Note: Go code will load from setting.json if env vars not set
# We only set them here if explicitly provided
if [ -z "$KMS_DB_PATH" ]; then
    # Don't set default - let Go code use setting.json
    unset KMS_DB_PATH
fi
if [ -z "$KMS_PORT" ]; then
    # Don't set default - let Go code use setting.json
    unset KMS_PORT
fi

print_info "Building KMS service..."
print_info "Main path: $MAIN_PATH"
print_info "Binary path: $BINARY_PATH"
print_info "Working directory: $(pwd)"

# Verify main.go exists before building
if [ ! -f "$MAIN_PATH" ]; then
    print_error "main.go not found at: $MAIN_PATH"
    print_error "Current directory: $(pwd)"
    print_error "Please check that cmd/server/main.go exists"
    exit 1
fi

if ! go build -o "$BINARY_PATH" "$MAIN_PATH"; then
    print_error "Failed to build KMS service"
    print_error "Check that Go is installed and working: $(which go)"
    print_error "Verify the main.go file is valid"
    exit 1
fi
print_info "Build successful"

print_info "Starting KMS service..."
print_info "Configuration:"
if [ -n "$KMS_DB_PATH" ]; then
    print_info "  DB Path: $KMS_DB_PATH (from env)"
else
    print_info "  DB Path: (will load from setting.json)"
fi
if [ -n "$KMS_PORT" ]; then
    print_info "  Port: $KMS_PORT (from env)"
else
    print_info "  Port: (will load from setting.json)"
fi
print_info "  Log File: $LOG_FILE"

# Start the service in background
nohup "$BINARY_PATH" > "$LOG_FILE" 2>&1 &
PID=$!

# Save PID
echo $PID > "$PID_FILE"

# Wait a moment to check if service started successfully
sleep 2

# Verify the process is still running
if ! ps -p "$PID" > /dev/null 2>&1; then
    print_error "Service failed to start. Check logs: $LOG_FILE"
    rm -f "$PID_FILE"
    exit 1
fi

print_info "KMS service started successfully (PID: $PID)"
print_info "Logs are being written to: $LOG_FILE"
print_info "To stop the service, run: ./stop.sh"
print_info "To view logs: tail -f $LOG_FILE"

# Test health endpoint
sleep 1
# Normalize port format (ensure it has colon prefix)
HEALTH_PORT="$KMS_PORT"
if [[ ! "$HEALTH_PORT" =~ ^: ]]; then
    HEALTH_PORT=":$HEALTH_PORT"
fi
if curl -s http://localhost${HEALTH_PORT}/health > /dev/null 2>&1; then
    print_info "Health check passed ✓"
else
    print_warning "Health check failed, but service may still be starting..."
fi

