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

# Create a test directory for output
TEST_DIR="/tmp/dijester_test_output"
rm -rf $TEST_DIR
mkdir -p $TEST_DIR

# Run with example config and generate output
echo
echo "Testing digest generation with example configuration..."
./bin/dijester --config example-config-md.toml --output $TEST_DIR/digest.md &> /tmp/config_test_output.txt
if ! grep -q "Fetching from source:  hackernews" /tmp/config_test_output.txt; then
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

echo -e "${GREEN}✓ Markdown digest generated successfully${NC}"
echo "First 10 lines of output:"
head -n 10 $TEST_DIR/digest.md

# Test EPUB output
echo
echo "Testing EPUB digest generation..."
./bin/dijester --config example-config-epub.toml --output $TEST_DIR/digest.epub &> /tmp/epub_test_output.txt
if ! grep -q "Fetching from source:  hackernews" /tmp/epub_test_output.txt; then
    echo -e "${RED}✗ EPUB configuration test failed${NC}"
    cat /tmp/epub_test_output.txt
    exit 1
fi

# Check if the EPUB output file exists
if [ ! -f $TEST_DIR/digest.epub ]; then
    echo -e "${RED}✗ EPUB output file not created!${NC}"
    exit 1
fi

# Check if the EPUB file is valid
if ! file $TEST_DIR/digest.epub | grep -q "EPUB document"; then
    echo -e "${RED}✗ EPUB file is not a valid EPUB document!${NC}"
    file $TEST_DIR/digest.epub
    exit 1
fi

echo -e "${GREEN}✓ EPUB digest generated successfully${NC}"
echo "EPUB file size: $(du -h $TEST_DIR/digest.epub | cut -f1)"

echo
echo -e "${GREEN}All end-to-end tests passed!${NC}"
