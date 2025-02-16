## Export Fix

The undefined method errors may be caused by unexported struct fields. Changes made:

1. Made model fields public:
   - Changed WordModel.db to WordModel.DB
   - Changed GroupModel.db to GroupModel.DB
   - Updated constructor functions

2. To apply:
```bash
cd lang-portal/backend_go
go clean -cache
mage build
```

This ensures proper visibility of model fields and methods.