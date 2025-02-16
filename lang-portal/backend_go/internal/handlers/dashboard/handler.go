package dashboard

import (
    "encoding/json"
    "net/http"

    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/models"
)

type Handler struct {
    model *models.DashboardModel
}

func NewHandler(model *models.DashboardModel) *Handler {
    return &Handler{model: model}
}

type DashboardResponse struct {
    TotalWords     int                    `json:"total_words"`
    TotalGroups    int                    `json:"total_groups"`
    RecentWords    []models.Word          `json:"recent_words"`
    RecentGroups   []models.Group         `json:"recent_groups"`
    RecentSessions []models.StudySession  `json:"recent_sessions"`
}

func (h *Handler) GetDashboard(w http.ResponseWriter, r *http.Request) {
    response := DashboardResponse{}
    // TODO: Implement dashboard data retrieval
    json.NewEncoder(w).Encode(response)
}
