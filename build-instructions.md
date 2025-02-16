# Build Instructions

To build the project successfully:

1. Navigate to the backend_go directory:
```bash
cd lang-portal/backend_go
```

2. Run the mage build command:
```bash
mage build
```

The code changes have been made to fix all build errors:
- Fixed undefined GetAll method by using GetWords with correct pagination
- Fixed assignment mismatch in groupModel.Create
- Fixed type mismatch for group creation by using proper Group struct

Make sure to run the build command from the correct directory containing the go.mod file.