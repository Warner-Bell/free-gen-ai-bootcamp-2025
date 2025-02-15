// internal/handlers/words/handler.go
package words

import (
    "encoding/json"
    "net/http"
    "strconv"

    "lang-portal/internal/models"
)

type Handler struct {
    wordModel *models.WordModel
}

func NewHandler(wordModel *models.WordModel) *Handler {
    return &Handler{wordModel: wordModel}
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

    words, err := h.wordModel.GetAll(page, perPage)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "words": words,
        "pagination": map[string]int{
            "current_page": page,
            "per_page":    perPage,
        },
    })
}
