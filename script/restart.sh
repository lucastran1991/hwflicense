#!/bin/bash

# Master Restart Script
# This script stops and starts both KMS backend and Next.js frontend services

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

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

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_service() {
    echo -e "${BLUE}[SERVICE]${NC} $1"
}

echo ""
echo "=========================================="
echo "Restarting KMS Backend and Frontend Services"
echo "=========================================="
echo ""

print_service "Restarting services..."

# Stop services
print_info "Stopping services..."
if [ -f "$SCRIPT_DIR/stop.sh" ]; then
    bash "$SCRIPT_DIR/stop.sh"
else
    print_warning "stop.sh not found, trying to stop manually..."
fi

# Wait for complete shutdown
print_info "Waiting for services to shut down completely..."
sleep 3

# Start services
print_info "Starting services..."
if [ -f "$SCRIPT_DIR/start.sh" ]; then
    bash "$SCRIPT_DIR/start.sh"
else
    print_warning "start.sh not found"
    exit 1
fi

print_info "Restart completed"
echo ""

