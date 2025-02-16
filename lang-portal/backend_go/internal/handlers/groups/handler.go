package groups

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "backend_go/internal/models"
)

type Handler struct {
    groupModel *models.GroupModel
}

func NewHandler(groupModel *models.GroupModel) *Handler {
    return &Handler{groupModel: groupModel}
}

func (h *Handler) GetGroups(w http.ResponseWriter, r *http.Request) {
    groups, err := h.groupModel.GetAll()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "groups": groups,
    })
}

func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name string `json:"name"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    group := &Group{
        Name: req.Name,
    }
    err := h.groupModel.Create(group)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(group)
}

func (h *Handler) AddWordToGroup(w http.ResponseWriter, r *http.Request) {
    groupID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    var req struct {
        WordID int64 `json:"word_id"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.groupModel.AddWord(groupID, req.WordID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) RemoveWordFromGroup(w http.ResponseWriter, r *http.Request) {
    groupID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    var req struct {
        WordID int64 `json:"word_id"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.groupModel.RemoveWord(groupID, req.WordID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetGroupWords(w http.ResponseWriter, r *http.Request) {
    groupID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    words, err := h.groupModel.GetGroupWords(groupID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "words": words,
    })
}