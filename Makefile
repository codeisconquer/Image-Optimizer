.PHONY: build build-all clean install deps

# Standard Build für aktuelle Plattform
build:
	go build -o image-optimizer main.go

# Build für verschiedene Plattformen
build-all:
	@echo "Building for macOS (Apple Silicon)..."
	GOOS=darwin GOARCH=arm64 go build -o dist/image-optimizer-darwin-arm64 main.go
	@echo "Building for macOS (Intel)..."
	GOOS=darwin GOARCH=amd64 go build -o dist/image-optimizer-darwin-amd64 main.go
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o dist/image-optimizer-linux-amd64 main.go
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o dist/image-optimizer-windows-amd64.exe main.go

# Dependencies installieren
deps:
	go mod download
	go mod tidy

# Clean
clean:
	rm -f image-optimizer image-optimizer.exe
	rm -rf dist/

# Install lokal
install: build
	cp image-optimizer /usr/local/bin/

