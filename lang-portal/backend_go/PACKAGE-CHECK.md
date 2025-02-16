## Go Package Build Check

The undefined method errors suggest the models package is not being compiled correctly. To verify:

1. Create a package check file in models directory to verify package builds:

Create file: internal/models/package_doc.go
```go
// Package models defines the data models and database operations
package models

// Export these methods to ensure they are compiled
var (
    _ = (*WordModel)(nil).Create
    _ = (*WordModel)(nil).GetByID
    _ = (*WordModel)(nil).Update
    _ = (*WordModel)(nil).Delete
    _ = (*GroupModel)(nil).Update
    _ = (*GroupModel)(nil).Delete
)
```

2. Clean build cache:
```bash
cd lang-portal/backend_go
go clean -cache
```

3. Rebuild checking all packages:
```bash
go build ./...
```

This ensures the models package is properly compiled with all required methods.