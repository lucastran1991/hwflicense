#!/bin/bash

# Backend Management Script
# Manages the backend server only
# Usage: ./scripts/backend.sh {start|stop|restart|status|build|logs}

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"
PID_FILE="$PROJECT_ROOT/.backend.pid"
LOG_FILE="$PROJECT_ROOT/.backend.log"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
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

is_running() {
    if [ -f "$PID_FILE" ]; then
        local pid=$(cat "$PID_FILE")
        if ps -p "$pid" > /dev/null 2>&1; then
            return 0
        else
            rm -f "$PID_FILE"
            return 1
        fi
    fi
    return 1
}

start() {
    log_info "Starting backend server..."
    
    if is_running; then
        log_warn "Backend is already running (PID: $(cat $PID_FILE))"
        return 0
    fi
    
    cd "$BACKEND_DIR"
    
    # Build if needed
    if [ ! -f "$BACKEND_DIR/server" ]; then
        log_info "Building backend..."
        go build -o server cmd/server/main.go
    fi
    
    # Create data directory
    mkdir -p data
    
    # Start server
    nohup ./server > "$LOG_FILE" 2>&1 &
    echo $! > "$PID_FILE"
    
    sleep 2
    
    if is_running; then
        log_info "Backend started (PID: $(cat $PID_FILE))"
        log_info "Server: http://localhost:8080"
        log_info "Logs: $LOG_FILE"
    else
        log_error "Backend failed to start. Check logs: $LOG_FILE"
        rm -f "$PID_FILE"
        exit 1
    fi
}

stop() {
    log_info "Stopping backend server..."
    
    if ! is_running; then
        log_warn "Backend is not running"
        rm -f "$PID_FILE"
        return 0
    fi
    
    local pid=$(cat "$PID_FILE")
    kill "$pid" || true
    wait "$pid" 2>/dev/null || true
    rm -f "$PID_FILE"
    
    log_info "Backend stopped"
}

restart() {
    stop
    sleep 2
    start
}

status() {
    if is_running; then
        local pid=$(cat "$PID_FILE")
        echo -e "Backend: ${GREEN}Running${NC} (PID: $pid, Port: 8080)"
        echo "URL: http://localhost:8080"
        echo "Logs: $LOG_FILE"
    else
        echo -e "Backend: ${RED}Stopped${NC}"
    fi
}

build() {
    log_info "Building backend..."
    cd "$BACKEND_DIR"
    go build -o server cmd/server/main.go
    log_info "✓ Build successful"
}

logs() {
    if [ -f "$LOG_FILE" ]; then
        tail -f "$LOG_FILE"
    else
        log_warn "No log file found. Start the server first."
    fi
}

case "$1" in
    start) start ;;
    stop) stop ;;
    restart) restart ;;
    status) status ;;
    build) build ;;
    logs) logs ;;
    *)
        echo "Usage: $0 {start|stop|restart|status|build|logs}"
        exit 1
        ;;
esac

