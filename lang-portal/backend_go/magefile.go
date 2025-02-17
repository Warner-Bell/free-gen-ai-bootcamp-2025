//go:build mage
// +build mage

package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"

    "github.com/magefile/mage/mg"
    "github.com/magefile/mage/sh"
    _ "github.com/mattn/go-sqlite3"
    "database/sql"
)

const (
    binName    = "backend_go"
    testPkgs   = "./..."
    buildDir   = "build"
    dbFile     = "words.db"
    migrateDir = "migrations"
)

// Build builds the application
func Build() error {
    // Ensure we're in module mode
    os.Setenv("GO111MODULE", "on")
    mg.Deps(Clean)
    fmt.Println("Building...")

    if err := os.MkdirAll(buildDir, 0755); err != nil {
        return err
    }

    // Build using module mode
    cmd := exec.Command("go", "build", "-mod=mod", "-o", filepath.Join(buildDir, binName), "./cmd/server")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

// Clean removes build artifacts
func Clean() error {
    fmt.Println("Cleaning...")
    return os.RemoveAll(buildDir)
}

// Test runs the test suite
func Test() error {
    fmt.Println("Running tests...")
    return sh.Run("go", "test", "-v", "-race", "-cover", testPkgs)
}

// Lint runs the linter
func Lint() error {
    fmt.Println("Running linter...")
    return sh.Run("golangci-lint", "run")
}

// Dev runs the application in development mode
func Dev() error {
    mg.Deps(InitDB)
    fmt.Println("Starting development server...")
    return sh.Run("go", "run", "./cmd/server")
}

// InitDB initializes the database with schema
func InitDB() error {
    fmt.Println("Initializing database...")

    db, err := sql.Open("sqlite3", dbFile)
    if err != nil {
        return fmt.Errorf("error opening database: %v", err)
    }
    defer db.Close()

    // Create the database schema
    schema := `
    PRAGMA foreign_keys = ON;

    CREATE TABLE IF NOT EXISTS words (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        word TEXT NOT NULL,
        translation TEXT,
        notes TEXT,
        japanese TEXT,
        romaji TEXT,
        english TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS groups (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE,
        description TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS word_groups (
        word_id INTEGER,
        group_id INTEGER,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (word_id, group_id),
        FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
        FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
    );

    CREATE TABLE IF NOT EXISTS study_sessions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        activity_name TEXT NOT NULL,
        group_id INTEGER,
        start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        end_time TIMESTAMP,
        FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE SET NULL
    );

    CREATE TABLE IF NOT EXISTS word_reviews (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        word_id INTEGER,
        study_session_id INTEGER,
        correct BOOLEAN NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
        FOREIGN KEY (study_session_id) REFERENCES study_sessions(id) ON DELETE CASCADE
    );

    -- Indexes for better performance
    CREATE INDEX IF NOT EXISTS idx_word_groups_group_id ON word_groups(group_id);
    CREATE INDEX IF NOT EXISTS idx_word_groups_word_id ON word_groups(word_id);
    CREATE INDEX IF NOT EXISTS idx_word_reviews_session_id ON word_reviews(study_session_id);
    CREATE INDEX IF NOT EXISTS idx_word_reviews_word_id ON word_reviews(word_id);
    CREATE INDEX IF NOT EXISTS idx_study_sessions_group_id ON study_sessions(group_id);
    `

    _, err = db.Exec(schema)
    if err != nil {
        return fmt.Errorf("error creating schema: %v", err)
    }

    return nil
}
