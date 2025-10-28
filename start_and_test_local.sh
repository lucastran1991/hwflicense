#!/bin/bash
# Start all services and test the 7 License Server APIs locally

echo "========================================="
echo "Starting Local Services and Testing"
echo "========================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Stop any existing services
echo "Cleaning up any existing services..."
pkill -f "backend/server" 2>/dev/null || true
pkill -f "astack-mock" 2>/dev/null || true
pkill -f "frontend" 2>/dev/null || true
sleep 1

# Start Mock A-Stack Server (port 8081)
echo -e "${GREEN}Starting Mock A-Stack Server on port 8081...${NC}"
cd backend
go run cmd/astack-mock/main.go > /tmp/astack-mock.log 2>&1 &
ASTACK_PID=$!
cd ..

# Wait for it to start
sleep 3
echo "Mock A-Stack PID: $ASTACK_PID"
echo ""

# Check if it's running
if curl -s http://localhost:8081/health > /dev/null; then
    echo -e "${GREEN}✓ Mock A-Stack is running on port 8081${NC}"
else
    echo -e "${RED}✗ Mock A-Stack failed to start${NC}"
    cat /tmp/astack-mock.log
    exit 1
fi
echo ""

# Start Backend/Hub (port 8080)
echo -e "${GREEN}Starting Hub Backend on port 8080...${NC}"
cd backend
if [ ! -f server ]; then
    echo "Building backend..."
    go build -o server cmd/server/main.go
fi

./server > /tmp/backend.log 2>&1 &
BACKEND_PID=$!
cd ..

# Wait for it to start
sleep 3
echo "Backend PID: $BACKEND_PID"
echo ""

# Check if it's running
if curl -s http://localhost:8080/api/health > /dev/null; then
    echo -e "${GREEN}✓ Hub Backend is running on port 8080${NC}"
else
    echo -e "${RED}✗ Hub Backend failed to start${NC}"
    cat /tmp/backend.log
fi
echo ""

echo "========================================="
echo "Running 7 License Server APIs Test"
echo "========================================="
echo ""

# Run the test script
./test_license_server_apis.sh http://localhost:8081

# Capture exit code
TEST_RESULT=$?

echo ""
echo "========================================="
echo "Test Complete"
echo "========================================="
echo ""
echo "Services are still running in the background."
echo "To stop them, run:"
echo "  kill $BACKEND_PID $ASTACK_PID"
echo ""

exit $TEST_RESULT

