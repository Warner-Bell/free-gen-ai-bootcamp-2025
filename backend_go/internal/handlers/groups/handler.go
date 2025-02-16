package groups

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"backend_go/internal/models"
)

type GroupsHandler struct {
	db *sql.DB
}

func NewGroupsHandler(db *sql.DB) *GroupsHandler {
	return &GroupsHandler{db: db}
}

func (h *GroupsHandler) GetGroupWords(groupID int64) ([]models.Word, error) {
	query := `
		SELECT w.id, w.japanese, w.romaji, w.english, w.created_at, w.updated_at
		FROM words w
		JOIN word_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
		ORDER BY w.id`
	
	rows, err := h.db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var word models.Word
		err := rows.Scan(
			&word.ID,
			&word.Japanese,
			&word.Romaji,
			&word.English,
			&word.CreatedAt,
			&word.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		words = append(words, word)
	}
	
	return words, nil
}

func (h *GroupsHandler) GetAllWords() ([]models.Word, error) {
	query := `SELECT id, japanese, romaji, english, created_at, updated_at FROM words`
	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var word models.Word
		err := rows.Scan(
			&word.ID,
			&word.Japanese,
			&word.Romaji,
			&word.English,
			&word.CreatedAt,
			&word.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		words = append(words, word)
	}
	return words, nil
}

func (h *GroupsHandler) AddWordToGroup(groupID int64, word models.Word) error {
	query := `INSERT INTO word_groups (group_id, word_id) VALUES (?, ?)`
	_, err := h.db.Exec(query, groupID, word.ID)
	return err
}

func (h *GroupsHandler) RemoveWordFromGroup(groupID int64, wordID int64) error {
	query := `DELETE FROM word_groups WHERE group_id = ? AND word_id = ?`
	_, err := h.db.Exec(query, groupID, wordID)
	return err
}

func (h *GroupsHandler) UpdateWordInGroup(groupID int64, word models.Word) error {
	// First ensure the word is in the group
	checkQuery := `SELECT 1 FROM word_groups WHERE group_id = ? AND word_id = ?`
	var exists bool
	err := h.db.QueryRow(checkQuery, groupID, word.ID).Scan(&exists)
	if err != nil {
		return err
	}
	
	// Update the word's properties
	updateQuery := `UPDATE words SET japanese = ?, romaji = ?, english = ?, updated_at = ? WHERE id = ?`
	_, err = h.db.Exec(updateQuery, word.Japanese, word.Romaji, word.English, time.Now(), word.ID)
	return err
}