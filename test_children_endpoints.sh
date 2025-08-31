#!/bin/bash

# RUMI Backend Children API Test Script
# This script tests the children endpoints

BASE_URL="http://localhost:8080/api/v1"
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Testing RUMI Children Endpoints${NC}"
echo "===================================="

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
    elif [ "$method" = "PUT" ] && [ -n "$data" ]; then
        response=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X ${method} \
            "${BASE_URL}${endpoint}" \
            -H "Content-Type: application/json" \
            -H "${headers}" \
            -d "${data}")
    elif [ -n "$headers" ]; then
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

# Test 1: Login to get JWT token
echo -e "\n${BLUE}Getting JWT Token for tests...${NC}"
login_data='{
  "email": "user@rumi.id",
  "password": "admin123"
}'

login_response=$(curl -s -X POST "${BASE_URL}/auth/login" \
    -H "Content-Type: application/json" \
    -d "$login_data")

# Extract token from login response
token=$(echo "$login_response" | jq -r '.data.token' 2>/dev/null)

if [ "$token" != "null" ] && [ -n "$token" ]; then
    echo -e "${GREEN}Token obtained: ${token:0:20}...${NC}"
    
    # Test 2: Create a child
    child_data='{
      "name": "Alice Johnson",
      "nick_name": "Ali",
      "birth_date": "2015-06-15"
    }'
    test_endpoint "POST" "/children" "$child_data" "Authorization: Bearer $token" "Create Child"
    
    # Test 3: Get user's children
    test_endpoint "GET" "/children" "" "Authorization: Bearer $token" "Get User Children"
    
    # Test 4: Get child by ID (assuming ID 1 exists)
    test_endpoint "GET" "/children/1" "" "Authorization: Bearer $token" "Get Child by ID"
    
    # Test 5: Update child (assuming ID 1 exists)
    update_data='{
      "name": "Alice Marie Johnson",
      "nick_name": "Allie",
      "birth_date": "2015-06-15"
    }'
    test_endpoint "PUT" "/children/1" "$update_data" "Authorization: Bearer $token" "Update Child"
    
    # Test 6: Admin login for admin tests
    echo -e "\n${BLUE}Getting Admin JWT Token...${NC}"
    admin_login_data='{
      "email": "admin@rumi.id",
      "password": "admin123"
    }'
    
    admin_login_response=$(curl -s -X POST "${BASE_URL}/auth/login" \
        -H "Content-Type: application/json" \
        -d "$admin_login_data")
    
    admin_token=$(echo "$admin_login_response" | jq -r '.data.token' 2>/dev/null)
    
    if [ "$admin_token" != "null" ] && [ -n "$admin_token" ]; then
        echo -e "${GREEN}Admin token obtained: ${admin_token:0:20}...${NC}"
        
        # Test 7: Get all children (admin only)
        test_endpoint "GET" "/admin/children" "" "Authorization: Bearer $admin_token" "Get All Children (Admin)"
        
        # Test 8: Set child active status (admin only)
        status_data='{"is_active": false}'
        test_endpoint "PUT" "/admin/children/1/active" "$status_data" "Authorization: Bearer $admin_token" "Deactivate Child (Admin)"
    else
        echo -e "${RED}❌ No admin token obtained, skipping admin endpoints${NC}"
    fi
    
else
    echo -e "${RED}❌ No token obtained, skipping authenticated endpoints${NC}"
fi

echo -e "\n${BLUE}===================================="
echo -e "🏁 Children API Testing Complete!${NC}"
echo ""
echo -e "${BLUE}Endpoints tested:${NC}"
echo "✓ POST /api/v1/children"
echo "✓ GET  /api/v1/children"
echo "✓ GET  /api/v1/children/:id"
echo "✓ PUT  /api/v1/children/:id"
echo "✓ GET  /api/v1/admin/children"
echo "✓ PUT  /api/v1/admin/children/:id/active"
