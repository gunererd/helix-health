.PHONY: build release clean install help

BINARY_NAME=helix-health
VERSION?=1.0.0

help:
	@echo "Available targets:"
	@echo "  make build     - Build for current platform"
	@echo "  make release   - Build for all platforms (Linux, macOS, Windows)"
	@echo "  make install   - Install to /usr/local/bin"
	@echo "  make clean     - Remove built binaries"

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME)
	@echo "Done! Run with ./$(BINARY_NAME)"

release:
	@echo "Building release binaries..."
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY_NAME)-linux-amd64
	GOOS=linux GOARCH=arm64 go build -o dist/$(BINARY_NAME)-linux-arm64
	GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY_NAME)-macos-intel
	GOOS=darwin GOARCH=arm64 go build -o dist/$(BINARY_NAME)-macos-arm64
	GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY_NAME)-windows-amd64.exe
	@echo "Release binaries built in dist/"
	@ls -lh dist/

install: build
	@echo "Installing to /usr/local/bin..."
	sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "Installed! Run with: $(BINARY_NAME)"

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -rf dist/
	@echo "Clean!"
