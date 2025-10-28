#!/bin/bash

# TaskMaster License System Management Script
# This script manages the backend server, license server, and frontend
# Usage: ./scripts/manage.sh {start|stop|restart|status|logs}

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"
FRONTEND_DIR="$PROJECT_ROOT/frontend"
LICENSE_SERVER_DIR="$PROJECT_ROOT/license-server"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# PIDs
BACKEND_PID_FILE="$PROJECT_ROOT/.backend.pid"
FRONTEND_PID_FILE="$PROJECT_ROOT/.frontend.pid"
LICENSE_SERVER_PID_FILE="$PROJECT_ROOT/.license_server.pid"

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_pid_file() {
    local pid_file=$1
    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if ps -p "$pid" > /dev/null 2>&1; then
            return 0
        else
            rm -f "$pid_file"
            return 1
        fi
    fi
    return 1
}

start_backend() {
    log_info "Starting backend server..."
    
    if check_pid_file "$BACKEND_PID_FILE"; then
        log_warn "Backend is already running (PID: $(cat $BACKEND_PID_FILE))"
        return 0
    fi
    
    cd "$BACKEND_DIR"
    
    # Create data directory if it doesn't exist
    mkdir -p data
    
    # Build if binary doesn't exist
    if [ ! -f "$BACKEND_DIR/server" ]; then
        log_info "Building backend..."
        go build -o server cmd/server/main.go
    fi
    
    # Start backend in background
    nohup ./server > "$PROJECT_ROOT/.backend.log" 2>&1 &
    BACKEND_PID=$!
    echo $BACKEND_PID > "$BACKEND_PID_FILE"
    
    log_info "Backend started (PID: $BACKEND_PID)"
    log_info "Logs: $PROJECT_ROOT/.backend.log"
    sleep 2
    
    # Check if it's actually running
    if ! ps -p $BACKEND_PID > /dev/null 2>&1; then
        log_error "Backend failed to start. Check logs for details."
        rm -f "$BACKEND_PID_FILE"
        return 1
    fi
    
    log_info "Backend is running on http://localhost:8080"
}

start_frontend() {
    log_info "Starting frontend server..."
    
    if check_pid_file "$FRONTEND_PID_FILE"; then
        log_warn "Frontend is already running (PID: $(cat $FRONTEND_PID_FILE))"
        return 0
    fi
    
    cd "$FRONTEND_DIR"
    
    # Check if node_modules exists
    if [ ! -d "$FRONTEND_DIR/node_modules" ]; then
        log_info "Installing frontend dependencies..."
        npm install
    fi
    
    # Start frontend in background
    nohup npm run dev > "$PROJECT_ROOT/.frontend.log" 2>&1 &
    FRONTEND_PID=$!
    echo $FRONTEND_PID > "$FRONTEND_PID_FILE"
    
    log_info "Frontend started (PID: $FRONTEND_PID)"
    log_info "Logs: $PROJECT_ROOT/.frontend.log"
    
    # Wait a bit for it to start
    sleep 5
    
    log_info "Frontend is running on http://localhost:3000"
}

stop_backend() {
    log_info "Stopping backend server..."
    
    if [ -f "$BACKEND_PID_FILE" ]; then
        BACKEND_PID=$(cat "$BACKEND_PID_FILE")
        if ps -p "$BACKEND_PID" > /dev/null 2>&1; then
            kill "$BACKEND_PID" || true
            wait "$BACKEND_PID" 2>/dev/null || true
            log_info "Backend stopped"
        else
            log_warn "Backend was not running"
        fi
        rm -f "$BACKEND_PID_FILE"
    else
        log_warn "Backend PID file not found"
    fi
}

stop_frontend() {
    log_info "Stopping frontend server..."
    
    if [ -f "$FRONTEND_PID_FILE" ]; then
        FRONTEND_PID=$(cat "$FRONTEND_PID_FILE")
        if ps -p "$FRONTEND_PID" > /dev/null 2>&1; then
            kill "$FRONTEND_PID" || true
            wait "$FRONTEND_PID" 2>/dev/null || true
            log_info "Frontend stopped"
        else
            log_warn "Frontend was not running"
        fi
        rm -f "$FRONTEND_PID_FILE"
    else
        log_warn "Frontend PID file not found"
    fi
}

status() {
    echo "=== TaskMaster License System Status ==="
    echo ""
    
    if check_pid_file "$BACKEND_PID_FILE"; then
        BACKEND_PID=$(cat "$BACKEND_PID_FILE")
        echo -e "Backend:  ${GREEN}Running${NC} (PID: $BACKEND_PID, Port: 8080)"
    else
        echo -e "Backend:  ${RED}Stopped${NC}"
    fi
    
    if check_pid_file "$FRONTEND_PID_FILE"; then
        FRONTEND_PID=$(cat "$FRONTEND_PID_FILE")
        echo -e "Frontend: ${GREEN}Running${NC} (PID: $FRONTEND_PID, Port: 3000)"
    else
        echo -e "Frontend: ${RED}Stopped${NC}"
    fi
    
    echo ""
    echo "Backend logs:  $PROJECT_ROOT/.backend.log"
    echo "Frontend logs: $PROJECT_ROOT/.frontend.log"
}

logs() {
    if [ "$1" == "backend" ]; then
        tail -f "$PROJECT_ROOT/.backend.log"
    elif [ "$1" == "frontend" ]; then
        tail -f "$PROJECT_ROOT/.frontend.log"
    else
        log_error "Usage: $0 logs {backend|frontend}"
        exit 1
    fi
}

# Main command handling
case "$1" in
    start)
        start_backend
        start_frontend
        echo ""
        status
        ;;
    stop)
        stop_frontend
        stop_backend
        ;;
    restart)
        stop_frontend
        stop_backend
        sleep 2
        start_backend
        start_frontend
        echo ""
        status
        ;;
    status)
        status
        ;;
    logs)
        if [ -z "$2" ]; then
            log_error "Usage: $0 logs {backend|frontend}"
            exit 1
        fi
        logs "$2"
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status|logs {backend|frontend}}"
        exit 1
        ;;
esac

