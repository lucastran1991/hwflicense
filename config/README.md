# Configuration Files

This folder contains centralized configuration files for all services in the TaskMaster License System.

## Files

- `backend.json` - Configuration for the Hub (backend) service
- `license-server.json` - Configuration for the License Server microservice
- `frontend.json` - Configuration for the Next.js frontend

## Changing Configuration

### Development Mode (default)

All services are configured for development mode by default with:
- SQLite databases
- Weak passwords (change for production!)
- Localhost URLs

To run in development:
```bash
./scripts/manage.sh start
```

### Production Mode

To switch to production mode:

1. Update the `mode` field in each JSON file:
   ```json
   {
     "mode": "production",
     ...
   }
   ```

2. Update sensitive values:
   - `jwt_secret` - Use a strong random secret (32+ characters)
   - `encryption_password` - Use a strong password for key encryption

3. Build for production:
   ```bash
   ./scripts/deploy.sh production
   ```

## Environment Variables

The system will use environment variables as fallback if config files are not found. Priority:

1. JSON config files (this folder)
2. Environment variables
3. Default values in code

## Security Notes

**⚠️ IMPORTANT for Production:**

1. Change all secrets and passwords
2. Use strong passwords (32+ characters)
3. Set `mode` to `production` in all configs
4. Use HTTPS in production
5. Backup database regularly

## Config Fields

### backend.json
- `mode` - Environment mode (dev/production)
- `db_path` - Path to SQLite database
- `jwt_secret` - JWT signing secret
- `api_port` - API server port (default: 8080)
- `api_env` - API environment (development/production)
- `encryption_password` - Password for AES key encryption
- `astack_mock_port` - Mock A-Stack server port
- `root_public_key` - Root public key for CML validation

### license-server.json
- `mode` - Environment mode (dev/production)
- `port` - License Server port (default: 8081)
- `database_path` - Path to License Server SQLite database
- `jwt_secret` - JWT secret for validation tokens
- `environment` - Environment type

### frontend.json
- `mode` - Environment mode (dev/production)
- `api_url` - Backend API URL
- `port` - Frontend dev server port

