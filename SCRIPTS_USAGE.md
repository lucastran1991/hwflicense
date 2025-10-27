# TaskMaster License System - Management Scripts

## ✅ Scripts Created Successfully!

**Date:** October 27, 2025  
**Location:** `scripts/` directory

---

## 📋 Available Scripts

### 1. Master Management Script (`manage.sh`)
Manages both backend and frontend together.

**Usage:**
```bash
./scripts/manage.sh {start|stop|restart|status|logs}
```

**Examples:**
```bash
# Start everything
./scripts/manage.sh start

# Stop everything
./scripts/manage.sh stop

# Restart everything
./scripts/manage.sh restart

# Check status
./scripts/manage.sh status

# View logs
./scripts/manage.sh logs backend    # Backend logs
./scripts/manage.sh logs frontend   # Frontend logs
```

---

### 2. Backend Script (`backend.sh`)
Manages the backend server only (Golang binary on port 8080).

**Usage:**
```bash
./scripts/backend.sh {start|stop|restart|status|build|logs}
```

**Examples:**
```bash
# Start backend
./scripts/backend.sh start

# Stop backend
./scripts/backend.sh stop

# Restart backend
./scripts/backend.sh restart

# Build backend binary
./scripts/backend.sh build

# Check status
./scripts/backend.sh status

# View logs
./scripts/backend.sh logs
```

---

### 3. Frontend Script (`frontend.sh`)
Manages the frontend server only (Next.js on port 3000).

**Usage:**
```bash
./scripts/frontend.sh {start|stop|restart|status|build|logs}
```

**Examples:**
```bash
# Start frontend
./scripts/frontend.sh start

# Stop frontend
./scripts/frontend.sh stop

# Restart frontend
./scripts/frontend.sh restart

# Build for production
./scripts/frontend.sh build

# Check status
./scripts/frontend.sh status

# View logs
./scripts/frontend.sh logs
```

---

### 4. Deployment Script (`deploy.sh`)
Builds and packages the system for production deployment.

**Usage:**
```bash
./scripts/deploy.sh [environment]
```

**Examples:**
```bash
# Build for production
./scripts/deploy.sh

# Build for staging
./scripts/deploy.sh staging
```

**What it does:**
1. Builds backend binary (optimized, stripped)
2. Builds frontend for production
3. Creates deployment package in `deploy/` directory
4. Includes all necessary files and scripts

---

## 🚀 Quick Start

### Start the Complete System
```bash
./scripts/manage.sh start
```

This will:
1. Start backend on http://localhost:8080
2. Start frontend on http://localhost:3000
3. Log processes IDs to `.backend.pid` and `.frontend.pid`
4. Create log files: `.backend.log` and `.frontend.log`

### Check Status
```bash
./scripts/manage.sh status
```

Output:
```
=== TaskMaster License System Status ===

Backend:  Running (PID: 12345, Port: 8080)
Frontend: Running (PID: 12346, Port: 3000)

Backend logs:  /path/to/.backend.log
Frontend logs: /path/to/.frontend.log
```

### Stop the System
```bash
./scripts/manage.sh stop
```

This will:
1. Stop frontend gracefully
2. Stop backend gracefully
3. Remove PID files
4. Keep log files for debugging

---

## 🏭 Production Deployment

### 1. Build for Production
```bash
./scripts/deploy.sh
```

This creates a `deploy/` directory with:
- ✅ Optimized backend binary
- ✅ Production-ready frontend build
- ✅ Migration files
- ✅ Start script for production
- ✅ Environment file template
- ✅ README with deployment instructions

### 2. Deploy to Server
```bash
# Copy deploy directory to server
scp -r deploy/ user@server:/opt/taskmaster-license/

# SSH to server
ssh user@server

# Navigate to deploy directory
cd /opt/taskmaster-license

# Configure environment
cp .env.example .env
nano .env  # Edit with your settings

# Generate keys (if not done)
go run genkeys/main.go org your_org dev

# Start the system
./start.sh
```

---

## 📝 Individual Service Control

### Backend Only
```bash
# Start just backend
./scripts/backend.sh start

# Stop just backend
./scripts/backend.sh stop

# Restart just backend
./scripts/backend.sh restart

# View backend logs in real-time
./scripts/backend.sh logs
```

### Frontend Only
```bash
# Start just frontend
./scripts/frontend.sh start

# Stop just frontend
./scripts/frontend.sh stop

# Restart just frontend
./scripts/frontend.sh restart

# View frontend logs in real-time
./scripts/frontend.sh logs
```

---

## 🔍 Monitoring & Debugging

### View Logs
```bash
# Backend logs
tail -f .backend.log

# Frontend logs
tail -f .frontend.log

# Or use the scripts
./scripts/manage.sh logs backend
./scripts/manage.sh logs frontend
```

### Check Process Status
```bash
# Using the scripts
./scripts/manage.sh status
./scripts/backend.sh status
./scripts/frontend.sh status

# Or manually
ps aux | grep server
ps aux | grep "npm run dev"
```

---

## 🛠️ Troubleshooting

### Port Already in Use
```bash
# Check what's using the port
lsof -i :8080  # Backend
lsof -i :3000  # Frontend

# Kill the process
kill -9 <PID>
```

### Clear Everything and Start Fresh
```bash
# Stop everything
./scripts/manage.sh stop

# Remove PID and log files
rm -f .backend.pid .frontend.pid .backend.log .frontend.log

# Restart
./scripts/manage.sh start
```

### Build Issues
```bash
# Rebuild backend
./scripts/backend.sh build

# Rebuild frontend
cd frontend
rm -rf node_modules .next
npm install
npm run build
```

---

## 📁 Files Created by Scripts

The scripts create these files (gitignored):

- `.backend.pid` - Backend process ID
- `.frontend.pid` - Frontend process ID
- `.backend.log` - Backend logs
- `.frontend.log` - Frontend logs
- `deploy/` - Production deployment package

---

## ⚙️ Environment Variables

Create a `.env` file in the project root (gitignored):

```bash
# Backend Configuration
DB_PATH=data/taskmaster_license.db
JWT_SECRET=your-secret-key-change-in-production
API_PORT=8080
API_ENV=production
ENCRYPTION_PASSWORD=your-strong-password-16-chars

# Frontend Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

Or set directly:
```bash
export JWT_SECRET=your-secret
export ENCRYPTION_PASSWORD=your-password
./scripts/manage.sh start
```

---

## 🔐 Production Checklist

Before deploying to production:

- [ ] Change `JWT_SECRET` to a strong random value
- [ ] Change `ENCRYPTION_PASSWORD` to a strong password (16+ characters)
- [ ] Generate root keys: `go run cmd/genkeys/main.go root`
- [ ] Generate org keys: `go run cmd/genkeys/main.go org org_id dev`
- [ ] Update `ROOT_PUBLIC_KEY` in `.env` with your root public key
- [ ] Configure firewall (allow ports 8080, 3000, or use reverse proxy)
- [ ] Set up SSL/TLS certificates for HTTPS
- [ ] Configure database backups
- [ ] Change default admin credentials
- [ ] Review security settings

---

## 📚 Script Descriptions

### `manage.sh` (Master Script)
- **Purpose:** Orchestrate all services
- **Commands:** start, stop, restart, status, logs
- **Use When:** You want to control everything at once

### `backend.sh` (Backend Only)
- **Purpose:** Manage Golang backend server
- **Commands:** start, stop, restart, status, build, logs
- **Use When:** You want to restart just the backend

### `frontend.sh` (Frontend Only)
- **Purpose:** Manage Next.js frontend server
- **Commands:** start, stop, restart, status, build, logs
- **Use When:** You want to restart just the frontend

### `deploy.sh` (Deployment)
- **Purpose:** Build for production
- **Commands:** Just run it: `./scripts/deploy.sh`
- **Use When:** Preparing for production deployment

---

## 🎯 Common Workflows

### Development Workflow
```bash
# Start everything
./scripts/manage.sh start

# Make changes to backend
cd backend
# edit files...

# Rebuild and restart backend
./scripts/backend.sh restart

# Check logs if something goes wrong
./scripts/backend.sh logs
```

### Production Deployment
```bash
# Build for production
./scripts/deploy.sh

# Copy deploy/ to server
scp -r deploy/ user@server:/opt/taskmaster-license/

# SSH to server and start
ssh user@server
cd /opt/taskmaster-license
./start.sh
```

### Debugging
```bash
# Check what's running
./scripts/manage.sh status

# View logs
./scripts/manage.sh logs backend
./scripts/manage.sh logs frontend

# Restart if needed
./scripts/manage.sh restart
```

---

## ✨ Summary

You now have **4 comprehensive scripts** to manage your system:

1. **`manage.sh`** - Master script (start/stop everything)
2. **`backend.sh`** - Backend management
3. **`frontend.sh`** - Frontend management
4. **`deploy.sh`** - Production deployment

**All scripts are executable and ready to use!** 🚀

