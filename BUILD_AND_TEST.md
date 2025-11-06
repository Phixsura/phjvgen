# Build and Test Guide for phjvgen

This document provides instructions for building, testing, and using the phjvgen Go binary.

## Prerequisites

- Go 1.24 or later
- Git
- Maven 3.x (for testing generated projects)
- MySQL 8.0+ (optional, for testing demo code)

## Building

### Quick Build

```bash
cd /Users/bytedance/develop/phix_project/template/phjvgen
go build -o phjvgen .
```

### Using Makefile

```bash
# Build optimized binary
make build

# Build for all platforms (Linux, macOS, Windows)
make build-all

# Development build (no optimization)
make dev

# Clean build artifacts
make clean
```

### Manual Build with Optimization

```bash
# Optimized build (smaller binary)
go build -ldflags="-s -w" -o phjvgen .

# Check binary size
ls -lh phjvgen
```

## Testing the Binary

### 1. Version Check

```bash
./phjvgen version
# Output: phjvgen version 1.0.0
```

### 2. Help Command

```bash
./phjvgen --help
./phjvgen generate --help
./phjvgen demo --help
./phjvgen add --help
```

### 3. Generate Example Project

```bash
# Quick example (no interaction)
./phjvgen example

# Verify generated files
ls -la demo-app/
tree demo-app/  # if tree is installed
```

### 4. Generate Custom Project

```bash
# Interactive generation
./phjvgen generate

# Sample inputs:
# Group ID: com.test
# Artifact ID: test-app
# Version: 1.0.0
# Project Name: Test Application
# Description: Test project
# Output Directory: ./test-app
```

### 5. Test Demo Code Generation

```bash
# Go to generated project
cd demo-app

# Generate CRUD demo
../phjvgen demo

# Verify generated files
find . -name "User*.java"
find . -name "*user*"
```

### 6. Test Module Addition

```bash
# Still in demo-app directory
../phjvgen add payment

# Verify new module
ls -la application/application-payment/
cat pom.xml | grep application-payment
```

### 7. Build Generated Project

```bash
# Build the generated Java project
cd demo-app
mvn clean install

# Check build artifacts
ls -la starter/target/
```

### 8. Run Generated Project

```bash
# Run without database (health check only)
java --enable-preview -jar starter/target/starter-1.0.0.jar

# In another terminal, test the health endpoint
curl http://localhost:8080/api/health
```

## Installation

### Method 1: Using the install command

```bash
./phjvgen install
```

This will:
- Copy the binary to `~/.local/bin/phjvgen` (or `/usr/local/bin` if root)
- Print PATH configuration instructions if needed

### Method 2: Manual installation

```bash
# Copy to user bin directory
mkdir -p ~/.local/bin
cp phjvgen ~/.local/bin/
chmod +x ~/.local/bin/phjvgen

# Add to PATH (if not already)
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# Verify installation
which phjvgen
phjvgen version
```

## Full Integration Test

This test creates a complete project and verifies all features:

```bash
#!/bin/bash
set -e

echo "=== Full Integration Test for phjvgen ==="

# 1. Build phjvgen
echo "Building phjvgen..."
go build -o phjvgen .

# 2. Generate example project
echo "Generating example project..."
./phjvgen example

# 3. Verify structure
echo "Verifying project structure..."
test -f demo-app/pom.xml || (echo "Parent POM missing" && exit 1)
test -d demo-app/common || (echo "Common module missing" && exit 1)
test -d demo-app/domain || (echo "Domain module missing" && exit 1)
test -d demo-app/infrastructure || (echo "Infrastructure module missing" && exit 1)
test -d demo-app/adapter || (echo "Adapter module missing" && exit 1)
test -d demo-app/application || (echo "Application module missing" && exit 1)
test -d demo-app/starter || (echo "Starter module missing" && exit 1)

# 4. Enter project and generate demo
echo "Generating CRUD demo..."
cd demo-app
../phjvgen demo <<EOF
y
EOF

# 5. Verify demo files
echo "Verifying demo files..."
test -f domain/src/main/java/com/example/demo/domain/model/User.java || exit 1
test -f infrastructure/src/main/resources/db/migration/V1__create_user_table.sql || exit 1
test -f application/application-user/src/main/java/com/example/demo/application/user/service/UserService.java || exit 1

# 6. Add new module
echo "Adding payment module..."
../phjvgen add payment <<EOF
y
EOF

# 7. Verify new module
test -d application/application-payment || exit 1
test -f application/application-payment/pom.xml || exit 1
grep -q "application-payment" pom.xml || exit 1

# 8. Build project
echo "Building Maven project..."
mvn clean install -DskipTests

# 9. Verify build
test -f starter/target/starter-1.0.0.jar || (echo "Build failed - JAR not found" && exit 1)

echo "=== All tests passed! ==="
cd ..
```

Save this as `integration-test.sh` and run:

```bash
chmod +x integration-test.sh
./integration-test.sh
```

## Common Issues and Solutions

### Issue: "no Go files" error

**Solution**: Make sure you're in the phjvgen directory when building:

```bash
cd /path/to/phjvgen
go build -o phjvgen .
```

### Issue: Missing dependencies

**Solution**: Download and tidy modules:

```bash
go mod download
go mod tidy
go build -o phjvgen .
```

### Issue: Permission denied when installing

**Solution**: Use sudo or install to user directory:

```bash
# Install to user directory
mkdir -p ~/.local/bin
cp phjvgen ~/.local/bin/

# Or use sudo for system-wide install
sudo ./phjvgen install
```

### Issue: Binary not found after install

**Solution**: Check and update PATH:

```bash
# Check if install directory is in PATH
echo $PATH | grep -q "$HOME/.local/bin" || echo "Not in PATH"

# Add to PATH
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

## Performance Benchmarks

Compare performance with original shell scripts:

```bash
# Time phjvgen
time ./phjvgen example

# Time original shell script (if available)
time ./script/example.sh

# Typical results:
# phjvgen: ~0.1s
# shell script: ~0.5s
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...
make test-coverage

# Run specific package tests
go test ./internal/generator/...
```

### Code Quality

```bash
# Format code
go fmt ./...
make fmt

# Vet code
go vet ./...
make vet

# Lint (requires golangci-lint)
make lint
```

### Debugging

```bash
# Run with verbose output
./phjvgen generate -v  # if verbose flag added

# Use delve debugger
dlv debug . -- generate
```

## Cross-Platform Builds

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o phjvgen-linux .

# macOS AMD64 (Intel)
GOOS=darwin GOARCH=amd64 go build -o phjvgen-darwin-amd64 .

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o phjvgen-darwin-arm64 .

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o phjvgen.exe .

# Or use make
make build-all
```

## Distribution

### Create Release Package

```bash
#!/bin/bash
VERSION="1.0.0"

# Create release directory
mkdir -p release

# Build for all platforms
make build-all

# Package for each platform
cd build

# Linux
tar -czf ../release/phjvgen-${VERSION}-linux-amd64.tar.gz phjvgen-linux-amd64
tar -czf ../release/phjvgen-${VERSION}-darwin-amd64.tar.gz phjvgen-darwin-amd64
tar -czf ../release/phjvgen-${VERSION}-darwin-arm64.tar.gz phjvgen-darwin-arm64
zip ../release/phjvgen-${VERSION}-windows-amd64.zip phjvgen-windows-amd64.exe

cd ..
ls -lh release/
```

## Next Steps

After successful build and test:

1. Install to system: `./phjvgen install`
2. Generate your first project: `phjvgen generate`
3. Read the generated README: `cat <project-name>/README.md`
4. Start development!

## Support

For issues or questions:
- Check the main README.md
- Review this build guide
- Check Go and Maven versions
- Verify PATH configuration
