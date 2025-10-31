#!/bin/bash

# Master Stop Script
# This script stops both Next.js frontend and KMS backend services

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
echo "Stopping KMS Backend and Frontend Services"
echo "=========================================="
echo ""

# Check if directories exist
if [ ! -d "$KMS_DIR" ]; then
    print_warning "KMS directory not found: $KMS_DIR (skipping)"
else
    # Stop Next.js Frontend first
    print_service "Stopping Next.js Frontend..."
    if [ -f "$INTERFACE_DIR/stop.sh" ]; then
        cd "$INTERFACE_DIR"
        if bash "$INTERFACE_DIR/stop.sh"; then
            print_info "Frontend stopped successfully"
        else
            print_warning "Frontend stop script encountered issues"
        fi
    else
        print_warning "Interface stop.sh not found: $INTERFACE_DIR/stop.sh"
    fi
    
    # Wait a moment before stopping backend
    sleep 1
fi

# Stop KMS Backend
if [ ! -d "$KMS_DIR" ]; then
    print_warning "KMS directory not found: $KMS_DIR (skipping)"
else
    print_service "Stopping KMS Backend..."
    if [ -f "$KMS_DIR/stop.sh" ]; then
        cd "$KMS_DIR"
        if bash "$KMS_DIR/stop.sh"; then
            print_info "Backend stopped successfully"
        else
            print_warning "Backend stop script encountered issues"
        fi
    else
        print_warning "KMS stop.sh not found: $KMS_DIR/stop.sh"
    fi
fi

# Summary
echo ""
echo "=========================================="
print_info "All services stopped successfully!"
echo "=========================================="
echo ""

