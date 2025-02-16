## Type Consistency Fix

The undefined method errors may be caused by type inconsistencies between models. Changes made:

1. Standardized timestamp types:
   - Changed Group model timestamps from string to time.Time
   - Now matches Word model timestamp types
   - Ensures consistent type handling across models

2. To apply changes:
```bash
cd lang-portal/backend_go
go clean -cache
mage build
```

This ensures consistent types are used across all models.