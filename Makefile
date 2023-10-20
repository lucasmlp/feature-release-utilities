.PHONY: build run clean

# Name of the binary we'll produce
BINARY_NAME=countCharts

all: build run

# Build the Go code into a binary
build:
	@echo "Building..."
	@go build -o $(BINARY_NAME)

# Run the binary
run: build
	@echo "Running..."
	@./$(BINARY_NAME)

# Clean up the built binary
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
