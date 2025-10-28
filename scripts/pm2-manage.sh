#!/bin/bash
# PM2 Management Script
# Usage: ./scripts/pm2-manage.sh {start|stop|restart|status|logs|monit}

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Check if PM2 is installed
if ! command -v pm2 &> /dev/null; then
    echo "PM2 is not installed. Installing..."
    npm install -g pm2
fi

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

cd "$PROJECT_ROOT"

case "$1" in
  start)
    log_info "Starting services with PM2..."
    pm2 start ecosystem.config.js
    pm2 save
    ;;
  stop)
    log_info "Stopping services..."
    pm2 stop ecosystem.config.js
    ;;
  restart)
    log_info "Restarting services..."
    pm2 restart ecosystem.config.js
    ;;
  status)
    pm2 status
    ;;
  logs)
    pm2 logs --lines 100
    ;;
  monit)
    pm2 monit
    ;;
  *)
    echo "Usage: $0 {start|stop|restart|status|logs|monit}"
    echo ""
    echo "Commands:"
    echo "  start   - Start all services"
    echo "  stop    - Stop all services"
    echo "  restart - Restart all services"
    echo "  status  - Show service status"
    echo "  logs    - Show unified logs"
    echo "  monit   - Open monitoring dashboard"
    exit 1
    ;;
esac

