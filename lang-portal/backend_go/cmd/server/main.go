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
