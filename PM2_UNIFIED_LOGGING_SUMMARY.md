# PM2 Integration and Unified Logging - Implementation Summary

## Overview

Successfully implemented PM2 process management and unified logging with service tags for the TaskMaster License System.

## What Was Implemented

### 1. PM2 Configuration (`ecosystem.config.js`)
- Centralized PM2 configuration for all 3 services
- Unified logging to single `logs/system.log` file
- Service tags: `[LS]` for License Server, `[BE]` for Backend, `[FE]` for Frontend
- Auto-restart on crash
- Memory limits and monitoring

### 2. Log Wrapper Scripts
Created wrapper scripts that add service tags to all logs:
- `wrapper-license-server.sh` - Tags logs with `[LS]`
- `wrapper-backend.sh` - Tags logs with `[BE]`
- `wrapper-frontend.sh` - Tags logs with `[FE]`

Each wrapper adds timestamp and service tag to log output.

### 3. Management Scripts
- **`scripts/pm2-manage.sh`** - PM2 management script with commands:
  - `start` - Start all services
  - `stop` - Stop all services
  - `restart` - Restart all services
  - `status` - Show status
  - `logs` - View unified logs
  - `monit` - Monitor resources

- **`scripts/view-logs.sh`** - Log viewer with filtering:
  - `./view-logs.sh` - All logs with colors
  - `./view-logs.sh BE` - Backend only
  - `./view-logs.sh FE` - Frontend only
  - `./view-logs.sh LS` - License Server only

### 4. Updated Deployment Script
Modified `scripts/deploy.sh` to:
- Copy PM2 configuration files
- Copy wrapper scripts
- Create logs directory
- Update `start.sh` to auto-detect PM2 and use it if available
- Fallback to shell scripts if PM2 not installed
- Create `pm2.sh` and `view-logs.sh` in deployment package

## Benefits

### Process Management
- ✅ **Auto-restart**: Services restart automatically on crash
- ✅ **Monitoring**: Built-in monitoring with `pm2 monit`
- ✅ **Production-ready**: Industry-standard process manager
- ✅ **Zero-downtime**: Reload without downtime
- ✅ **Startup scripts**: Auto-start on system boot

### Unified Logging
- ✅ **Single file**: All logs in `logs/system.log`
- ✅ **Service tags**: Easy to identify log source
- ✅ **Timestamps**: All logs have consistent timestamps
- ✅ **Filtering**: Easy to filter by service
- ✅ **Color coding**: Optional color tags for readability

## Usage

### Development
```bash
# Start with PM2
./scripts/pm2-manage.sh start

# View logs
./scripts/view-logs.sh        # All logs
./scripts/view-logs.sh BE     # Backend only
```

### Production Deployment
```bash
# Deploy
./scripts/deploy.sh

# On server - Install PM2 (if not installed)
npm install -g pm2

# Start with PM2
./start.sh  # Auto-detects PM2

# Or use PM2 commands directly
./pm2.sh start
./pm2.sh status
./pm2.sh logs
./pm2.sh monit
```

### Log Management
```bash
# View unified logs
tail -f logs/system.log

# View with service filtering
./view-logs.sh BE    # Backend logs only
./view-logs.sh FE    # Frontend logs only
./view-logs.sh LS    # License Server only

# Search logs
grep "ERROR" logs/system.log
grep "\[BE\]" logs/system.log
```

## Architecture

```
┌─────────────────────────────────────┐
│       PM2 Process Manager           │
├─────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌────┐│
│  │  [LS]    │  │  [BE]    │  │[FE]││
│  │Wrapper   │  │Wrapper   │  │Wrap││
│  └────┬─────┘  └────┬─────┘  └─┬──┘│
│       │             │           │   │
│       └─────────────┼───────────┘   │
│                     │               │
└─────────────────────┼───────────────┘
                     │
              ┌──────▼────────┐
              │ system.log   │
              │ [LS] logs    │
              │ [BE] logs    │
              │ [FE] logs    │
              └──────────────┘
```

## Files Created

1. `ecosystem.config.js` - PM2 configuration
2. `wrapper-license-server.sh` - License Server wrapper
3. `wrapper-backend.sh` - Backend wrapper
4. `wrapper-frontend.sh` - Frontend wrapper
5. `scripts/pm2-manage.sh` - PM2 management script
6. `scripts/view-logs.sh` - Log viewer with filtering

## Files Modified

1. `scripts/deploy.sh` - Added PM2 support
2. Deployment README - Updated with PM2 instructions

## Fallback Support

The system gracefully falls back to shell scripts if PM2 is not installed:
- `start.sh` detects PM2 availability
- Uses PM2 if available (recommended)
- Falls back to shell scripts if PM2 not found
- Both methods support unified logging

## Log Format

All logs follow this format:
```
[SERVICE_TAG] YYYY-MM-DD HH:MM:SS Log message
```

Example:
```
[LS] 2024-10-28 13:45:23 License Server started on port 8081
[BE] 2024-10-28 13:45:24 Backend started on port 8080
[FE] 2024-10-28 13:45:25 Frontend started on port 3000
```

## Next Steps

1. Test PM2 with full system
2. Set up PM2 startup script for production
3. Configure log rotation (PM2 includes this)
4. Monitor with PM2 dashboard
5. Set up alerts for critical errors

## References

- PM2 Documentation: https://pm2.keymetrics.io/docs/
- PM2 Startup Guide: https://pm2.keymetrics.io/docs/usage/startup/
- Ecosystem File Reference: https://pm2.keymetrics.io/docs/usage/application-declaration/

