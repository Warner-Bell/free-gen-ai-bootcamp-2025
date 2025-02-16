// cmd/server/main.go
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

    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/models"
    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/handlers/words"
    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/handlers/groups"
    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/handlers/sessions"
    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/handlers/dashboard"
    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/handlers/settings"
    "free-gen-ai-bootcamp-2025/lang-portal/backend_go/internal/handlers/study_activities"
)

func main() {
    // Initialize database
    db, err := sql.Open("sqlite3", "./words.db?_foreign_keys=on&_journal_mode=WAL")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Test database connection
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

    // Initialize models
    wordModel := models.NewWordModel(db)
    groupModel := models.NewGroupModel(db)
    sessionModel := models.NewStudySessionModel(db)
    dashboardModel := models.NewDashboardModel(db)

    // Initialize handlers
    wordHandler := words.NewHandler(wordModel)
    groupHandler := groups.NewHandler(groupModel)
    sessionHandler := sessions.NewHandler(sessionModel)
    dashboardHandler := dashboard.NewHandler(dashboardModel)

    // Initialize router
    r := chi.NewRouter()

    // Middleware
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.RequestID)
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    }))

    // Routes
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    // API routes
    r.Route("/api", func(r chi.Router) {
        // Words endpoints
        r.Get("/words", wordHandler.GetWords)
        r.Post("/words", wordHandler.CreateWord)
        r.Get("/words/{id}", wordHandler.GetWord)
        r.Put("/words/{id}", wordHandler.UpdateWord)
        r.Delete("/words/{id}", wordHandler.DeleteWord)

        // Groups endpoints
        r.Get("/groups", groupHandler.GetGroups)
        r.Post("/groups", groupHandler.CreateGroup)
        r.Get("/groups/{id}", groupHandler.GetGroup)
        r.Put("/groups/{id}", groupHandler.UpdateGroup)
        r.Delete("/groups/{id}", groupHandler.DeleteGroup)
        r.Post("/groups/{id}/words", groupHandler.AddWordToGroup)
        r.Delete("/groups/{id}/words/{wordId}", groupHandler.RemoveWordFromGroup)
        r.Get("/groups/{id}/words", groupHandler.GetGroupWords)

        // Study session endpoints
        r.Post("/sessions", sessionHandler.CreateSession)
        r.Put("/sessions/{id}/end", sessionHandler.EndSession)
        r.Post("/sessions/{id}/reviews", sessionHandler.AddReview)
        r.Get("/sessions/{id}/stats", sessionHandler.GetSessionStats)
        r.Get("/sessions", sessionHandler.GetSessions)
        r.Get("/sessions/{id}", sessionHandler.GetSession)
        r.Get("/sessions/group/{groupId}", sessionHandler.GetSessionsByGroup)

        // Dashboard endpoints
        r.Get("/dashboard/stats", dashboardHandler.GetDashboardStats)
        r.Get("/dashboard/recent-sessions", dashboardHandler.GetRecentSessions)
        r.Get("/dashboard/progress", dashboardHandler.GetLearningProgress)
    })

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
