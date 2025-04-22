.PHONY: all build test clean run fmt lint e2e-test tools check

all: clean fmt lint test build

build:
	go build -o bin/dijester cmd/dijester/main.go

test:
	go test -v ./pkg/...

run:
	go run cmd/dijester/main.go

e2e-test:
	@scripts/e2e_test.sh

fmt:
	golangci-lint fmt

lint:
	golangci-lint run

clean:
	rm -rf bin/
