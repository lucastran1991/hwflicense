# TaskMaster License System - Management Scripts

## Quick Start

### Start Everything
```bash
./scripts/manage.sh start
```

### Stop Everything
```bash
./scripts/manage.sh stop
```

### Restart Everything
```bash
./scripts/manage.sh restart
```

### Check Status
```bash
./scripts/manage.sh status
```

### View Logs
```bash
./scripts/manage.sh logs backend    # Backend logs
./scripts/manage.sh logs frontend   # Frontend logs
```

---

## Individual Scripts

### Backend Only
```bash
./scripts/backend.sh start    # Start backend
./scripts/backend.sh stop     # Stop backend
./scripts/backend.sh restart  # Restart backend
./scripts/backend.sh status   # Check status
./scripts/backend.sh build    # Build backend binary
./scripts/backend.sh logs     # View logs
```

### Frontend Only
```bash
./scripts/frontend.sh start    # Start frontend
./scripts/frontend.sh stop     # Stop frontend
./scripts/frontend.sh restart  # Restart frontend
./scripts/frontend.sh status   # Check status
./scripts/frontend.sh build    # Build for production
./scripts/frontend.sh logs     # View logs
```

---

## Production Deployment

### Build for Production
```bash
./scripts/deploy.sh
```

This will:
1. Build backend binary
2. Build frontend for production
3. Create deployment package in `deploy/` directory

### Deploy to Server
1. Copy `deploy/` directory to your production server
2. Configure `.env` file in the deploy directory
3. Generate keys (if not already done):
   ```bash
   go run genkeys/main.go org your_org_id dev
   ```
4. Start the system:
   ```bash
   cd deploy
   ./start.sh
   ```

---

## Available Commands

### Master Script (`manage.sh`)
- `start` - Start both backend and frontend
- `stop` - Stop both backend and frontend
- `restart` - Restart both services
- `status` - Show status of both services
- `logs {backend|frontend}` - View logs

### Backend Script (`backend.sh`)
- `start` - Start backend server (port 8080)
- `stop` - Stop backend server
- `restart` - Restart backend server
- `status` - Show backend status
- `build` - Build backend binary
- `logs` - View backend logs

### Frontend Script (`frontend.sh`)
- `start` - Start frontend server (port 3000)
- `stop` - Stop frontend server
- `restart` - Restart frontend server
- `status` - Show frontend status
- `build` - Build frontend for production
- `logs` - View frontend logs

### Deployment Script (`deploy.sh`)
- Creates production-ready deployment package
- Builds backend binary
- Builds frontend for production
- Creates deployment directory with all files

---

## Environment Variables

Create `.env` file in the project root:

```bash
# Backend Configuration
DB_PATH=data/taskmaster_license.db
JWT_SECRET=your-secret-key-change-in-production
API_PORT=8080
API_ENV=production
ENCRYPTION_PASSWORD=your-encryption-password-16-chars-min

# Frontend Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

Or use environment variables directly:

```bash
export DB_PATH=data/taskmaster_license.db
export JWT_SECRET=your-secret
export ENCRYPTION_PASSWORD=your-encryption-password-16-chars-min
```

---

## Monitoring

### Check Process Status
```bash
./scripts/manage.sh status
```

### View Logs
```bash
# Backend logs
tail -f .backend.log

# Frontend logs
tail -f .frontend.log

# Both (using manage script)
./scripts/manage.sh logs backend
./scripts/manage.sh logs frontend
```

### Kill Stuck Processes
```bash
# Find and kill backend
ps aux | grep server | grep -v grep | awk '{print $2}' | xargs kill

# Find and kill frontend
ps aux | grep "npm run dev" | grep -v grep | awk '{print $2}' | xargs kill
```

---

## Troubleshooting

### Port Already in Use
```bash
# Check what's using the port
lsof -i :8080  # Backend
lsof -i :3000  # Frontend

# Kill the process
kill -9 <PID>
```

### Build Failures
```bash
# Clean and rebuild backend
cd backend
rm server
go build -o server cmd/server/main.go

# Clean and rebuild frontend
cd frontend
rm -rf .next node_modules
npm install
npm run build
```

### Database Issues
```bash
# Reset database (BACKUP FIRST!)
cd backend
rm data/taskmaster_license.db
# Database will be recreated on next start
```

---

## Production Checklist

Before deploying to production:

- [ ] Change `JWT_SECRET` to a strong random value
- [ ] Change `ENCRYPTION_PASSWORD` to a strong password (16+ chars)
- [ ] Update `ROOT_PUBLIC_KEY` with your actual root public key
- [ ] Generate org keys for your organization
- [ ] Configure firewall rules (allow 8080, 3000)
- [ ] Set up SSL/TLS certificates for HTTPS
- [ ] Configure regular database backups
- [ ] Set up monitoring and logging
- [ ] Change default admin credentials
- [ ] Review and configure security settings

---

## Files Created by Scripts

- `.backend.pid` - Backend process ID
- `.frontend.pid` - Frontend process ID
- `.backend.log` - Backend logs
- `.frontend.log` - Frontend logs
- `deploy/` - Production deployment package

These files are gitignored for cleanliness.

---

## Quick Reference

```bash
# Start everything
./scripts/manage.sh start

# Stop everything
./scripts/manage.sh stop

# Restart everything
./scripts/manage.sh restart

# Check status
./scripts/manage.sh status

# Deploy for production
./scripts/deploy.sh

# Backend only
./scripts/backend.sh start

# Frontend only
./scripts/frontend.sh start
```

