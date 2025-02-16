// Package models defines the data models and database operations
package models

// Export these methods to ensure they are compiled
var (
    _ = (*WordModel)(nil).Create
    _ = (*WordModel)(nil).GetByID
    _ = (*WordModel)(nil).Update
    _ = (*WordModel)(nil).Delete
    _ = (*GroupModel)(nil).Update
    _ = (*GroupModel)(nil).Delete
)