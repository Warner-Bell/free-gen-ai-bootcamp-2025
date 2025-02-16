## Fix Build Cache Instructions

To resolve the "undefined method" errors, try clearing Go's build cache before rebuilding:

1. Change to the backend_go directory:
```bash
cd lang-portal/backend_go
```

2. Clear Go's build cache:
```bash
go clean -cache
```

3. Then run the build command:
```bash
mage build
```

These steps will ensure Go rebuilds everything from scratch, picking up the newly added methods.