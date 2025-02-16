# Handler Fix Notes

To fix the build errors:

1. Files should be created in `/mnt/h/Cloud-Lab/free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/handlers/`, where the go.mod file is located

2. Current errors:
   - `undefined: time` in words/handler.go
   - `undefined: Word` in groups/handler.go 

3. Solution:
   - Delete any files in the wrong backend_go directory
   - Create handlers in lang-portal/backend_go/internal/handlers/
   - Use correct imports relative to the module root (backend_go)
   - Keep all imports including "time" and "backend_go/internal/models"

4. Ensure all code is under the same module directory as go.mod for proper dependency resolution