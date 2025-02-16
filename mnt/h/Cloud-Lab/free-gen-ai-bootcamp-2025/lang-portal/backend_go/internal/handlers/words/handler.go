package words

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"
    "time"

    "backend_go/internal/models"
)

type Handler struct {
    db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
    return &Handler{db: db}
}

func (h *Handler) GetWords(w http.ResponseWriter, r *http.Request) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    if page < 1 {
        page = 1
    }
    
    perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
    if perPage < 1 {
        perPage = 10
    }

    model := models.NewWordModel(h.db)
    words, total, err := model.GetWords((page-1)*perPage, perPage)
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

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}