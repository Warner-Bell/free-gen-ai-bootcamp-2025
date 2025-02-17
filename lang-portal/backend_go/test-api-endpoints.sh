#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color
BLUE='\033[0;34m'

# Base URL
BASE_URL="http://localhost:8080/api"

# Function to print test header
print_test() {
    echo -e "\n${BLUE}=== Testing: $1 ===${NC}"
}

# Function to check if command succeeded
check_result() {
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ PASS${NC}: $1"
    else
        echo -e "${RED}✗ FAIL${NC}: $1"
        exit 1
    fi
}

# Test GET all groups
print_test "GET all groups"
response=$(curl -s -w "\n%{http_code}" -H "Accept: application/json" "${BASE_URL}/groups")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 200 ]; then
    check_result "GET /groups returned 200"
    if echo "$body" | jq . >/dev/null 2>&1; then
        check_result "Response is valid JSON"
    else
        echo -e "${RED}✗ FAIL${NC}: Invalid JSON response"
        exit 1
    fi
else
    echo -e "${RED}✗ FAIL${NC}: GET /groups returned $http_code"
    exit 1
fi

# Test GET single group
print_test "GET single group (ID: 1)"
response=$(curl -s -w "\n%{http_code}" -H "Accept: application/json" "${BASE_URL}/groups/1")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 200 ]; then
    check_result "GET /groups/1 returned 200"
    if echo "$body" | jq . >/dev/null 2>&1; then
        check_result "Response is valid JSON"
        
        # Check required fields
        group_id=$(echo "$body" | jq -r '.id')
        if [ "$group_id" -eq 1 ]; then
            check_result "Group ID matches expected value"
        else
            echo -e "${RED}✗ FAIL${NC}: Unexpected group ID: $group_id"
            exit 1
        fi
    else
        echo -e "${RED}✗ FAIL${NC}: Invalid JSON response"
        exit 1
    fi
else
    echo -e "${RED}✗ FAIL${NC}: GET /groups/1 returned $http_code"
    exit 1
fi

# Test CREATE new group
print_test "CREATE new group"
response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "Accept: application/json" \
    -d '{"name":"Test Group","description":"Test Description"}' \
    "${BASE_URL}/groups")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 201 ] || [ "$http_code" -eq 200 ]; then
    check_result "POST /groups returned ${http_code}"
    if echo "$body" | jq . >/dev/null 2>&1; then
        check_result "Response is valid JSON"
        new_group_id=$(echo "$body" | jq -r '.id')
        echo -e "${GREEN}Created group ID: $new_group_id${NC}"
    else
        echo -e "${RED}✗ FAIL${NC}: Invalid JSON response"
        exit 1
    fi
else
    echo -e "${RED}✗ FAIL${NC}: POST /groups returned $http_code"
    exit 1
fi

# Test UPDATE group
if [ ! -z "$new_group_id" ]; then
    print_test "UPDATE group"
    response=$(curl -s -w "\n%{http_code}" -X PUT \
        -H "Content-Type: application/json" \
        -H "Accept: application/json" \
        -d '{"name":"Updated Test Group","description":"Updated Description"}' \
        "${BASE_URL}/groups/${new_group_id}")
    http_code=$(echo "$response" | tail -n1)

    if [ "$http_code" -eq 200 ]; then
        check_result "PUT /groups/${new_group_id} returned 200"
    else
        echo -e "${RED}✗ FAIL${NC}: PUT /groups/${new_group_id} returned $http_code"
        exit 1
    fi
fi

# Test DELETE group
if [ ! -z "$new_group_id" ]; then
    print_test "DELETE group"
    response=$(curl -s -w "\n%{http_code}" -X DELETE \
        -H "Accept: application/json" \
        "${BASE_URL}/groups/${new_group_id}")
    http_code=$(echo "$response" | tail -n1)

    if [ "$http_code" -eq 200 ] || [ "$http_code" -eq 204 ]; then
        check_result "DELETE /groups/${new_group_id} returned ${http_code}"
    else
        echo -e "${RED}✗ FAIL${NC}: DELETE /groups/${new_group_id} returned $http_code"
        exit 1
    fi
fi


# Test non-existent group
print_test "GET non-existent group"
response=$(curl -s -w "\n%{http_code}" -H "Accept: application/json" "${BASE_URL}/groups/999999")
http_code=$(echo "$response" | tail -n1)

if [ "$http_code" -eq 404 ]; then
    check_result "GET /groups/999999 returned 404 as expected"
else
    echo -e "${RED}✗ FAIL${NC}: GET /groups/999999 returned $http_code (expected 404)"
    exit 1
fi

echo -e "\n${GREEN}All tests completed successfully!${NC}"
