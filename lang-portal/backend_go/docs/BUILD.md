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

## Code Changes Made
The following changes have been implemented to fix the build errors:

1. In `internal/handlers/words/handler.go`:
   - Changed `GetAll` to `GetWords` with correct pagination parameters
   - Updated to use `(page-1)*perPage` for proper offset calculation

2. In `internal/handlers/groups/handler.go`:
   - Fixed group creation to use proper `Group` struct
   - Updated return value handling for `Create` method

3. Added `internal/handlers/groups/types.go`:
   - Added type alias for proper Group type handling

## Troubleshooting

If you see "go: no modules specified":
- Make sure you're in the `lang-portal/backend_go` directory
- Verify go.mod exists
- Run `go mod download` before building

If build still fails, verify GOPATH and GOROOT are set correctly.