// internal/models/word.go

package models

import (
    "database/sql"
)

// internal/models/word.go

type Word struct {
    ID          int       `json:"id"`
    Word        string    `json:"word"`
    Translation string    `json:"translation"`
    Notes       string    `json:"notes"`
    Japanese    string    `json:"japanese"`
    Romaji      string    `json:"romaji"`
    English     string    `json:"english"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}


type WordModel struct {
    db *sql.DB
}

func NewWordModel(db *sql.DB) *WordModel {
    return &WordModel{db: db}
}

// GetWords returns words with pagination
func (m *WordModel) GetWords(offset, limit int) ([]Word, int, error) {
    // First, get total count
    var total int
    err := m.db.QueryRow("SELECT COUNT(*) FROM words").Scan(&total)
    if err != nil {
        return nil, 0, err
    }

    // Then get the words for the current page
    rows, err := m.db.Query("SELECT id, word, translation, notes FROM words LIMIT ? OFFSET ?", limit, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var words []Word
    for rows.Next() {
        var w Word
        err := rows.Scan(&w.ID, &w.Word, &w.Translation, &w.Notes)
        if err != nil {
            return nil, 0, err
        }
        words = append(words, w)
    }

    return words, total, nil
}

// SearchWords searches words based on a query string
func (m *WordModel) SearchWords(query string, offset, limit int) ([]Word, int, error) {
    // First, get total count for the search
    var total int
    err := m.db.QueryRow(
        "SELECT COUNT(*) FROM words WHERE word LIKE ? OR translation LIKE ?",
        "%"+query+"%", "%"+query+"%",
    ).Scan(&total)
    if err != nil {
        return nil, 0, err
    }

    // Then get the matching words
    rows, err := m.db.Query(
        "SELECT id, word, translation, notes FROM words WHERE word LIKE ? OR translation LIKE ? LIMIT ? OFFSET ?",
        "%"+query+"%", "%"+query+"%", limit, offset,
    )
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var words []Word
    for rows.Next() {
        var w Word
        err := rows.Scan(&w.ID, &w.Word, &w.Translation, &w.Notes)
        if err != nil {
            return nil, 0, err
        }
        words = append(words, w)
    }

    return words, total, nil
}

// Add other methods like GetWord, CreateWord, UpdateWord, DeleteWord...
