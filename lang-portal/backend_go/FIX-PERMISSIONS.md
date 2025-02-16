## Fix File Permissions

The undefined method errors may be caused by file permission issues preventing Go from properly reading the source files. To fix:

1. Change to the backend_go directory:
```bash
cd lang-portal/backend_go
```

2. Fix file permissions:
```bash
chmod -R u+rw .
```

3. Clean build cache and dependencies:
```bash
go clean -cache
go clean -modcache
rm -f go.sum
```

4. Regenerate dependencies:
```bash
go mod tidy
```

5. Rebuild:
```bash
mage build
```

The code appears correct but Go may not be able to properly read the files due to permissions. These steps ensure all files are readable and writable by the current user.