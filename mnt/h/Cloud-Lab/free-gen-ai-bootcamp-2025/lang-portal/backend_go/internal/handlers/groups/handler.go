package groups

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"

    "backend_go/internal/models"
)

type Handler struct {
    db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
    return &Handler{db: db}
}

func (h *Handler) GetGroupWords(w http.ResponseWriter, r *http.Request) {
    groupID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    model := models.NewGroupModel(h.db)
    var words []models.Word
    words, err = model.GetGroupWords(groupID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(words)
}