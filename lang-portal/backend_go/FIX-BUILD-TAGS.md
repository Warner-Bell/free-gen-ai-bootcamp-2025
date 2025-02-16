## Fix Build Tags and Compilation

The undefined method errors suggest files may not be compiling properly. To fix:

1. Check build constraints - ensure no files have build tags preventing compilation:
```bash
cd lang-portal/backend_go
find . -name "*.go" -exec grep -l "^//go:build" {} \;
```

2. If any files have build constraints, verify they are correct for your environment.

3. Build with verbose output to see which files are being compiled:
```bash
go build -v ./...
```

4. Clean and rebuild:
```bash
go clean -cache
go clean -modcache
mage build
```

The methods exist in the code but may not be included in compilation due to build constraints. These steps will verify all files are being properly included.