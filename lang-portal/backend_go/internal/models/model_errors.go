package models

import "database/sql"

// Common errors
var (
	ErrNotFound = sql.ErrNoRows
)