#!/bin/bash
# Test CML Upload and Site Creation Flow

BASE_URL="http://localhost:8080"

# Get auth token
echo "1. Getting auth token..."
TOKEN=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.token')

echo "Token: ${TOKEN:0:50}..."
echo ""

# Test CML upload (if it exists)
echo "2. Attempting CML upload..."
curl -s -X POST "$BASE_URL/api/cml/upload?org_id=test_org_001" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}' | jq '.'

echo ""
echo "3. Checking what CML data format is needed..."
