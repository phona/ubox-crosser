SOURCES := $(shell find . -name "*.go")
MODULE := github.com/phona/ubox-crosser
BINARIES := client server auth_server
COVERAGE_DIR := coverage
BUILD_ID ?= $(shell date +%s)

# ===========================================
# Build Commands
# ===========================================

.PHONY: build clean fmt vet test unit-test unit-test-coverage

GIT_SHA := $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")

build: $(SOURCES)
	@echo "=== Building binaries ==="
	CGO_ENABLED=0 go build -ldflags="-s -w -X main.GitSHA=$(GIT_SHA)" -o bin/client ./cmd/client
	CGO_ENABLED=0 go build -ldflags="-s -w -X main.GitSHA=$(GIT_SHA)" -o bin/server ./cmd/server
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/auth_server ./cmd/auth_server

clean:
	rm -rf bin/ $(COVERAGE_DIR)/

fmt:
	@echo "=== Formatting Code ==="
	go fmt ./...

vet:
	@echo "=== Running go vet ==="
	go vet ./...

test:
	@echo "=== Running Tests ==="
	go test -v -count=1 ./...

# Run unit tests only (no integration)
unit-test:
	@echo "=== Running Unit Tests ==="
	go test -short -v -count=1 ./...

# Run unit tests with coverage report
unit-test-coverage:
	@echo "=== Running Unit Tests with Coverage ==="
	@mkdir -p $(COVERAGE_DIR)
	go test -short -v -count=1 -coverprofile=$(COVERAGE_DIR)/unit.out -covermode=set ./...

# ===========================================
# Code Quality Commands
# ===========================================

.PHONY: lint

lint:
	@echo "=== Running Linter ==="
	@GOPATH_BIN=$$(go env GOPATH)/bin; \
	export PATH="$$GOPATH_BIN:$$PATH"; \
	if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "golangci-lint not found, installing to $$GOPATH_BIN ..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$GOPATH_BIN v1.62.2; \
	fi; \
	golangci-lint run ./...

# ===========================================
# Integration Testing Commands
# ===========================================

.PHONY: test-integration test-clean test-help sonar

test-help:
	@echo "UBox-Crosser Test Commands"
	@echo ""
	@echo "  make unit-test             Run unit tests"
	@echo "  make unit-test-coverage    Run unit tests with coverage"
	@echo "  make test-integration      Run integration tests (Docker required)"
	@echo "  make test-clean            Clean up test containers and volumes"
	@echo ""
	@echo "Options:"
	@echo "  BUILD_ID=<id>              Unique identifier for test run (default: timestamp)"

# Run SonarQube analysis locally
sonar:
	@echo "=== Running SonarQube Analysis ==="
	sonar-scanner -Dproject.settings=sonar-project.properties

test-integration:
	@echo "=== Running Integration Tests ==="
	@echo "=== Cleaning up any previous run ==="
	docker compose -p crosser-test-$(BUILD_ID) -f tests/docker-compose.yml down -v --remove-orphans 2>/dev/null || true
	@echo "=== Starting containers ==="
	@mkdir -p $(COVERAGE_DIR)/raw
	@chmod -R 777 $(COVERAGE_DIR)
	COVERAGE_HOST_DIR=$(CURDIR)/$(COVERAGE_DIR)/raw \
		docker compose -p crosser-test-$(BUILD_ID) -f tests/docker-compose.yml up --build --exit-code-from test-runner
	@echo "=== Merging coverage data ==="
	@if ls $(COVERAGE_DIR)/raw/cov* 1>/dev/null 2>&1; then \
		go tool covdata textfmt -i=$(COVERAGE_DIR)/raw -o=$(COVERAGE_DIR)/integration.out; \
		echo "Coverage report: $(COVERAGE_DIR)/integration.out"; \
	else \
		echo "No coverage data collected"; \
	fi
	@echo "=== Cleaning up test containers and volumes ==="
	docker compose -p crosser-test-$(BUILD_ID) -f tests/docker-compose.yml down -v --remove-orphans 2>/dev/null || true

test-clean:
	@echo "=== Cleaning up all test containers ==="
	docker compose -p crosser-test-$(BUILD_ID) -f tests/docker-compose.yml down -v --remove-orphans 2>/dev/null || true
	@rm -rf $(COVERAGE_DIR)/*
	@echo "Done."

# ═══════════════════════════════════════════════════
# CI Standard Interface
# ═══════════════════════════════════════════════════

.PHONY: ci-env ci-setup ci-lint ci-unit-test ci-integration-test ci-build ci-test dev-cross-check

ci-env:
	@echo "GO_VERSION=1.23"
	@echo "NEEDS_DOCKER=true"

ci-setup:
	go mod download
	@which golangci-lint > /dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.62.2

# Code lint (parallel go vet + golangci-lint, BASE_REV for incremental scan)
ci-lint:
	@golangci-lint run $${BASE_REV:+--new-from-rev=$$BASE_REV}

ci-unit-test:
	@mkdir -p $(COVERAGE_DIR)
	go test -short -v -count=1 -coverprofile=$(COVERAGE_DIR)/unit.out -covermode=set ./...

ci-integration-test:
	$(MAKE) test-integration BUILD_ID="$(BUILD_ID)"

ci-build:
	@echo "Building Docker images..."
	docker build -t ubox-crosser-client --build-arg BINARY=client -f Dockerfile .
	docker build -t ubox-crosser-server --build-arg BINARY=server -f Dockerfile .
	docker build -t ubox-crosser-auth-server --build-arg BINARY=auth_server -f Dockerfile .

# ci-test: unit tests (staging_test gate)
ci-test:
	$(MAKE) ci-unit-test

# dev-cross-check: lint + unit tests (dev_cross_check gate)
dev-cross-check:
	$(MAKE) ci-lint
	$(MAKE) ci-unit-test
