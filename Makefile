BINARY_NAME=chartcli
BUILD_DIR=build
PROJECT_ROOT=github.com/lucasmlp/release-yaml-utils

.PHONY: all build clean generate tobereleased

all: build

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/chartcli

generate:
	@$(BUILD_DIR)/$(BINARY_NAME) generate

tobereleased:
	@$(BUILD_DIR)/$(BINARY_NAME) tobereleased

clean:
	@rm -rf $(BUILD_DIR)
