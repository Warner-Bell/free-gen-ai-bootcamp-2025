package words

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/models"
)

type Handler struct {
    model *models.WordModel
}

func NewHandler(model *models.WordModel) *Handler {
    return &Handler{model: model}
}

func (h *Handler) GetWords(w http.ResponseWriter, r *http.Request) {
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

    if limit == 0 {
        limit = 10
    }

    words, total, err := h.model.GetWords(offset, limit)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response := struct {
        Words []models.Word `json:"words"`
        Total int          `json:"total"`
    }{
        Words: words,
        Total: total,
    }

    json.NewEncoder(w).Encode(response)
}

func (h *Handler) SearchWords(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

    if limit == 0 {
        limit = 10
    }

    words, total, err := h.model.SearchWords(query, offset, limit)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response := struct {
        Words []models.Word `json:"words"`
        Total int          `json:"total"`
    }{
        Words: words,
        Total: total,
    }

    json.NewEncoder(w).Encode(response)
}

func (h *Handler) CreateWord(w http.ResponseWriter, r *http.Request) {
    var word models.Word
    if err := json.NewDecoder(r.Body).Decode(&word); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.model.Create(&word); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(word)
}

func (h *Handler) GetWord(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    word, err := h.model.GetByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if word == nil {
        http.NotFound(w, r)
        return
    }

    json.NewEncoder(w).Encode(word)
}

func (h *Handler) UpdateWord(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var word models.Word
    if err := json.NewDecoder(r.Body).Decode(&word); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    word.ID = id
    if err := h.model.Update(&word); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(word)
}

func (h *Handler) DeleteWord(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    if err := h.model.Delete(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
