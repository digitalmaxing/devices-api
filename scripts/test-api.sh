#!/bin/bash

# ============================================
# Devices API - Comprehensive Test Script
# Run with: ./scripts/test-api.sh
# ============================================

set -e

BASE_URL="http://localhost:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

DEVICE_ID=""
IN_USE_DEVICE_ID=""

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}   Devices API - Test Suite${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# Function to print test result
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ PASS${NC} - $2"
    else
        echo -e "${RED}✗ FAIL${NC} - $2"
    fi
}

# ============================================
# 1. BASIC CRUD TESTS
# ============================================
echo -e "${YELLOW}[1/5] Basic CRUD Tests${NC}"

# Create device
echo -n "Creating device... "
RESPONSE=$(curl -s -X POST $BASE_URL/devices \
  -H "Content-Type: application/json" \
  -d '{"name":"Test iPhone", "brand":"Apple", "state":"available"}')

DEVICE_ID=$(echo $RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -n "$DEVICE_ID" ]; then
    print_result 0 "Create device (ID: $DEVICE_ID)"
else
    print_result 1 "Create device"
    exit 1
fi

# Get all devices
echo -n "Fetching all devices... "
STATUS=$(curl -s -o /dev/null -w "%{http_code}" $BASE_URL/devices)
if [ "$STATUS" = "200" ]; then
    print_result 0 "List all devices"
else
    print_result 1 "List all devices (HTTP $STATUS)"
fi

# Get single device
echo -n "Fetching single device... "
STATUS=$(curl -s -o /dev/null -w "%{http_code}" $BASE_URL/devices/$DEVICE_ID)
if [ "$STATUS" = "200" ]; then
    print_result 0 "Get single device"
else
    print_result 1 "Get single device (HTTP $STATUS)"
fi

# Partial update
echo -n "Partial update (change state)... "
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X PATCH $BASE_URL/devices/$DEVICE_ID \
  -H "Content-Type: application/json" \
  -d '{"state":"in-use"}')
if [ "$STATUS" = "200" ]; then
    print_result 0 "Partial update"
else
    print_result 1 "Partial update (HTTP $STATUS)"
fi

# ============================================
# 2. FILTERING TESTS
# ============================================
echo ""
echo -e "${YELLOW}[2/5] Filtering Tests${NC}"


echo -n "Filter by brand (Apple)... "
STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/devices?brand=Apple")
if [ "$STATUS" = "200" ]; then
    print_result 0 "Filter by brand"
else
    print_result 1 "Filter by brand"
fi

echo -n "Filter by state (in-use)... "
STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/devices?state=in-use")
if [ "$STATUS" = "200" ]; then
    print_result 0 "Filter by state"
else
    print_result 1 "Filter by state"
fi

# ============================================
# 3. DOMAIN VALIDATION TESTS (Critical)
# ============================================
echo ""
echo -e "${YELLOW}[3/5] Domain Validation Tests (These should FAIL)${NC}"

# Create in-use device for validation tests
echo -n "Creating in-use device for validation tests... "
RESPONSE=$(curl -s -X POST $BASE_URL/devices \
  -H "Content-Type: application/json" \
  -d '{"name":"Validation Test", "brand":"Samsung", "state":"in-use"}')
IN_USE_DEVICE_ID=$(echo $RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
print_result 0 "Created in-use device"

# Test 1: Try to update creation time
echo -n "Trying to update created_at (should fail)... "
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X PATCH $BASE_URL/devices/$IN_USE_DEVICE_ID \
  -H "Content-Type: application/json" \
  -d '{"created_at":"2020-01-01T00:00:00Z"}')
if [ "$STATUS" = "400" ]; then
    print_result 0 "Block created_at update (got 400)"
else
    print_result 1 "Block created_at update (got $STATUS instead of 400)"
fi

# Test 2: Try to update name while in-use
echo -n "Trying to update name while in-use (should fail)... "
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X PATCH $BASE_URL/devices/$IN_USE_DEVICE_ID \
  -H "Content-Type: application/json" \
  -d '{"name":"Hacked Name"}')
if [ "$STATUS" = "400" ]; then
    print_result 0 "Block name update on in-use device"
else
    print_result 1 "Block name update (got $STATUS)"
fi

# Test 3: Try to delete in-use device
echo -n "Trying to delete in-use device (should fail)... "
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE $BASE_URL/devices/$IN_USE_DEVICE_ID)
if [ "$STATUS" = "409" ]; then
    print_result 0 "Block delete of in-use device (got 409)"
else
    print_result 1 "Block delete (got $STATUS instead of 409)"
fi

# ============================================
# 4. PERSISTENCE TEST
# ============================================
echo ""
echo -e "${YELLOW}[4/5] Persistence Test${NC}"

echo -n "Creating device for persistence test... "
RESPONSE=$(curl -s -X POST $BASE_URL/devices \
  -H "Content-Type: application/json" \
  -d '{"name":"Persistence Check", "brand":"Test", "state":"available"}')
PERSIST_ID=$(echo $RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)

echo "Restarting containers (this may take 10-15 seconds)..."
docker-compose restart > /dev/null 2>&1
sleep 12

echo -n "Checking if device survived restart... "
STATUS=$(curl -s -o /dev/null -w "%{http_code}" $BASE_URL/devices/$PERSIST_ID)
if [ "$STATUS" = "200" ]; then
    print_result 0 "Data persisted after restart"
else
    print_result 1 "Data did NOT persist (HTTP $STATUS)"
fi

# ============================================
# 5. EDGE CASES
# ============================================
echo ""
echo -e "${YELLOW}[5/5] Edge Cases${NC}"

echo -n "Get non-existent device... "
STATUS=$(curl -s -o /dev/null -w "%{http_code}" $BASE_URL/devices/00000000-0000-0000-0000-000000000000)
if [ "$STATUS" = "404" ]; then
    print_result 0 "Returns 404 for non-existent device"
else
    print_result 1 "Expected 404, got $STATUS"
fi

echo -n "Delete non-existent device... "
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE $BASE_URL/devices/00000000-0000-0000-0000-000000000000)
if [ "$STATUS" = "404" ]; then
    print_result 0 "Returns 404 for deleting non-existent device"
else
    print_result 1 "Expected 404, got $STATUS"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   Test Suite Completed${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Tip: Run this script again after any changes:"
echo "  ./scripts/test-api.sh"
echo ""