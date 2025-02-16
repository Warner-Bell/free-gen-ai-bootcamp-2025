## Rebuild Instructions

After fixing code formatting and exports:

1. Save all files and run:
```bash
cd lang-portal/backend_go
go mod tidy
go clean -cache
mage build
```

This should resolve any undefined method errors by ensuring consistent code formatting and proper Go exports.