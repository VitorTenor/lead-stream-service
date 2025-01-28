#!/bin/bash

# Define colors for pretty output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to run tests and check for errors
run_tests() {
    local path=$1
    echo -e "${YELLOW}Running tests in $path...${NC}"
    if go test -v "$path"; then
        echo -e "${GREEN}Tests passed in $path${NC}"
    else
        echo -e "${RED}Tests failed in $path${NC}"
        exit 1
    fi
}

ENVIRONMENT=${ENVIRONMENT:-LOCAL}

# Set prefix based on environment
if [ "$ENVIRONMENT" == "LOCAL" ]; then
    PREFIX="../"
elif [ "$ENVIRONMENT" == "TEST" ]; then
    PREFIX="./"
else
    PREFIX=""
fi

# Service path
SERVICE_PATH="${PREFIX}internal/services"

# Integration test path
INTEGRATION_TEST_PATH="${PREFIX}internal/integration"

# Run all tests
echo -e "${YELLOW}Running all Go tests in environment: $ENVIRONMENT...${NC}"
run_tests "$SERVICE_PATH"
run_tests "$INTEGRATION_TEST_PATH"

echo -e "${GREEN}All tests passed successfully.${NC}"