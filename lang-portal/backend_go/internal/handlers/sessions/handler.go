// Package sessions handles study session-related HTTP endpoints
package sessions

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"

    "github.com/go-chi/chi"
    "backend_go/internal/models"
)

type Handler struct {
    sessionModel *models.StudySessionModel
}

func NewHandler(sessionModel *models.StudySessionModel) *Handler {
    return &Handler{sessionModel: sessionModel}
}

type createSessionRequest struct {
    GroupID     int       `json:"groupId"`
    StartedAt   time.Time `json:"startedAt"`
    TargetWords []int     `json:"targetWords"`
}

type addReviewRequest struct {
    WordID int  `json:"wordId"`
    Known  bool `json:"known"`
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
    var req createSessionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    session, err := h.sessionModel.CreateSession(req.GroupID, req.StartedAt, req.TargetWords)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(session)
}

func (h *Handler) EndSession(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    sessionID, err := strconv.Atoi(id)
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }

    if err := h.sessionModel.EndSession(sessionID, time.Now()); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) AddReview(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    sessionID, err := strconv.Atoi(id)
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
    sessionID, err := strconv.Atoi(id)
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
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    sessions, err := h.sessionModel.GetSessions(page)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(sessions)
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    sessionID, err := strconv.Atoi(id)
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }

    session, err := h.sessionModel.GetSessionByID(sessionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(session)
}

func (h *Handler) GetSessionsByGroup(w http.ResponseWriter, r *http.Request) {
    groupID, err := strconv.Atoi(chi.URLParam(r, "groupId"))
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    sessions, err := h.sessionModel.GetSessionsByGroup(groupID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(sessions)
}