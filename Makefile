# Variables
BINARY_NAME=journl
INSTALL_DIR=$(HOME)/.local/bin
GO_FILES=$(shell find . -type f -name '*.go')

.PHONY: all build install clean

# Default target
all: build

# Build the binary in the current directory
build: $(GO_FILES)
	go build -o $(BINARY_NAME) main.go

# Build and copy to .local/bin
install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Successfully installed $(BINARY_NAME) to $(INSTALL_DIR)"

# Clean up local binary
clean:
	rm -f $(BINARY_NAME)
