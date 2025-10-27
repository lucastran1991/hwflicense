#!/bin/bash

# Frontend Management Script
# Manages the frontend server only
# Usage: ./scripts/frontend.sh {start|stop|restart|status|build|logs}

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
FRONTEND_DIR="$PROJECT_ROOT/frontend"
PID_FILE="$PROJECT_ROOT/.frontend.pid"
LOG_FILE="$PROJECT_ROOT/.frontend.log"

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
    log_info "Starting frontend server..."
    
    if is_running; then
        log_warn "Frontend is already running (PID: $(cat $PID_FILE))"
        return 0
    fi
    
    cd "$FRONTEND_DIR"
    
    # Install dependencies if needed
    if [ ! -d "$FRONTEND_DIR/node_modules" ]; then
        log_info "Installing dependencies..."
        npm install
    fi
    
    # Start frontend
    nohup npm run dev > "$LOG_FILE" 2>&1 &
    echo $! > "$PID_FILE"
    
    sleep 5
    
    if is_running; then
        log_info "Frontend started (PID: $(cat $PID_FILE))"
        log_info "Server: http://localhost:3000"
        log_info "Logs: $LOG_FILE"
    else
        log_error "Frontend failed to start. Check logs: $LOG_FILE"
        rm -f "$PID_FILE"
        exit 1
    fi
}

stop() {
    log_info "Stopping frontend server..."
    
    if ! is_running; then
        log_warn "Frontend is not running"
        rm -f "$PID_FILE"
        return 0
    fi
    
    local pid=$(cat "$PID_FILE")
    kill "$pid" || true
    wait "$pid" 2>/dev/null || true
    rm -f "$PID_FILE"
    
    log_info "Frontend stopped"
}

restart() {
    stop
    sleep 2
    start
}

status() {
    if is_running; then
        local pid=$(cat "$PID_FILE")
        echo -e "Frontend: ${GREEN}Running${NC} (PID: $pid, Port: 3000)"
        echo "URL: http://localhost:3000"
        echo "Logs: $LOG_FILE"
    else
        echo -e "Frontend: ${RED}Stopped${NC}"
    fi
}

build() {
    log_info "Building frontend for production..."
    cd "$FRONTEND_DIR"
    
    if [ ! -d "$FRONTEND_DIR/node_modules" ]; then
        log_info "Installing dependencies..."
        npm install
    fi
    
    npm run build
    log_info "✓ Build successful"
    log_info "Output: $FRONTEND_DIR/.next"
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

