// Package sessions handles study session-related HTTP endpoints
package sessions

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"

    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/models"
    "github.com/go-chi/chi/v5"
)

type Handler struct {
    sessionModel *models.StudySessionModel
}

func NewHandler(sessionModel *models.StudySessionModel) *Handler {
    return &Handler{sessionModel: sessionModel}
}

type createSessionRequest struct {
    GroupID      int64  `json:"group_id"`
    ActivityName string `json:"activity_name"`
    Settings     json.RawMessage `json:"settings,omitempty"`
}

type addReviewRequest struct {
    WordID int64 `json:"word_id"`
    Known  bool  `json:"known"`
    Time   int   `json:"time_ms,omitempty"` // Response time in milliseconds
}

type sessionResponse struct {
    ID           int64           `json:"id"`
    GroupID      int64           `json:"group_id"`
    ActivityName string         `json:"activity_name"`
    StartTime    time.Time      `json:"start_time"`
    EndTime      *time.Time     `json:"end_time,omitempty"`
    Settings     json.RawMessage `json:"settings,omitempty"`
    Stats        *sessionStats   `json:"stats,omitempty"`
}

type sessionStats struct {
    TotalWords     int     `json:"total_words"`
    KnownWords     int     `json:"known_words"`
    UnknownWords   int     `json:"unknown_words"`
    AverageTimeMs  float64 `json:"average_time_ms,omitempty"`
    CompletionRate float64 `json:"completion_rate"`
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
    var req createSessionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    session, err := h.sessionModel.Create(req.GroupID, req.ActivityName)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(session)
}

func (h *Handler) EndSession(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    sessionID, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }

    if err := h.sessionModel.End(sessionID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) AddReview(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    sessionID, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }

    var req addReviewRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.sessionModel.AddReview(sessionID, req.WordID, req.Known); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetSessionStats(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    sessionID, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }

    stats, err := h.sessionModel.GetSessionStats(sessionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(stats)
}

func (h *Handler) GetSessions(w http.ResponseWriter, r *http.Request) {
    sessions, err := h.sessionModel.GetAll()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(sessions)
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    sessionID, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }

    session, err := h.sessionModel.GetByID(sessionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(session)
}

func (h *Handler) GetSessionsByGroup(w http.ResponseWriter, r *http.Request) {
    groupID := chi.URLParam(r, "groupId")
    groupIDInt, err := strconv.ParseInt(groupID, 10, 64)
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    sessions, err := h.sessionModel.GetByGroupID(groupIDInt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(sessions)
}
