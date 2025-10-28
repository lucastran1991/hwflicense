# Local Testing Instructions

## Start and Test the 7 License Server APIs Locally

### Quick Start

Run the automated script:
```bash
./start_and_test_local.sh
```

### Manual Steps

#### Step 1: Start Mock A-Stack Server (port 8081)

```bash
cd backend
go run cmd/astack-mock/main.go &
cd ..
```

Wait 3 seconds for it to start.

#### Step 2: Start Hub Backend (port 8080)

```bash
cd backend
# Build if needed
go build -o server cmd/server/main.go

# Start server
./server &
cd ..
```

Wait 3 seconds for it to start.

#### Step 3: Run the 7 APIs Test

```bash
./test_license_server_apis.sh http://localhost:8081
```

### What Gets Tested

The test script will test all 7 License Server APIs:

1. **POST /api/v1/sites/create** - Create site key
2. **GET /api/v1/sites** - List site keys  
3. **PUT /api/v1/sites/:id** - Update site key
4. **POST /api/v1/keys/refresh** - Refresh key
5. **POST /api/v1/stats/aggregate** - Send stats
6. **POST /api/v1/keys/validate** - Validate key
7. **POST /api/v1/alerts** - Send alert

### Expected Results

Since the Mock A-Stack server (`backend/cmd/astack-mock/main.go`) doesn't implement all these endpoints, most tests will fail with 404.

However, you can test:
- Health check: `curl http://localhost:8081/health`
- CML Issue: `curl -X POST http://localhost:8081/api/cml/issue ...`
- Manifest Receive: `curl -X POST http://localhost:8081/api/manifests/receive ...`

### Stop Services

```bash
pkill -f "backend/server"
pkill -f "astack-mock"
```

## Notes

The actual License Server (with all 7 APIs) is not yet implemented. The current setup has:
- Hub Backend (port 8080) - Full implementation
- Mock A-Stack (port 8081) - Partial mock for testing
- License Server - Not implemented yet (commented out)

To implement the full License Server with all 7 APIs, you would need to create a separate microservice.

