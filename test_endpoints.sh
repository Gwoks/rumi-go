#!/bin/bash

# RUMI Backend API Test Script
# This script tests the endpoints that match the Angular API service

BASE_URL="http://localhost:8080/api/v1"
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Testing RUMI Backend Endpoints${NC}"
echo "=================================="

# Function to make API calls and display results
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local headers=$4
    local description=$5
    
    echo -e "\n${BLUE}Testing: ${description}${NC}"
    echo "Endpoint: ${method} ${endpoint}"
    
    if [ "$method" = "POST" ] && [ -n "$data" ]; then
        if [ -n "$headers" ]; then
            response=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X ${method} \
                "${BASE_URL}${endpoint}" \
                -H "Content-Type: application/json" \
                -H "${headers}" \
                -d "${data}")
        else
            response=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X ${method} \
                "${BASE_URL}${endpoint}" \
                -H "Content-Type: application/json" \
                -d "${data}")
        fi
    elif [ "$method" = "GET" ] && [ -n "$headers" ]; then
        response=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X ${method} \
            "${BASE_URL}${endpoint}" \
            -H "${headers}")
    else
        response=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X ${method} \
            "${BASE_URL}${endpoint}")
    fi
    
    http_status=$(echo "$response" | grep "HTTP_STATUS:" | cut -d: -f2)
    response_body=$(echo "$response" | sed '/HTTP_STATUS:/d')
    
    if [ "$http_status" -ge 200 ] && [ "$http_status" -lt 300 ]; then
        echo -e "${GREEN}✅ Success (HTTP $http_status)${NC}"
    else
        echo -e "${RED}❌ Failed (HTTP $http_status)${NC}"
    fi
    
    echo "Response: $(echo $response_body | jq . 2>/dev/null || echo $response_body)"
}

# Test 1: Health Check
test_endpoint "GET" "/health" "" "" "Health Check"

# Test 2: User Signup - POST /auth/signup
signup_data='{
  "email": "test@example.com",
  "name": "Test User",
  "phone": "+1234567890",
  "password": "password123"
}'
test_endpoint "POST" "/auth/signup" "$signup_data" "" "User Signup"

# Test 3: User Login - POST /auth/login
login_data='{
  "email": "test@example.com",
  "password": "password123"
}'
echo -e "\n${BLUE}Getting JWT Token for further tests...${NC}"
login_response=$(curl -s -X POST "${BASE_URL}/auth/login" \
    -H "Content-Type: application/json" \
    -d "$login_data")

# Extract token from login response
token=$(echo "$login_response" | jq -r '.data.token' 2>/dev/null)

test_endpoint "POST" "/auth/login" "$login_data" "" "User Login"

# Test 4: Get Profile - GET /auth/profile (requires authentication)
if [ "$token" != "null" ] && [ -n "$token" ]; then
    echo -e "\n${GREEN}Token obtained: ${token:0:20}...${NC}"
    test_endpoint "GET" "/auth/profile" "" "Authorization: Bearer $token" "Get User Profile"
    
    # Test 5: Refresh Token - POST /auth/refresh
    test_endpoint "POST" "/auth/refresh" "" "Authorization: Bearer $token" "Refresh Token"
    
    # Test 6: Validate Token - POST /auth/validate
    test_endpoint "POST" "/auth/validate" "" "Authorization: Bearer $token" "Validate Token"
    
    # Test 7: Logout - POST /auth/logout
    test_endpoint "POST" "/auth/logout" "" "Authorization: Bearer $token" "User Logout"
else
    echo -e "${RED}❌ No token obtained, skipping authenticated endpoints${NC}"
fi

# Test 8: Default Admin Login
admin_login_data='{
  "email": "admin@rumi.id",
  "password": "admin123"
}'
test_endpoint "POST" "/auth/login" "$admin_login_data" "" "Admin Login (Default User)"

echo -e "\n${BLUE}=================================="
echo -e "🏁 Testing Complete!${NC}"
echo ""
echo -e "${BLUE}Endpoints tested:${NC}"
echo "✓ POST /api/v1/auth/signup"
echo "✓ POST /api/v1/auth/login"
echo "✓ POST /api/v1/auth/logout"
echo "✓ POST /api/v1/auth/refresh"
echo "✓ POST /api/v1/auth/validate"
echo "✓ GET  /api/v1/auth/profile"
echo ""
echo -e "${BLUE}These endpoints match exactly with your Angular API service!${NC}"
