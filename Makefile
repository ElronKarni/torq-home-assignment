.PHONY: build run test test-coverage clean dev rate-limit-test test-package

# Go parameters
BINARY_NAME=ip2country-api
MAIN_PATH=./cmd/main.go
COVER_PROFILE=coverage.out

# Build the application
build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# Run the application
run: build
	./$(BINARY_NAME)

# Run tests
test:
	go test ./... -v

# Run tests with coverage
test-coverage:
	go test ./... -coverprofile=$(COVER_PROFILE)
	go tool cover -html=$(COVER_PROFILE)

# Test a specific package with coverage
test-package:
	@echo "Usage: make test-package PKG=./pkg/ratelimit"
	@if [ "$(PKG)" != "" ]; then \
		go test $(PKG) -coverprofile=$(COVER_PROFILE) -v; \
		go tool cover -html=$(COVER_PROFILE); \
	fi

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(COVER_PROFILE)


# Real run rate limit test script
rate-limit-test:
	bash test-rate-limit.sh 