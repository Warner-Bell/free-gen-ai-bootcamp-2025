package models

import (
    "database/sql"
    "time"
)

type Word struct {
    ID          int64     `json:"id"`
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
    DB *sql.DB
}

func NewWordModel(db *sql.DB) *WordModel {
    return &WordModel{DB: db}
}

func (m *WordModel) GetWords(offset, limit int) ([]Word, int, error) {
    var total int
    err := m.DB.QueryRow("SELECT COUNT(*) FROM words").Scan(&total)
    if err != nil {
        return nil, 0, err
    }

    rows, err := m.DB.Query("SELECT id, word, translation, notes, japanese, romaji, english, created_at, updated_at FROM words ORDER BY id LIMIT ? OFFSET ?", limit, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var words []Word
    for rows.Next() {
        var w Word
        err := rows.Scan(&w.ID, &w.Word, &w.Translation, &w.Notes, &w.Japanese, &w.Romaji, &w.English, &w.CreatedAt, &w.UpdatedAt)
        if err != nil {
            return nil, 0, err
        }
        words = append(words, w)
    }
    return words, total, nil
}

func (m *WordModel) Create(word *Word) error {
    query := "INSERT INTO words (word, translation, notes, japanese, romaji, english, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, DATETIME('now'), DATETIME('now')) RETURNING id"
    
    err := m.DB.QueryRow(query,
        word.Word,
        word.Translation,
        word.Notes,
        word.Japanese,
        word.Romaji,
        word.English).Scan(&word.ID)
    return err
}

func (m *WordModel) GetByID(id int64) (*Word, error) {
    word := &Word{}
    query := "SELECT id, word, translation, notes, japanese, romaji, english, created_at, updated_at FROM words WHERE id = ?"
    
    err := m.DB.QueryRow(query, id).Scan(
        &word.ID,
        &word.Word,
        &word.Translation,
        &word.Notes,
        &word.Japanese,
        &word.Romaji,
        &word.English,
        &word.CreatedAt,
        &word.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, err
    }
    return word, err
}

func (m *WordModel) Update(word *Word) error {
    query := "UPDATE words SET word = ?, translation = ?, notes = ?, japanese = ?, romaji = ?, english = ?, updated_at = DATETIME('now') WHERE id = ?"
    
    result, err := m.DB.Exec(query,
        word.Word,
        word.Translation,
        word.Notes,
        word.Japanese,
        word.Romaji,
        word.English,
        word.ID)
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }
    return nil
}

func (m *WordModel) Delete(id int64) error {
    query := "DELETE FROM words WHERE id = ?"
    result, err := m.DB.Exec(query, id)
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }
    return nil
}

func (m *WordModel) SearchWords(query string, offset, limit int) ([]Word, int, error) {
    var total int
    err := m.DB.QueryRow("SELECT COUNT(*) FROM words WHERE word LIKE ? OR translation LIKE ? OR notes LIKE ?",
        "%"+query+"%", "%"+query+"%", "%"+query+"%").Scan(&total)
    if err != nil {
        return nil, 0, err
    }

    rows, err := m.DB.Query("SELECT id, word, translation, notes, japanese, romaji, english, created_at, updated_at FROM words WHERE word LIKE ? OR translation LIKE ? OR notes LIKE ? ORDER BY id LIMIT ? OFFSET ?",
        "%"+query+"%", "%"+query+"%", "%"+query+"%", limit, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var words []Word
    for rows.Next() {
        var w Word
        err := rows.Scan(&w.ID, &w.Word, &w.Translation, &w.Notes, &w.Japanese, &w.Romaji, &w.English, &w.CreatedAt, &w.UpdatedAt)
        if err != nil {
            return nil, 0, err
        }
        words = append(words, w)
    }
    return words, total, nil
}
