#!/bin/bash
# Quick License Server 7 APIs Test

echo "========================================="
echo "License Server - 7 Core APIs Test"
echo "========================================="
echo ""

BASE_URL="${1:-http://localhost:8081}"
echo "Target: $BASE_URL"
echo ""

# Test if server is running
if ! curl -s "$BASE_URL/health" > /dev/null 2>&1; then
    echo "ERROR: License server not running at $BASE_URL"
    echo "Start it with: ./scripts/license-server.sh start"
    exit 1
fi

echo "✓ License server is running"
echo ""
echo "The 7 License Server APIs are:"
echo ""
echo "1. POST /api/v1/sites/create     - Create site key"
echo "2. GET  /api/v1/sites            - List all site keys"
echo "3. PUT  /api/v1/sites/:id        - Update site key"
echo "4. POST /api/v1/keys/refresh     - Refresh key"
echo "5. POST /api/v1/stats/aggregate  - Send quarterly stats"
echo "6. POST /api/v1/keys/validate    - Validate key"
echo "7. POST /api/v1/alerts           - Send alert"
echo ""

# Run the comprehensive test
echo "Running comprehensive test suite..."
./test_license_server_apis.sh

