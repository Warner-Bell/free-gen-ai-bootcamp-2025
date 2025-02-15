// internal/handlers/sessions/handler.go
package sessions

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "lang-portal/internal/models"
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
}

type addReviewRequest struct {
    WordID  int64 `json:"word_id"`
    Correct bool  `json:"correct"`
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

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(session)
}

func (h *Handler) EndSession(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }

    if err := h.sessionModel.End(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) AddReview(w http.ResponseWriter, r *http.Request) {
    sessionIDStr := chi.URLParam(r, "id")
    sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }

    var req addReviewRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.sessionModel.AddReview(sessionID, req.WordID, req.Correct); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetSessionStats(w http.ResponseWriter, r *http.Request) {
    sessionIDStr := chi.URLParam(r, "id")
    sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }

    stats, err := h.sessionModel.GetSessionStats(sessionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}
