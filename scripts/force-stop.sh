#!/bin/bash
# Force Stop All Services
# This script forcefully stops all TaskMaster License System services
# Usage: ./scripts/force-stop.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

echo "========================================="
echo "Force Stopping All Services"
echo "========================================="
echo ""

# Stop PM2 processes
if command -v pm2 &> /dev/null; then
    log_info "Stopping PM2 processes..."
    pm2 stop all 2>/dev/null || true
    pm2 delete all 2>/dev/null || true
    log_info "✓ PM2 processes stopped"
fi

# Kill backend server
log_info "Stopping backend server..."
pkill -f "./server" 2>/dev/null || true
pkill -f "backend/.*server" 2>/dev/null || true
pkill -f "cmd/server/main.go" 2>/dev/null || true
# Kill by port
lsof -ti:8080 | xargs kill -9 2>/dev/null || true
log_info "✓ Backend stopped"

# Kill License Server
log_info "Stopping License Server..."
pkill -f "./license-server" 2>/dev/null || true
pkill -f "license-server/.*license-server" 2>/dev/null || true
pkill -f "cmd/license-server" 2>/dev/null || true
# Kill by port
lsof -ti:8081 | xargs kill -9 2>/dev/null || true
log_info "✓ License Server stopped"

# Kill frontend
log_info "Stopping frontend..."
pkill -f "next start" 2>/dev/null || true
pkill -f "npm run dev" 2>/dev/null || true
pkill -f "npm run start" 2>/dev/null || true
# Kill by port
lsof -ti:3000 | xargs kill -9 2>/dev/null || true
log_info "✓ Frontend stopped"

# Clean up PID files
log_info "Cleaning up PID files..."
rm -f "$PROJECT_ROOT/.backend.pid" 2>/dev/null || true
rm -f "$PROJECT_ROOT/.frontend.pid" 2>/dev/null || true
rm -f "$PROJECT_ROOT/.license_server.pid" 2>/dev/null || true
rm -f "$PROJECT_ROOT/deploy/.backend.pid" 2>/dev/null || true
rm -f "$PROJECT_ROOT/deploy/.frontend.pid" 2>/dev/null || true
rm -f "$PROJECT_ROOT/deploy/.license_server.pid" 2>/dev/null || true
log_info "✓ PID files removed"

# Wait a moment for processes to fully terminate
sleep 1

# Final check
echo ""
log_info "Checking for remaining processes..."

BACKEND_COUNT=$(ps aux | grep -E "./server|backend.*server" | grep -v grep | wc -l | tr -d ' ')
LICENSE_COUNT=$(ps aux | grep -E "./license-server|license-server.*license" | grep -v grep | wc -l | tr -d ' ')
FRONTEND_COUNT=$(ps aux | grep -E "next|npm.*dev|npm.*start" | grep -v grep | wc -l | tr -d ' ')

if [ "$BACKEND_COUNT" -gt 0 ] || [ "$LICENSE_COUNT" -gt 0 ] || [ "$FRONTEND_COUNT" -gt 0 ]; then
    log_warn "Some processes may still be running:"
    
    if [ "$BACKEND_COUNT" -gt 0 ]; then
        log_warn "Backend: $BACKEND_COUNT processes"
        ps aux | grep -E "./server|backend.*server" | grep -v grep
    fi
    
    if [ "$LICENSE_COUNT" -gt 0 ]; then
        log_warn "License Server: $LICENSE_COUNT processes"
        ps aux | grep -E "./license-server|license-server.*license" | grep -v grep
    fi
    
    if [ "$FRONTEND_COUNT" -gt 0 ]; then
        log_warn "Frontend: $FRONTEND_COUNT processes"
        ps aux | grep -E "next|npm.*dev|npm.*start" | grep -v grep
    fi
    
    echo ""
    log_warn "You may need to manually kill these processes:"
    echo "  pkill -9 -f '<process-name>'"
else
    log_info "✓ All processes stopped successfully"
fi

# Check ports
echo ""
log_info "Checking ports..."
if lsof -ti:8080 &> /dev/null; then
    log_warn "Port 8080 still in use"
    lsof -ti:8080 | xargs kill -9 2>/dev/null || true
fi

if lsof -ti:8081 &> /dev/null; then
    log_warn "Port 8081 still in use"
    lsof -ti:8081 | xargs kill -9 2>/dev/null || true
fi

if lsof -ti:3000 &> /dev/null; then
    log_warn "Port 3000 still in use"
    lsof -ti:3000 | xargs kill -9 2>/dev/null || true
fi

echo ""
log_info "========================================="
log_info "Force Stop Complete"
log_info "========================================="

