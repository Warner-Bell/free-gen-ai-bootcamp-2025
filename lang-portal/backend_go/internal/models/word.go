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
	db *sql.DB
}

func NewWordModel(db *sql.DB) *WordModel {
	return &WordModel{db: db}
}

func (m *WordModel) GetWords(offset, limit int) ([]Word, int, error) {
	var total int
	err := m.db.QueryRow("SELECT COUNT(*) FROM words").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := m.db.Query(`
		SELECT id, word, translation, notes, japanese, romaji, english, created_at, updated_at 
		FROM words 
		ORDER BY id 
		LIMIT ? OFFSET ?`,
		limit, offset)
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

func (m *WordModel) SearchWords(query string, offset, limit int) ([]Word, int, error) {
	var total int
	err := m.db.QueryRow(`
		SELECT COUNT(*) FROM words 
		WHERE word LIKE ? OR translation LIKE ? OR notes LIKE ?`,
		"%"+query+"%", "%"+query+"%", "%"+query+"%").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := m.db.Query(`
		SELECT id, word, translation, notes, japanese, romaji, english, created_at, updated_at
		FROM words
		WHERE word LIKE ? OR translation LIKE ? OR notes LIKE ?
		ORDER BY id
		LIMIT ? OFFSET ?`,
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