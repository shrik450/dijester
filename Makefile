.PHONY: all build test clean run fmt lint e2e-test tools check

all: fmt lint test build

build:
	go build -o bin/dijester cmd/dijester/main.go

test:
	go test -v ./pkg/...

run:
	go run cmd/dijester/main.go

test-rss:
	go run cmd/dijester/main.go --test-source rss

test-hn:
	go run cmd/dijester/main.go --test-source hackernews

e2e-test:
	@chmod +x scripts/e2e_test.sh
	@scripts/e2e_test.sh

tools:
	@echo "Installing tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	go install mvdan.cc/gofumpt

fmt: tools
	@echo "Formatting code..."
	$(shell go env GOPATH)/bin/gofumpt -l -w .

lint: tools
	@echo "Linting code..."
	$(shell go env GOPATH)/bin/golangci-lint run

check: fmt lint test

clean:
	rm -rf bin/