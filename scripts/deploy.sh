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

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}TaskMaster License System Deployment${NC}"
echo -e "${GREEN}Environment: $ENVIRONMENT${NC}"
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

# Copy backend
cp -r "$BACKEND_DIR/server" "$DEPLOY_DIR/"
cp -r "$BACKEND_DIR/migrations" "$DEPLOY_DIR/"
cp -r "$BACKEND_DIR/keys" "$DEPLOY_DIR/" 2>/dev/null || true
mkdir -p "$DEPLOY_DIR/data"

# Copy frontend build
cp -r "$FRONTEND_DIR/.next" "$DEPLOY_DIR/frontend/"
cp -r "$FRONTEND_DIR/public" "$DEPLOY_DIR/frontend/" 2>/dev/null || true
cp "$FRONTEND_DIR/package.json" "$DEPLOY_DIR/frontend/"
cp "$FRONTEND_DIR/next.config.ts" "$DEPLOY_DIR/frontend/"

# Create scripts
cat > "$DEPLOY_DIR/start.sh" << 'EOF'
#!/bin/bash
# Start the production server

# Start backend
./server > backend.log 2>&1 &
BACKEND_PID=$!
echo $BACKEND_PID > .backend.pid

# Start frontend (requires npm/node in production)
cd frontend
npm run start > ../frontend.log 2>&1 &
FRONTEND_PID=$!
echo $FRONTEND_PID > ../.frontend.pid

echo "TaskMaster License System started"
echo "Backend PID: $BACKEND_PID"
echo "Frontend PID: $FRONTEND_PID"
EOF

chmod +x "$DEPLOY_DIR/start.sh"

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

1. Generate root keys (if not already done):
   ```bash
   cd /path/to/backend
   go run cmd/genkeys/main.go root
   # Copy keys/root_public.pem content to ROOT_PUBLIC_KEY in .env
   ```

2. Generate org keys:
   ```bash
   cd /path/to/backend
   go run cmd/genkeys/main.go org your_org_id dev
   ```

3. Configure environment:
   ```bash
   cp .env.example .env
   nano .env  # Edit with your settings
   ```

4. Start the system:
   ```bash
   ./start.sh
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
echo "Total Size: $(du -sh $DEPLOY_DIR | awk '{print $1}')"
echo ""
log_info "Next steps:"
echo "  1. Copy $DEPLOY_DIR to your production server"
echo "  2. Configure .env file"
echo "  3. Generate keys (root and org keys)"
echo "  4. Run ./start.sh"
echo ""
log_info "✓ Deployment package ready!"

