// Package words handles word-related HTTP endpoints
package words

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/warner/lang-portal/backend_go/internal/models"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	wordModel *models.WordModel
}

func NewHandler(wordModel *models.WordModel) *Handler {
	return &Handler{wordModel: wordModel}
}

func (h *Handler) GetWords(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}
	offset := page * limit
	
	words, total, err := h.wordModel.GetWords(offset, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"words": words,
		"total": total,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) CreateWord(w http.ResponseWriter, r *http.Request) {
	var word models.Word
	if err := json.NewDecoder(r.Body).Decode(&word); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.wordModel.Create(&word); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(word)
}

func (h *Handler) GetWord(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	wordID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Invalid word ID", http.StatusBadRequest)
		return
	}

	word, err := h.wordModel.GetByID(wordID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(word)
}

func (h *Handler) UpdateWord(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	wordID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Invalid word ID", http.StatusBadRequest)
		return
	}

	var word models.Word
	if err := json.NewDecoder(r.Body).Decode(&word); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	word.ID = wordID

	if err := h.wordModel.Update(&word); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(word)
}

func (h *Handler) DeleteWord(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	wordID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Invalid word ID", http.StatusBadRequest)
		return
	}

	if err := h.wordModel.Delete(wordID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}