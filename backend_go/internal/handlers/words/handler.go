package words

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"backend_go/internal/models"
)

type WordsHandler struct {
	db *sql.DB
}

func NewWordsHandler(db *sql.DB) *WordsHandler {
	return &WordsHandler{db: db}
}

func (h *WordsHandler) GetWords(w http.ResponseWriter, r *http.Request) {
	var words []models.Word
	rows, err := h.db.Query(`SELECT id, japanese, romaji, english, created_at, updated_at FROM words`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var word models.Word
		if err := rows.Scan(&word.ID, &word.Japanese, &word.Romaji, &word.English, &word.CreatedAt, &word.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		words = append(words, word)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(words)
}