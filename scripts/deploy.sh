#!/bin/bash

# TaskMaster License System - Production Deployment Script
# This script builds and prepares the system for production deployment
# Usage: ./scripts/deploy.sh [environment]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"
FRONTEND_DIR="$PROJECT_ROOT/frontend"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

ENVIRONMENT="${1:-production}"

# Function to parse JSON config
parse_config() {
    local config_file="$1"
    local key="$2"
    python3 -c "import json, sys; print(json.load(sys.stdin).get('$key', ''))" < "$config_file" 2>/dev/null || echo ""
}

# Try to read mode from config files
BACKEND_MODE=$(parse_config "$PROJECT_ROOT/config/backend.json" "mode" 2>/dev/null || echo "$ENVIRONMENT")
FRONTEND_MODE=$(parse_config "$PROJECT_ROOT/config/frontend.json" "mode" 2>/dev/null || echo "$ENVIRONMENT")

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}TaskMaster License System Deployment${NC}"
echo -e "${GREEN}Environment: $ENVIRONMENT${NC}"
echo -e "${GREEN}Backend Mode: $BACKEND_MODE${NC}"
echo -e "${GREEN}Frontend Mode: $FRONTEND_MODE${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Function to print colored output
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check prerequisites
log_info "Checking prerequisites..."

if ! command -v go &> /dev/null; then
    log_error "Go is not installed. Please install Go 1.21+"
    exit 1
fi

if ! command -v npm &> /dev/null && ! command -v yarn &> /dev/null; then
    log_error "Neither npm nor yarn is installed"
    exit 1
fi

if ! command -v node &> /dev/null; then
    log_error "Node.js is not installed"
    exit 1
fi

log_info "✓ All prerequisites met"
echo ""

# Step 1: Build Backend
log_info "Building backend..."
cd "$BACKEND_DIR"

# Clean previous builds
if [ -f "$BACKEND_DIR/server" ]; then
    log_info "Removing previous backend build..."
    rm -f "$BACKEND_DIR/server"
fi

# Build backend for production
log_info "Compiling Go backend..."
go build -ldflags="-s -w" -o server cmd/server/main.go

if [ -f "$BACKEND_DIR/server" ]; then
    log_info "✓ Backend built successfully"
    ls -lh "$BACKEND_DIR/server"
else
    log_error "Backend build failed"
    exit 1
fi
echo ""

# Step 2: Build Frontend
log_info "Building frontend..."
cd "$FRONTEND_DIR"

# Install dependencies if needed
if [ ! -d "$FRONTEND_DIR/node_modules" ]; then
    log_info "Installing frontend dependencies..."
    npm install
fi

# Build frontend
log_info "Building Next.js frontend..."
npm run build

if [ -d "$FRONTEND_DIR/.next" ]; then
    log_info "✓ Frontend built successfully"
else
    log_error "Frontend build failed"
    exit 1
fi
echo ""

# Step 3: Create deployment package
log_info "Creating deployment package..."

DEPLOY_DIR="$PROJECT_ROOT/deploy"
mkdir -p "$DEPLOY_DIR"

# Copy config folder
cp -r "$PROJECT_ROOT/config" "$DEPLOY_DIR/"

# Copy backend
cp -r "$BACKEND_DIR/server" "$DEPLOY_DIR/"
cp -r "$BACKEND_DIR/migrations" "$DEPLOY_DIR/"
cp -r "$BACKEND_DIR/keys" "$DEPLOY_DIR/" 2>/dev/null || true
mkdir -p "$DEPLOY_DIR/data"

# Copy frontend build
mkdir -p "$DEPLOY_DIR/frontend"
# Remove old .next if exists to ensure clean copy
rm -rf "$DEPLOY_DIR/frontend/.next"
cp -r "$FRONTEND_DIR/.next" "$DEPLOY_DIR/frontend/"
cp -r "$FRONTEND_DIR/public" "$DEPLOY_DIR/frontend/" 2>/dev/null || true
cp "$FRONTEND_DIR/package.json" "$DEPLOY_DIR/frontend/"
cp "$FRONTEND_DIR/next.config.js" "$DEPLOY_DIR/frontend/" 2>/dev/null || true

# Note: node_modules will be installed by wrapper-frontend.sh when starting
# This avoids issues with symlinks and platform-specific binaries

# Copy PM2 files
cp "$PROJECT_ROOT/ecosystem.config.js" "$DEPLOY_DIR/"

# Copy wrapper scripts
cp "$PROJECT_ROOT/wrapper-backend.sh" "$DEPLOY_DIR/" 2>/dev/null || true
cp "$PROJECT_ROOT/wrapper-frontend.sh" "$DEPLOY_DIR/" 2>/dev/null || true

# Make wrappers executable
chmod +x "$DEPLOY_DIR/wrapper-"*.sh 2>/dev/null || true

# Create logs directory
mkdir -p "$DEPLOY_DIR/logs"

# Create scripts
cat > "$DEPLOY_DIR/start.sh" << 'EOF'
#!/bin/bash
# Start the production server

echo "Starting TaskMaster License System..."
echo "Using configurations from config/ folder"
echo ""

# Check if PM2 is available
if command -v pm2 &> /dev/null; then
    echo "PM2 detected - Using PM2 for process management"
    echo ""
    
    # Use PM2 with unified logging
    pm2 start ecosystem.config.js
    pm2 save
    
    echo ""
    echo "====================================="
    echo "TaskMaster License System started with PM2"
    echo "====================================="
    echo "License Server: http://localhost:8081"
    echo "Backend: http://localhost:8080"
    echo "Frontend: http://localhost:3000"
    echo ""
    echo "Unified logs: logs/system.log"
    echo ""
    echo "Useful commands:"
    echo "  pm2 status          - Check status"
    echo "  pm2 logs             - View logs"
    echo "  pm2 monit             - Monitor resources"
    echo "  pm2 restart all      - Restart all"
    echo "  pm2 stop all          - Stop all"
    echo ""
else
    echo "PM2 not found - Using shell scripts"
    echo ""
    
    # Fallback to shell scripts
    mkdir -p logs
    
    # Start backend with wrapper
    ./wrapper-backend.sh >> logs/system.log 2>&1 &
    BACKEND_PID=$!
    echo $BACKEND_PID > .backend.pid
    
    # Start frontend with wrapper
    ./wrapper-frontend.sh >> logs/system.log 2>&1 &
    FRONTEND_PID=$!
    echo $FRONTEND_PID > .frontend.pid
    
    echo "====================================="
    echo "TaskMaster License System started"
    echo "====================================="
    echo "Backend PID: $BACKEND_PID (http://localhost:8080)"
    echo "Frontend PID: $FRONTEND_PID (http://localhost:3000)"
    echo ""
    echo "Unified logs: logs/system.log"
    echo ""
    echo "Note: Install PM2 for better process management:"
    echo "  npm install -g pm2"
    echo ""
fi
EOF

chmod +x "$DEPLOY_DIR/start.sh"

# Create PM2 management script
cat > "$DEPLOY_DIR/pm2.sh" << 'EOF'
#!/bin/bash
# PM2 Management Script
# Usage: ./pm2.sh {start|stop|restart|status|logs|monit}

if ! command -v pm2 &> /dev/null; then
    echo "PM2 is not installed."
    echo "Install it with: npm install -g pm2"
    exit 1
fi

case "$1" in
  start)
    pm2 start ecosystem.config.js
    pm2 save
    ;;
  stop)
    pm2 stop ecosystem.config.js
    ;;
  restart)
    pm2 restart ecosystem.config.js
    ;;
  status)
    pm2 status
    ;;
  logs)
    pm2 logs
    ;;
  monit)
    pm2 monit
    ;;
  *)
    echo "Usage: $0 {start|stop|restart|status|logs|monit}"
    exit 1
    ;;
esac
EOF

chmod +x "$DEPLOY_DIR/pm2.sh"

# Create log viewer script
cat > "$DEPLOY_DIR/view-logs.sh" << 'EOF'
#!/bin/bash
# View unified logs with optional filtering
# Usage: ./view-logs.sh [BE|FE|LS]

LOG_FILE="logs/system.log"

if [ -z "$1" ]; then
  tail -f "$LOG_FILE" 2>/dev/null
else
  if [ "$1" == "BE" ] || [ "$1" == "FE" ] || [ "$1" == "LS" ]; then
    tail -f "$LOG_FILE" 2>/dev/null | grep "\[$1\]"
  else
    echo "Usage: $0 [BE|FE|LS]"
    exit 1
  fi
fi
EOF

chmod +x "$DEPLOY_DIR/view-logs.sh"

# Create environment file template
cat > "$DEPLOY_DIR/.env.example" << 'EOF'
# Database Configuration
DB_PATH=data/taskmaster_license.db

# JWT Configuration
JWT_SECRET=change-this-secret-in-production

# Server Configuration
API_PORT=8080
API_ENV=production

# Encryption Password (for private keys)
ENCRYPTION_PASSWORD=change-this-password-in-production-12345

# A-Stack Mock Server (if running)
ASTACK_MOCK_PORT=8081

# Root Public Key (for CML validation)
ROOT_PUBLIC_KEY=-----BEGIN PUBLIC KEY-----
YOUR_ROOT_PUBLIC_KEY_HERE
-----END PUBLIC KEY-----

# Frontend API URL
NEXT_PUBLIC_API_URL=http://localhost:8080/api
EOF

# Create README for deployment
cat > "$DEPLOY_DIR/README.md" << 'EOF'
# TaskMaster License System - Production Deployment

## Quick Start

### Prerequisites

Install PM2 (recommended for production):
```bash
npm install -g pm2
```

### Configuration

1. Configure the system by editing JSON files in `config/` folder:
   ```bash
   nano config/backend.json
   nano config/frontend.json
   ```

2. Generate root keys (if not already done):
   ```bash
   cd /path/to/backend
   go run cmd/genkeys/main.go root
   # Copy keys/root_public.pem content to config/backend.json
   ```

3. Generate org keys:
   ```bash
   cd /path/to/backend
   go run cmd/genkeys/main.go org your_org_id dev
   ```

4. Set mode to "production" in all config files:
   ```json
   {
     "mode": "production",
     ...
   }
   ```

### Starting the System

**With PM2 (Recommended):**
```bash
./start.sh  # Will detect PM2 and use it
# Or directly:
./pm2.sh start
```

**Without PM2 (Fallback):**
```bash
./start.sh  # Will use shell scripts automatically
```

### Managing the System

**PM2 Commands:**
```bash
./pm2.sh status    # Check status
./pm2.sh restart   # Restart all services
./pm2.sh stop      # Stop all services
./pm2.sh logs      # View logs
./pm2.sh monit     # Monitor resources
```

**Viewing Logs:**
```bash
./view-logs.sh         # All logs
./view-logs.sh BE      # Backend only
./view-logs.sh FE      # Frontend only
./view-logs.sh LS      # License Server only
```

## Environment Variables

See `.env.example` for all available configuration options.

## Data Directory

The system stores data in the `data/` directory:
- SQLite database: `data/taskmaster_license.db`
- Logs: `backend.log`, `frontend.log`

## Backing Up

Important: Back up your `data/` directory regularly:
```bash
tar -czf backup-$(date +%Y%m%d).tar.gz data/
```

## Monitoring

Check logs:
```bash
tail -f backend.log
tail -f frontend.log
```

## Security Notes

1. Change `JWT_SECRET` in production
2. Change `ENCRYPTION_PASSWORD` in production
3. Use strong passwords (16+ characters for encryption)
4. Keep root keys secure
5. Regular backups of database

## Architecture

- Backend: Golang binary (server)
- Frontend: Next.js build (.next)
- Database: SQLite (data/taskmaster_license.db)
- Ports: 8080 (backend), 3000 (frontend)

## Endpoints

- Backend API: http://localhost:8080/api
- Frontend: http://localhost:3000
- Health Check: http://localhost:8080/api/health

## Default Credentials

Username: `admin`
Password: `admin123`

**⚠️ CHANGE THESE IN PRODUCTION!**
EOF

log_info "✓ Deployment package created at: $DEPLOY_DIR"
echo ""

# Summary
log_info "=========================================="
log_info "Deployment Summary"
log_info "=========================================="
echo ""
echo "Backend:"
echo "  Binary: $DEPLOY_DIR/server"
echo "  Size: $(ls -lh $DEPLOY_DIR/server | awk '{print $5}')"
echo ""
echo "Frontend:"
echo "  Build: $DEPLOY_DIR/frontend"
echo "  Size: $(du -sh $DEPLOY_DIR/frontend | awk '{print $1}')"
echo ""
echo "Configuration:"
echo "  Config files: $DEPLOY_DIR/config/"
echo ""
echo "Total Size: $(du -sh $DEPLOY_DIR | awk '{print $1}')"
echo ""
log_info "Next steps:"
echo "  1. Copy $DEPLOY_DIR to your production server"
echo "  2. Edit config files in $DEPLOY_DIR/config/ for your environment"
echo "  3. Set mode to 'production' in all config files"
echo "  4. Generate keys (root and org keys)"
echo "  5. Run ./start.sh"
echo ""
log_info "✓ Deployment package ready!"

