## Code Formatting Fix

The undefined method errors may be caused by inconsistent code formatting. To fix:

1. Format all Go files:
```bash
cd lang-portal/backend_go
gofmt -w .
```

2. Clean and rebuild:
```bash
go clean -cache
mage build
```

This ensures consistent formatting using tabs instead of spaces, which is Go's standard.