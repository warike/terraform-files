# Makefile for wtf (Warike Terraform Files)

.PHONY: all build test test-unit test-integration clean help

all: build

## build: Compile the binary
build:
	go build -o build/tfinit ./cmd/tfinit

## test: Run all tests
test: test-unit test-integration

## test-unit: Run only unit tests (fast, no network)
test-unit:
	go test -v -short ./...

## test-integration: Run integration tests (hits real API)
test-integration:
	go test -v -run TestGetLatestVersion_RealRegistry ./internal/providers/...

## test-updater: Run updater tests
test-updater:
	go test -v ./internal/updater/...

## clean: Remove build artifacts
clean:
	rm -rf build/
	rm -f provider.tf variables.tf terraform.tfvars main.tf

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'