## Fix Module Dependencies

To ensure all code changes are properly saved and dependencies are resolved:

1. Change to the backend_go directory:
```bash
cd lang-portal/backend_go
```

2. Save all changes:
```bash
git add .
git commit -m "feat: add CRUD methods to word and group models"
```

3. Run go mod tidy to ensure dependencies are resolved:
```bash
go mod tidy
```

4. Clean and rebuild:
```bash
go clean -cache
mage build
```

Make sure all files are properly saved before running these commands.