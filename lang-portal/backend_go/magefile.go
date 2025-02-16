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

    // Initialize vendor directory
    if err := exec.Command("bash", "vendor-init.sh").Run(); err != nil {
        return fmt.Errorf("failed to initialize vendor: %v", err)
    }

    // Then build
    cmd := exec.Command("go", "build", "-mod=vendor", "-o", filepath.Join(buildDir, binName), "./cmd/server")
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
    
    // Create migrations directory if it doesn't exist
    if err := os.MkdirAll(migrateDir, 0755); err != nil {
        return err
    }

    // Create the database schema
    schema := `
    PRAGMA foreign_keys = ON;

    CREATE TABLE IF NOT EXISTS words (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        japanese TEXT NOT NULL,
        romaji TEXT NOT NULL,
        english TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS groups (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        description TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS word_groups (
        word_id INTEGER,
        group_id INTEGER,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (word_id, group_id),
        FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
        FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
    );

    CREATE TABLE IF NOT EXISTS study_sessions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        group_id INTEGER,
        activity_name TEXT NOT NULL,
        start_time DATETIME DEFAULT CURRENT_TIMESTAMP,
        end_time DATETIME,
        FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE SET NULL
    );

    CREATE TABLE IF NOT EXISTS word_reviews (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        word_id INTEGER,
        study_session_id INTEGER,
        correct BOOLEAN NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
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

    // Write schema to migration file
    schemaFile := filepath.Join(migrateDir, "001_initial_schema.sql")
    if err := os.WriteFile(schemaFile, []byte(schema), 0644); err != nil {
        return err
    }

    // Execute schema
    db, err := exec.Command("sqlite3", dbFile).StdinPipe()
    if err != nil {
        return err
    }

    cmd := exec.Command("sqlite3", dbFile)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Start(); err != nil {
        return err
    }

    if _, err := db.Write([]byte(schema)); err != nil {
        return err
    }

    db.Close()
    return cmd.Wait()
}

// All runs all the main tasks
func All() {
    mg.SerialDeps(Clean, Lint, Test, Build)
}
