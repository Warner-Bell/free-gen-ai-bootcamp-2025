# Build Instructions

## Prerequisites
- Go 1.20 or later
- Mage build tool

## Build Steps

1. Ensure you are in the correct directory structure:
```bash
cd lang-portal/backend_go
```

2. Verify the go.mod file is present:
```bash
ls go.mod
```

3. If needed, download dependencies:
```bash
go mod download
```

4. Run the build:
```bash
mage build
```

# Start the server
mage run

# Stop the server
Ctrl+C (in the terminal where server is running)

# Clean build artifacts
mage clean

# Tidy dependencies
go mod tidy

# Full rebuild (clean and build)
mage clean && mage build

# Reset database (not implemented)
mage resetdb

# Run migrations (not implemented)
mage migrate

# Run API endpoint tests
./test-api-endpoints.sh

# Run unit tests (not implemented)
mage test
