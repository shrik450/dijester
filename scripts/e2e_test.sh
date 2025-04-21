#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "=== Dijester End-to-End Test ==="
echo

# Build the binary
echo "Building dijester..."
go build -o bin/dijester cmd/dijester/main.go

# Test RSS source
echo
echo "Testing RSS source..."
./bin/dijester --test-source rss &> /tmp/rss_test_output.txt
if grep -q "Fetched" /tmp/rss_test_output.txt; then
    echo -e "${GREEN}✓ RSS source test passed${NC}"
else
    echo -e "${RED}✗ RSS source test failed${NC}"
    cat /tmp/rss_test_output.txt
    exit 1
fi

# Test Hacker News source
echo
echo "Testing Hacker News source..."
./bin/dijester --test-source hackernews &> /tmp/hn_test_output.txt
if grep -q "Fetched" /tmp/hn_test_output.txt; then
    echo -e "${GREEN}✓ Hacker News source test passed${NC}"
else
    echo -e "${RED}✗ Hacker News source test failed${NC}"
    cat /tmp/hn_test_output.txt
    exit 1
fi

# Create a test directory for output
TEST_DIR="/tmp/dijester_test_output"
rm -rf $TEST_DIR
mkdir -p $TEST_DIR

# Run with example config and generate output
echo
echo "Testing digest generation with example configuration..."
./bin/dijester --config example-config.toml --output $TEST_DIR/digest.md &> /tmp/config_test_output.txt
if ! grep -q "Initialized 2 active sources" /tmp/config_test_output.txt; then
    echo -e "${RED}✗ Configuration test failed${NC}"
    cat /tmp/config_test_output.txt
    exit 1
fi

# Check if the output file exists
if [ ! -f $TEST_DIR/digest.md ]; then
    echo -e "${RED}✗ Output file not created!${NC}"
    exit 1
fi

# Check the content of the output file
if ! grep -q "My Daily News Digest" $TEST_DIR/digest.md; then
    echo -e "${RED}✗ Output file does not contain expected title!${NC}"
    cat $TEST_DIR/digest.md | head -n 10
    exit 1
fi

echo -e "${GREEN}✓ Digest generated successfully${NC}"
echo "First 10 lines of output:"
head -n 10 $TEST_DIR/digest.md

echo
echo -e "${GREEN}All end-to-end tests passed!${NC}"