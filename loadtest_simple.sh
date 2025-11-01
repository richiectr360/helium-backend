#!/bin/bash

# Simple Load Testing Script using Apache Bench and curl
BASE_URL="${BASE_URL:-http://localhost:8000}"

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}=== Backend Load Testing ===${NC}\n"

# Check if server is running
if ! curl -s "${BASE_URL}/health" > /dev/null; then
    echo -e "${RED}❌ Server is not running at ${BASE_URL}${NC}"
    echo "Please start the server first: cd go-localization-manager-backend && make run"
    exit 1
fi
echo -e "${GREEN}✓ Server is running${NC}\n"

# Test 1: Single request - cold cache
echo -e "${BLUE}Test 1: Cold Cache (Single Request)${NC}"
echo "Clearing cache by making first request..."
curl -s "${BASE_URL}/api/component/welcome?lang=en" > /dev/null
sleep 1
echo "Measuring cold cache request..."
time curl -s "${BASE_URL}/api/component/welcome?lang=en" > /dev/null
echo ""

# Test 2: Warm cache - 100 requests
echo -e "${BLUE}Test 2: Warm Cache (100 requests, 10 concurrent)${NC}"
ab -n 100 -c 10 "${BASE_URL}/api/component/welcome?lang=en" 2>&1 | grep -E "(Requests per second|Time per request|Failed requests|Transfer rate)"
echo ""

# Test 3: Different languages
echo -e "${BLUE}Test 3: Different Languages (50 requests each)${NC}"
for lang in es fr de; do
    echo "Testing lang=$lang..."
    ab -n 50 -c 5 "${BASE_URL}/api/component/welcome?lang=${lang}" 2>&1 | grep -E "(Requests per second|Time per request)"
done
echo ""

# Test 4: Different components
echo -e "${BLUE}Test 4: Different Components (50 requests each)${NC}"
for comp in navigation user_profile footer; do
    echo "Testing component=$comp..."
    ab -n 50 -c 5 "${BASE_URL}/api/component/${comp}?lang=en" 2>&1 | grep -E "(Requests per second|Time per request)"
done
echo ""

# Test 5: Concurrency test (beyond default limit of 20)
echo -e "${BLUE}Test 5: Concurrency Test (5 concurrent, default limit=20, expect few/no 503s)${NC}"
ab -n 50 -c 5 "${BASE_URL}/api/component/welcome?lang=en" 2>&1 | grep -E "(Requests per second|Time per request|Failed requests)"
echo ""

# Test 6: Sustained load
echo -e "${BLUE}Test 6: Sustained Load (500 requests, 20 concurrent)${NC}"
ab -n 500 -c 20 "${BASE_URL}/api/component/welcome?lang=en" 2>&1 | grep -E "(Requests per second|Time per request|Failed requests|Transfer rate)"
echo ""

# Test 7: Invalid language (should return 400)
echo -e "${BLUE}Test 7: Invalid Language Validation${NC}"
response=$(curl -s -w "\n%{http_code}" "${BASE_URL}/api/component/welcome?lang=invalid")
http_code=$(echo "$response" | tail -n1)
if [ "$http_code" = "400" ]; then
    echo -e "${GREEN}✓ Correctly returned 400 for invalid language${NC}"
else
    echo -e "${RED}✗ Expected 400, got $http_code${NC}"
fi
echo ""

# Test 8: Invalid component (should return 404)
echo -e "${BLUE}Test 8: Invalid Component Validation${NC}"
response=$(curl -s -w "\n%{http_code}" "${BASE_URL}/api/component/invalid?lang=en")
http_code=$(echo "$response" | tail -n1)
if [ "$http_code" = "404" ]; then
    echo -e "${GREEN}✓ Correctly returned 404 for invalid component${NC}"
else
    echo -e "${RED}✗ Expected 404, got $http_code${NC}"
fi
echo ""

echo -e "${BLUE}=== Load Testing Complete ===${NC}"

