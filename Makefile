# Build settings

APP_NAME?=krakend-endpoints-tool
IMAGE_REF?=ghcr.io/infratographer/$(APP_NAME)/$(APP_NAME)
IMAGE_TAG?=latest

# Utility settings
TOOLS_DIR := .tools
GOLANGCI_LINT_VERSION = v1.50.1

# Targets

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./tests/...

.PHONY: build
build:
	@echo "Building..."
	@go build -o bin/$(APP_NAME) main.go
	@echo "Build complete. Run with: ./bin/$(APP_NAME)"

.PHONY: image
image:
	@echo "Building docker image... (IMAGE_REF=$(IMAGE_REF) IMAGE_TAG=$(IMAGE_TAG))"
	@docker build -t $(IMAGE_REF):$(IMAGE_TAG) .

lint: golint  ## Runs all lint checks.

golint: $(TOOLS_DIR)/golangci-lint  ## Runs Go lint checks.
	@echo Linting Go files...
	@$(TOOLS_DIR)/golangci-lint run

# Tools setup
$(TOOLS_DIR):
	mkdir -p $(TOOLS_DIR)

$(TOOLS_DIR)/golangci-lint: $(TOOLS_DIR)
	export \
		VERSION=$(GOLANGCI_LINT_VERSION) \
		URL=https://raw.githubusercontent.com/golangci/golangci-lint \
		BINDIR=$(TOOLS_DIR) && \
	curl -sfL $$URL/$$VERSION/install.sh | sh -s $$VERSION
	$(TOOLS_DIR)/golangci-lint version
	$(TOOLS_DIR)/golangci-lint linters

