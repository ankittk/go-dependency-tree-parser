BINARY_NAME = dtree
BINARY_DIR = bin

.PHONY: all clean test build

all: build

clean:
	@echo "Cleaning up..."
	@rm -rf $(BINARY_DIR)
	@rm -f output.json
	@echo "Clean done."

test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "Tests completed."

lint:
	@echo "Running golangci-lint..."
	@golangci-lint run ./...
	@echo "Linting completed."

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_DIR)/$(BINARY_NAME) .
	@echo "Build completed: $(BINARY_DIR)/$(BINARY_NAME)"

.PHONY: run
run: build
	@./$(BINARY_DIR)/$(BINARY_NAME) parse $(repo) $(tag)
