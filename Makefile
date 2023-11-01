BINARY_NAME=chartcli
BUILD_DIR=build
PROJECT_ROOT=github.com/lucasmlp/release-yaml-utils

.PHONY: all build clean generate tobereleased count merge

all: build clean

clean:
	@rm -rf $(BUILD_DIR)

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/chartcli

generate: build
	@$(BUILD_DIR)/$(BINARY_NAME) generate

tobereleased: build
	@$(BUILD_DIR)/$(BINARY_NAME) tobereleased

count: build
	@$(BUILD_DIR)/$(BINARY_NAME) count

merge: build
	@$(BUILD_DIR)/$(BINARY_NAME) merge
