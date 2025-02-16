// internal/handlers/dashboard/handler.go
package dashboard

import (
    "encoding/json"
    "net/http"

    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/models"
)

type Handler struct {
    dashboardModel *models.DashboardModel
}

func NewHandler(dashboardModel *models.DashboardModel) *Handler {
    return &Handler{dashboardModel: dashboardModel}
}

func (h *Handler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
    stats, err := h.dashboardModel.GetStats()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}

func (h *Handler) GetRecentSessions(w http.ResponseWriter, r *http.Request) {
    sessions, err := h.dashboardModel.GetRecentSessions(10) // Get last 10 sessions
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "sessions": sessions,
    })
}

func (h *Handler) GetLearningProgress(w http.ResponseWriter, r *http.Request) {
    progress, err := h.dashboardModel.GetLearningProgress()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(progress)
}
