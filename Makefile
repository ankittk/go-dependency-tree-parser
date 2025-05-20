BINARY_NAME = dtree
BINARY_DIR = bin

.PHONY: all
all: build

.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_DIR)/$(BINARY_NAME) main.go
	@echo "Build completed: $(BINARY_DIR)/$(BINARY_NAME)"

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BINARY_DIR)
	@rm -rf output.json
	@echo "Clean done."
