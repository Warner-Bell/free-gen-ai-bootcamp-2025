package settings

import (
    "encoding/json"
    "net/http"

    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/models"
)

type Handler struct {
    sessionModel  *models.StudySessionModel
    activityModel *models.StudyActivityModel
}

func NewHandler(sessionModel *models.StudySessionModel, activityModel *models.StudyActivityModel) *Handler {
    return &Handler{
        sessionModel:  sessionModel,
        activityModel: activityModel,
    }
}

type Settings struct {
    DefaultSessionDuration int  `json:"default_session_duration"`
    ShowTimer             bool `json:"show_timer"`
    AutoAdvance           bool `json:"auto_advance"`
}

func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) {
    settings := Settings{
        DefaultSessionDuration: 30,
        ShowTimer:             true,
        AutoAdvance:           false,
    }
    json.NewEncoder(w).Encode(settings)
}

func (h *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
    var settings Settings
    if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    json.NewEncoder(w).Encode(settings)
}
