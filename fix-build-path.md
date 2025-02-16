# Build Path Fix Required

## Issue
The Go build is failing because handler methods can't be found, even though they are correctly
implemented. This is happening because files are being saved to the wrong location.

## Current Paths
- Files are being written to: `lang-portal/backend_go/`
- But Go is looking in: `/mnt/h/Cloud-Lab/free-gen-ai-bootcamp-2025/lang-portal/backend_go/`

## Required Fix
The handler files need to be moved/copied to the correct build path:

```bash
# Move/copy handler files to the correct location:
cp -r lang-portal/backend_go/internal/handlers/* /mnt/h/Cloud-Lab/free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/handlers/

# Or update GOPATH/workspace settings to point to the correct location:
export GOPATH=/path/to/workspace
cd /mnt/h/Cloud-Lab/free-gen-ai-bootcamp-2025/lang-portal/backend_go
go build ./...
```

## Affected Files
The following files need to be in the correct location:
1. internal/handlers/words/handler.go
2. internal/handlers/groups/handler.go 
3. internal/handlers/sessions/handler.go

## Verification
You can verify file locations with:
```bash
find /mnt/h/Cloud-Lab/free-gen-ai-bootcamp-2025/lang-portal/backend_go -name "handler.go"
```

All handler files must exist in /mnt/h/Cloud-Lab/free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/handlers/[module]/handler.go