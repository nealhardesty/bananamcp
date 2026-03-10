# BananaMCP Makefile

BINARY     := banana
CMD        := ./cmd/banana
VERSION    := $(shell grep 'const Version' cmd/banana/version.go | sed 's/.*"\(.*\)".*/\1/')

.PHONY: all build test lint fmt tidy clean run help \
        version version-increment release push

## all: Build the banana binary (default)
all: build

## build: Build the banana binary
build:
	go build -o $(BINARY) $(CMD)

## test: Run all tests with race detector
test:
	go test -race ./...

## lint: Run go vet and golangci-lint (if available)
lint:
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping"; \
	fi

## fmt: Format code with gofmt
fmt:
	gofmt -w .

## tidy: Run go mod tidy
tidy:
	go mod tidy

## clean: Remove build artifacts
clean:
	rm -f $(BINARY) output.*

## run-mcp: Build and run the MCP server
run-mcp: build
	OPENROUTER_API_KEY=$${OPENROUTER_API_KEY} ./$(BINARY) mcp

## version: Display current version
version:
	@echo $(VERSION)

## version-increment: Increment patch version in cmd/banana/version.go
version-increment:
	@current=$(VERSION); \
	major=$$(echo $$current | cut -d. -f1); \
	minor=$$(echo $$current | cut -d. -f2); \
	patch=$$(echo $$current | cut -d. -f3); \
	newpatch=$$((patch + 1)); \
	newver="$$major.$$minor.$$newpatch"; \
	sed -i "s/const Version = \"$$current\"/const Version = \"$$newver\"/" cmd/banana/version.go; \
	echo "Version bumped: $$current → $$newver"

## release: Bump version, build, tag, and push
release: version-increment build
	@newver=$$(grep 'const Version' cmd/banana/version.go | sed 's/.*"\(.*\)".*/\1/'); \
	git add -A && \
	git commit -m "Release v$$newver. $$(gitsum)" && \
	git tag "v$$newver" && \
	git push && \
	git push --tags && \
	echo "Released v$$newver"

## push: Commit all changes and push
push:
	git add -A
	git commit -m "$$(gitsum)"
	git push

## help: Show this help message
help:
	@grep -E '^## ' Makefile | sed 's/## //' | column -t -s ':'
