// internal/models/word.go
package models

import (
    "database/sql"
    "time"
)

type Word struct {
    ID        int64     `json:"id"`
    Japanese  string    `json:"japanese"`
    Romaji    string    `json:"romaji"`
    English   string    `json:"english"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type WordModel struct {
    db *sql.DB
}

func NewWordModel(db *sql.DB) *WordModel {
    return &WordModel{db: db}
}

func (m *WordModel) GetAll(page, perPage int) ([]Word, error) {
    offset := (page - 1) * perPage
    rows, err := m.db.Query(`
        SELECT id, japanese, romaji, english, created_at, updated_at 
        FROM words 
        ORDER BY id 
        LIMIT ? OFFSET ?`, 
        perPage, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var words []Word
    for rows.Next() {
        var w Word
        err := rows.Scan(&w.ID, &w.Japanese, &w.Romaji, &w.English, &w.CreatedAt, &w.UpdatedAt)
        if err != nil {
            return nil, err
        }
        words = append(words, w)
    }
    return words, nil
}
