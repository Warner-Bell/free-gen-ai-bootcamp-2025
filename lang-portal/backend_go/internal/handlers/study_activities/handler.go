// Package study_activities handles study activity-related HTTP endpoints
package study_activities

import (
    "encoding/json"
    "net/http"
    "strconv"

    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/models"
    "github.com/go-chi/chi/v5"
)

type Handler struct {
    activityModel *models.StudyActivityModel
}

func NewHandler(activityModel *models.StudyActivityModel) *Handler {
    return &Handler{activityModel: activityModel}
}

type createActivityRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Type        string `json:"type"`
    Settings    json.RawMessage `json:"settings,omitempty"`
}

type updateActivityRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Type        string `json:"type"`
    Settings    json.RawMessage `json:"settings,omitempty"`
    Active      bool   `json:"active"`
}

// CreateActivity handles the creation of a new study activity
func (h *Handler) CreateActivity(w http.ResponseWriter, r *http.Request) {
    var req createActivityRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    activity, err := h.activityModel.Create(req.Name, req.Description, req.Type, req.Settings)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(activity)
}

// GetActivity retrieves a specific study activity by ID
func (h *Handler) GetActivity(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    activityID, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        http.Error(w, "Invalid activity ID", http.StatusBadRequest)
        return
    }

    activity, err := h.activityModel.GetByID(activityID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(activity)
}

// UpdateActivity updates an existing study activity
func (h *Handler) UpdateActivity(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    activityID, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        http.Error(w, "Invalid activity ID", http.StatusBadRequest)
        return
    }

    var req updateActivityRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    activity, err := h.activityModel.Update(activityID, req.Name, req.Description, req.Type, req.Settings, req.Active)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(activity)
}

// DeleteActivity deletes a study activity
func (h *Handler) DeleteActivity(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    activityID, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        http.Error(w, "Invalid activity ID", http.StatusBadRequest)
        return
    }

    if err := h.activityModel.Delete(activityID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// ListActivities retrieves all study activities
func (h *Handler) ListActivities(w http.ResponseWriter, r *http.Request) {
    activities, err := h.activityModel.GetAll()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(activities)
}

// GetActivityStats retrieves statistics for a specific activity
func (h *Handler) GetActivityStats(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    activityID, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        http.Error(w, "Invalid activity ID", http.StatusBadRequest)
        return
    }

    stats, err := h.activityModel.GetStats(activityID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(stats)
}
