## Clean Build Instructions

To resolve module resolution issues:

1. Change to the backend_go directory:
```bash
cd lang-portal/backend_go
```

2. Remove build artifacts:
```bash
rm -f go.sum
go clean -modcache
go clean -cache
```

3. Regenerate dependencies:
```bash
go mod tidy
```

4. Rebuild:
```bash
mage build
```

This should force Go to recalculate all module paths and rebuild from scratch.