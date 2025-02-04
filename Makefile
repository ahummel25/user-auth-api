# Build configuration
BUILD_PREFIX          = env GOOS=linux GOARCH=arm64 go build
COMMON_LDFLAGS        = -s -w -extldflags '-static'
COMMON_TAGS           = lambda.norpc netgo
DEPLOYMENT_STAGE      = ${TF_WORKSPACE}
BUILD_DIR            = build
PACKAGES_DIR         = $(BUILD_DIR)/packages

# Environment configuration
ENV_TAGS              = $(if $(filter prod,$(DEPLOYMENT_STAGE)),prod,$(if $(filter dev,$(DEPLOYMENT_STAGE)),dev,))

# Lambda configuration
LAMBDA_DIR           = lambda
LAMBDA_CMD_DIR       = cmd/lambda

.PHONY: all build clean deploy gomodgen local dev-deps logs logs-mongo help check-deps

# Default target
all: check-deps build

# Help target
help:
	@echo "Available targets:"
	@echo "  all          - Build everything (default target)"
	@echo "  build        - Build Lambda functions"
	@echo "  clean        - Remove build artifacts"
	@echo "  deploy       - Deploy to AWS using serverless"
	@echo "  gomodgen     - Generate go.mod file"
	@echo "  local        - Run local development server"
	@echo "  dev-deps     - Start development dependencies"
	@echo "  logs         - View all logs"
	@echo "  logs-mongo   - View MongoDB logs"
	@echo "  help         - Show this help message"
	@echo "  check-deps   - Check required dependencies"

# Check for required dependencies
check-deps:
	@command -v go >/dev/null 2>&1 || { echo "Error: Go is not installed"; exit 1; }
	@command -v docker >/dev/null 2>&1 || { echo "Error: Docker is not installed"; exit 1; }
	@command -v docker compose >/dev/null 2>&1 || { echo "Error: docker-compose is not installed"; exit 1; }
	@command -v sls >/dev/null 2>&1 || { echo "Error: Serverless Framework is not installed"; exit 1; }

build: 
	@echo "Building Lambda functions..."
	@mkdir -p $(PACKAGES_DIR)
	@for func in $(LAMBDA_DIR)/*/$(LAMBDA_CMD_DIR); do \
		funcname=$$(echo $$func | cut -d'/' -f2); \
		mkdir -p $(PACKAGES_DIR)/$$funcname; \
		$(BUILD_PREFIX) -ldflags "$(COMMON_LDFLAGS)" -tags "$(COMMON_TAGS) $(ENV_TAGS)" \
			-trimpath \
			-o $(PACKAGES_DIR)/$$funcname/bootstrap $$func/*.go; \
		zip -j $(PACKAGES_DIR)/$$funcname.zip $(PACKAGES_DIR)/$$funcname/bootstrap; \
	done
	@echo "Build completed successfully!"

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf ./$(BUILD_DIR) ./vendor Gopkg.lock
	@echo "Clean completed!"

deploy: clean build
	@echo "Deploying to $(DEPLOYMENT_STAGE)..."
	@sls deploy --stage $(DEPLOYMENT_STAGE) --verbose

gomodgen:
	@chmod u+x gomod.sh
	@./gomod.sh

# Local development targets
local:
	@echo "Starting local development server..."
	@go run cmd/local/main.go

dev-deps:
	@echo "Starting development dependencies..."
	@docker compose up -d

# Logging targets
logs-mongo:
	@docker compose logs -f mongodb

logs: logs-mongo