// internal/handlers/groups/handler.go
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

type createGroupRequest struct {
    Name string `json:"name"`
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
    var req createGroupRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    group, err := h.groupModel.Create(req.Name)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(group)
}

func (h *Handler) GetGroup(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    group, err := h.groupModel.GetByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(group)
}



type addWordRequest struct {
    WordID int64 `json:"word_id"`
}

func (h *Handler) AddWordToGroup(w http.ResponseWriter, r *http.Request) {
    groupID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    var req addWordRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.groupModel.AddWord(groupID, req.WordID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func (h *Handler) RemoveWordFromGroup(w http.ResponseWriter, r *http.Request) {
    groupID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    wordID, err := strconv.ParseInt(chi.URLParam(r, "wordId"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid word ID", http.StatusBadRequest)
        return
    }

    if err := h.groupModel.RemoveWord(groupID, wordID); err != nil {
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
