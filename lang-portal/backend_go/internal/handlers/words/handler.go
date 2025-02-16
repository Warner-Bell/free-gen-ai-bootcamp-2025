// Package words handles word-related HTTP endpoints
package words

import (
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/go-chi/chi"

    "backend_go/internal/models"
)

type Handler struct {
    wordModel *models.WordModel
}

func NewHandler(wordModel *models.WordModel) *Handler {
    return &Handler{wordModel: wordModel}
}

func (h *Handler) GetWords(w http.ResponseWriter, r *http.Request) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    words, err := h.wordModel.GetWords(page)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(words)
}

func (h *Handler) CreateWord(w http.ResponseWriter, r *http.Request) {
    var word models.Word
    if err := json.NewDecoder(r.Body).Decode(&word); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    if err := h.wordModel.CreateWord(&word); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(word)
}

func (h *Handler) GetWord(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    wordID, err := strconv.Atoi(id)
    if err != nil {
        http.Error(w, "Invalid word ID", http.StatusBadRequest)
        return
    }

    word, err := h.wordModel.GetWordByID(wordID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(word)
}

func (h *Handler) UpdateWord(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    wordID, err := strconv.Atoi(id)
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
    if err := h.wordModel.UpdateWord(&word); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(word)
}

func (h *Handler) DeleteWord(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    wordID, err := strconv.Atoi(id)
    if err != nil {
        http.Error(w, "Invalid word ID", http.StatusBadRequest)
        return
    }

    if err := h.wordModel.DeleteWord(wordID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}