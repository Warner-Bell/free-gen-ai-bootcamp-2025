## Complete Module Reset Instructions

The undefined method errors are likely due to cached module information. The module path was changed from `backend_go` to `github.com/warner/lang-portal/backend_go`, which requires a complete reset:

1. Make the fix script executable and run it:
```bash
chmod +x FIX-MODULES.sh
./FIX-MODULES.sh
```

Or run these commands manually:

1. Change to the backend directory:
```bash
cd lang-portal/backend_go
```

2. Remove all module and build caches:
```bash
go clean -modcache
go clean -cache
rm -f go.sum
rm -rf vendor/
```

3. Reinitialize module:
```bash
go mod edit -module github.com/warner/lang-portal/backend_go
go mod tidy
```

4. Rebuild:
```bash
mage build
```

This ensures a complete reset of Go's module system to recognize the new module path and exposed methods.