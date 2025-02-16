// internal/models/word.go
package models

import (
    "database/sql"
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

func (m *WordModel) GetAll(page, perPage int) ([]Word, int, error) {
    // First get total count
    var totalCount int
    err := m.db.QueryRow("SELECT COUNT(*) FROM words").Scan(&totalCount)
    if err != nil {
        return nil, 0, err
    }

    // Calculate offset
    offset := (page - 1) * perPage

    // Get paginated results
    query := `
        SELECT id, japanese, romaji, english, created_at, updated_at 
        FROM words 
        ORDER BY id 
        LIMIT ? OFFSET ?`
    
    rows, err := m.db.Query(query, perPage, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var words []Word
    for rows.Next() {
        var word Word
        err := rows.Scan(
            &word.ID,
            &word.Japanese,
            &word.Romaji,
            &word.English,
            &word.CreatedAt,
            &word.UpdatedAt,
        )
        if err != nil {
            return nil, 0, err
        }
        words = append(words, word)
    }

    if err = rows.Err(); err != nil {
        return nil, 0, err
    }

    return words, totalCount, nil
}
