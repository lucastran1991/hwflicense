#!/bin/bash

# License Server Management Script
# Manages the License Server microservice
# Usage: ./scripts/license-server.sh {start|stop|restart|status|build|logs}

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
LICENSE_SERVER_DIR="$PROJECT_ROOT/license-server"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

PID_FILE="$PROJECT_ROOT/.license_server.pid"

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

start_license_server() {
    log_info "Starting License Server..."
    
    if [ -f "$PID_FILE" ]; then
        local pid=$(cat "$PID_FILE")
        if ps -p "$pid" > /dev/null 2>&1; then
            log_warn "License Server is already running (PID: $pid)"
            return 0
        fi
    fi
    
    cd "$LICENSE_SERVER_DIR"
    
    # Build if needed
    if [ ! -f "license-server" ]; then
        log_info "Building License Server..."
        go build -o license-server ./cmd/license-server
    fi
    
    # Start license server in background
    nohup ./license-server > "$PROJECT_ROOT/.license_server.log" 2>&1 &
    LICENSE_SERVER_PID=$!
    echo $LICENSE_SERVER_PID > "$PID_FILE"
    
    log_info "License Server started (PID: $LICENSE_SERVER_PID)"
    log_info "Logs: $PROJECT_ROOT/.license_server.log"
    log_info "License Server is running on http://localhost:8081"
}

stop_license_server() {
    log_info "Stopping License Server..."
    
    if [ -f "$PID_FILE" ]; then
        local pid=$(cat "$PID_FILE")
        if ps -p "$pid" > /dev/null 2>&1; then
            kill "$pid" || true
            wait "$pid" 2>/dev/null || true
            log_info "License Server stopped"
        else
            log_warn "License Server was not running"
        fi
        rm -f "$PID_FILE"
    else
        log_warn "License Server PID file not found"
    fi
}

show_status() {
    if [ -f "$PID_FILE" ]; then
        local pid=$(cat "$PID_FILE")
        if ps -p "$pid" > /dev/null 2>&1; then
            log_info "License Server is running (PID: $pid)"
            return 0
        else
            rm -f "$PID_FILE"
        fi
    fi
    log_warn "License Server is not running"
}

show_logs() {
    if [ -f "$PROJECT_ROOT/.license_server.log" ]; then
        tail -f "$PROJECT_ROOT/.license_server.log"
    else
        log_error "License Server log file not found"
    fi
}

build_license_server() {
    log_info "Building License Server..."
    cd "$LICENSE_SERVER_DIR"
    go build -o license-server ./cmd/license-server
    log_info "License Server built successfully"
}

# Main command handler
case "${1:-}" in
    start)
        start_license_server
        ;;
    stop)
        stop_license_server
        ;;
    restart)
        stop_license_server
        sleep 2
        start_license_server
        ;;
    status)
        show_status
        ;;
    build)
        build_license_server
        ;;
    logs)
        show_logs
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status|build|logs}"
        exit 1
        ;;
esac

