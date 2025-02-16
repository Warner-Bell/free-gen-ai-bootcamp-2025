#!/bin/bash

# Create directories if they don't exist
mkdir -p internal/models
mkdir -p internal/handlers/dashboard
mkdir -p internal/handlers/settings
mkdir -p internal/handlers/words
mkdir -p internal/handlers/groups
mkdir -p internal/handlers/sessions
mkdir -p internal/handlers/study_activities
mkdir -p cmd/server

# Create dashboard.go
cat > internal/models/dashboard.go << 'EOL'
package models

import "database/sql"

type DashboardModel struct {
    DB           *sql.DB
    WordModel    *WordModel
    GroupModel   *GroupModel
    SessionModel *StudySessionModel
}

func NewDashboardModel(db *sql.DB, wm *WordModel, gm *GroupModel, sm *StudySessionModel) *DashboardModel {
    return &DashboardModel{
        DB:           db,
        WordModel:    wm,
        GroupModel:   gm,
        SessionModel: sm,
    }
}
EOL

# Create dashboard handler
cat > internal/handlers/dashboard/handler.go << 'EOL'
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
EOL

# Create settings handler
cat > internal/handlers/settings/handler.go << 'EOL'
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
EOL

# Update words handler
cat > internal/handlers/words/handler.go << 'EOL'
package words

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/models"
)

type Handler struct {
    model *models.WordModel
}

func NewHandler(model *models.WordModel) *Handler {
    return &Handler{model: model}
}

func (h *Handler) GetWords(w http.ResponseWriter, r *http.Request) {
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

    if limit == 0 {
        limit = 10
    }

    words, total, err := h.model.GetWords(offset, limit)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response := struct {
        Words []models.Word `json:"words"`
        Total int          `json:"total"`
    }{
        Words: words,
        Total: total,
    }

    json.NewEncoder(w).Encode(response)
}

func (h *Handler) SearchWords(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

    if limit == 0 {
        limit = 10
    }

    words, total, err := h.model.SearchWords(query, offset, limit)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response := struct {
        Words []models.Word `json:"words"`
        Total int          `json:"total"`
    }{
        Words: words,
        Total: total,
    }

    json.NewEncoder(w).Encode(response)
}

func (h *Handler) CreateWord(w http.ResponseWriter, r *http.Request) {
    var word models.Word
    if err := json.NewDecoder(r.Body).Decode(&word); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.model.Create(&word); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(word)
}

func (h *Handler) GetWord(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    word, err := h.model.GetByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if word == nil {
        http.NotFound(w, r)
        return
    }

    json.NewEncoder(w).Encode(word)
}

func (h *Handler) UpdateWord(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var word models.Word
    if err := json.NewDecoder(r.Body).Decode(&word); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    word.ID = id
    if err := h.model.Update(&word); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(word)
}

func (h *Handler) DeleteWord(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    if err := h.model.Delete(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
EOL

# Create main.go
cat > cmd/server/main.go << 'EOL'
package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/cors"
    _ "github.com/mattn/go-sqlite3"

    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/handlers/dashboard"
    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/handlers/groups"
    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/handlers/sessions"
    settingsHandler "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/handlers/settings"
    activityHandler "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/handlers/study_activities"
    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/handlers/words"
    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/models"
)

func main() {
    // Initialize database
    db, err := sql.Open("sqlite3", "words.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Initialize models
    wordModel := models.NewWordModel(db)
    groupModel := models.NewGroupModel(db)
    sessionModel := models.NewStudySessionModel(db)
    activityModel := models.NewStudyActivityModel(db)
    dashboardModel := models.NewDashboardModel(db, wordModel, groupModel, sessionModel)

    // Initialize handlers
    wordHandler := words.NewHandler(wordModel)
    groupHandler := groups.NewHandler(groupModel)
    sessionHandler := sessions.NewHandler(sessionModel)
    dashboardHandler := dashboard.NewHandler(dashboardModel)
    settingsHandler := settingsHandler.NewHandler(sessionModel, activityModel)
    activityHandler := activityHandler.NewHandler(activityModel)

    // Initialize router
    r := chi.NewRouter()

    // Middleware
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:          300,
    }))

    // Routes
    r.Route("/api", func(r chi.Router) {
        // Dashboard routes
        r.Get("/dashboard", dashboardHandler.GetDashboard)

        // Word routes
        r.Route("/words", func(r chi.Router) {
            r.Get("/", wordHandler.GetWords)
            r.Post("/", wordHandler.CreateWord)
            r.Get("/search", wordHandler.SearchWords)
            r.Get("/{id}", wordHandler.GetWord)
            r.Put("/{id}", wordHandler.UpdateWord)
            r.Delete("/{id}", wordHandler.DeleteWord)
        })

        // Group routes
        r.Route("/groups", func(r chi.Router) {
            r.Get("/", groupHandler.GetGroups)
            r.Post("/", groupHandler.CreateGroup)
            r.Get("/{id}", groupHandler.GetGroup)
            r.Put("/{id}", groupHandler.UpdateGroup)
            r.Delete("/{id}", groupHandler.DeleteGroup)
            r.Get("/{id}/words", groupHandler.GetGroupWords)
            r.Post("/{id}/words", groupHandler.AddWordToGroup)
            r.Delete("/{id}/words/{wordId}", groupHandler.RemoveWordFromGroup)
        })

        // Study session routes
        r.Route("/sessions", func(r chi.Router) {
            r.Get("/", sessionHandler.GetSessions)
            r.Post("/", sessionHandler.CreateSession)
            r.Get("/{id}", sessionHandler.GetSession)
            r.Put("/{id}/end", sessionHandler.EndSession)
            r.Post("/{id}/reviews", sessionHandler.AddReview)
            r.Get("/{id}/stats", sessionHandler.GetSessionStats)
        })

        // Study activity routes
        r.Route("/activities", func(r chi.Router) {
            r.Get("/", activityHandler.ListActivities)
            r.Post("/", activityHandler.CreateActivity)
            r.Get("/{id}", activityHandler.GetActivity)
            r.Put("/{id}", activityHandler.UpdateActivity)
            r.Delete("/{id}", activityHandler.DeleteActivity)
            r.Get("/{id}/stats", activityHandler.GetActivityStats)
        })

        // Settings routes
        r.Route("/settings", func(r chi.Router) {
            r.Get("/", settingsHandler.GetSettings)
            r.Put("/", settingsHandler.UpdateSettings)
        })
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s...\n", port)
    if err := http.ListenAndServe(":"+port, r); err != nil {
        log.Fatal(err)
    }
}
EOL

echo "Files have been created and updated successfully!"
echo "Now run: go mod tidy && mage build"
