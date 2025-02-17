## Build Instructions

To build the backend_go application:

1. Change to the backend_go directory:
```bash
cd lang-portal/backend_go
```

2. Run the build command:
```bash
mage build
```

Note: The build command must be run from the `lang-portal/backend_go` directory where the `go.mod` file is located.

## Check Files Included in Build

The undefined method errors could be caused by files not being included in the build. To check:

1. List all Go files that should be compiled:
```bash
cd lang-portal/backend_go
find . -name "*.go" ! -name "*_test.go"
```

2. Verify models directory is in GOPATH:
```bash
go list -f '{{.Dir}}' free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/models
```

3. Check module dependencies:
```bash
go mod why free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/models
```

This will help verify that all necessary files are being found and included in the build process.