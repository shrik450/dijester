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
TEST_DIR=$(mktemp -d)
if [ ! -d $TEST_DIR ]; then
    echo -e "${RED}✗ Failed to create temporary directory!${NC}"
    exit 1
fi
echo "Temporary directory created at $TEST_DIR"

# Run with example config and generate output
echo
echo "Testing digest generation with example configuration..."
./bin/dijester --config example-config-md.toml --output-dir $TEST_DIR &> "$TEST_DIR/config_test_output.txt"
if ! grep -q "Fetching from source:  hackernews" "$TEST_DIR/config_test_output.txt"; then
    echo -e "${RED}✗ Configuration test failed${NC}"
    cat "$TEST_DIR/config_test_output.txt"
    exit 1
fi

# Check if the output file exists
if [ ! -f "$TEST_DIR/digest-$(date +%Y)-$(date +%-m)-$(date +%-d).md" ]; then
    echo -e "${RED}✗ Output file not created!${NC}"
    exit 1
fi

# Check the content of the output file
if ! grep -q "My Daily News Digest" "$TEST_DIR/digest-$(date +%Y)-$(date +%-m)-$(date +%-d).md"; then
    echo -e "${RED}✗ Output file does not contain expected title!${NC}"
    cat "$TEST_DIR/digest-$(date +%Y)-$(date +%-m)-$(date +%-d).md" | head -n 10
    exit 1
fi

echo -e "${GREEN}✓ Markdown digest generated successfully${NC}"
echo "First 10 lines of output:"
head -n 10 "$TEST_DIR/digest-$(date +%Y)-$(date +%-m)-$(date +%-d).md"

# Test EPUB output
echo
echo "Testing EPUB digest generation..."
./bin/dijester --config example-config-epub.toml --output-dir $TEST_DIR &> "$TEST_DIR/epub_test_output.txt"
if ! grep -q "Fetching from source:  hackernews" "$TEST_DIR/epub_test_output.txt"; then
    echo -e "${RED}✗ EPUB configuration test failed${NC}"
    cat "$TEST_DIR/epub_test_output.txt"
    exit 1
fi

if grep -q "Error fetching from rss:" "$TEST_DIR/epub_test_output.txt"; then
    echo -e "${RED}✗ EPUB test failed with errors!${NC}"
    cat "$TEST_DIR/epub_test_output.txt"
    exit 1
fi

# Check if the EPUB output file exists
# Look for any .epub file in the output directory since the filename uses a date template
EPUB_FILE=$(find $TEST_DIR -name "*.epub" | head -n 1)
if [ -z "$EPUB_FILE" ]; then
    echo -e "${RED}✗ EPUB output file not created!${NC}"
    exit 1
fi

# Check if the EPUB file is valid
if ! file "$EPUB_FILE" | grep -q "EPUB document"; then
    echo -e "${RED}✗ EPUB file is not a valid EPUB document!${NC}"
    file "$EPUB_FILE"
    exit 1
fi

echo -e "${GREEN}✓ EPUB digest generated successfully${NC}"
echo "EPUB file size: $(du -h "$EPUB_FILE" | cut -f1)"

echo
echo -e "${GREEN}All end-to-end tests passed!${NC}"
