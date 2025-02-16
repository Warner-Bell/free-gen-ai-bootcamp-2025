## Fix Build Instructions

The build errors appear because the build command is not being run from the correct directory. To fix:

1. Change to the correct directory that contains go.mod:
```bash
cd lang-portal/backend_go
```

2. Then run the build command:
```bash
mage build
```

The current command is being run from `/mnt/h/Cloud-Lab/free-gen-ai-bootcamp-2025/lang-portal/backend_go` which may not be resolving the module paths correctly.

Ensure you are in the directory containing go.mod before running the build command.