// Package groups handles group-related HTTP endpoints
package groups

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi"
    "backend_go/internal/models"
)

type Handler struct {
    groupModel *models.GroupModel
}

func NewHandler(groupModel *models.GroupModel) *Handler {
    return &Handler{groupModel: groupModel}
}

func (h *Handler) GetGroups(w http.ResponseWriter, r *http.Request) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    groups, err := h.groupModel.GetGroups(page)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(groups)
}

func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
    var group models.Group
    if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.groupModel.CreateGroup(&group); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(group)
}

func (h *Handler) GetGroup(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    groupID, err := strconv.Atoi(id)
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    group, err := h.groupModel.GetGroupByID(groupID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(group)
}

func (h *Handler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    groupID, err := strconv.Atoi(id)
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    var group models.Group
    if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    group.ID = groupID
    if err := h.groupModel.UpdateGroup(&group); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(group)
}

func (h *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    groupID, err := strconv.Atoi(id)
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    if err := h.groupModel.DeleteGroup(groupID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AddWordToGroup(w http.ResponseWriter, r *http.Request) {
    groupID, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    var req struct {
        WordID int `json:"wordId"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.groupModel.AddWordToGroup(groupID, req.WordID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) RemoveWordFromGroup(w http.ResponseWriter, r *http.Request) {
    groupID, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    wordID, err := strconv.Atoi(chi.URLParam(r, "wordId"))
    if err != nil {
        http.Error(w, "Invalid word ID", http.StatusBadRequest)
        return
    }

    if err := h.groupModel.RemoveWordFromGroup(groupID, wordID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetGroupWords(w http.ResponseWriter, r *http.Request) {
    groupID, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    words, err := h.groupModel.GetGroupWords(groupID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(words)
}