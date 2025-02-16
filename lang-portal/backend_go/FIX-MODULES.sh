#!/bin/bash

# Reset Go modules and rebuild with new module path
cd lang-portal/backend_go

# Remove all cached modules and build artifacts
go clean -modcache
go clean -cache
rm -f go.sum

# Clear any vendor directory if it exists
rm -rf vendor/

# Reinitialize modules with new path
go mod edit -module github.com/warner/lang-portal/backend_go
go mod tidy

# Clean build cache and rebuild
mage build