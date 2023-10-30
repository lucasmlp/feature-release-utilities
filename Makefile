# Makefile

BINARY_NAME=chartcli
SOURCE_DIR=./
BUILD_DIR=./build

all: build

build:
	@echo "Building..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SOURCE_DIR)

run: build
	@echo "Running..."
	$(BUILD_DIR)/$(BINARY_NAME) generate

to-be-released: build
	@echo "Running..."
	$(BUILD_DIR)/$(BINARY_NAME) tobereleased

clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)

.PHONY: all build run clean
