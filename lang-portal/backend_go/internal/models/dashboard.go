// internal/models/dashboard.go
package models

import (
    "database/sql"
    "time"
)

type DashboardModel struct {
    db *sql.DB
}

type DashboardStats struct {
    TotalWords         int     `json:"total_words"`
    WordsStudied       int     `json:"words_studied"`
    TotalSessions     int     `json:"total_sessions"`
    AverageAccuracy   float64 `json:"average_accuracy"`
    StudyStreakDays   int     `json:"study_streak_days"`
    ActiveGroups      int     `json:"active_groups"`
}

type RecentSession struct {
    ID           int64      `json:"id"`
    ActivityName string     `json:"activity_name"`
    GroupName    string     `json:"group_name"`
    StartTime    time.Time  `json:"start_time"`
    EndTime      *time.Time `json:"end_time"`
    Accuracy     float64    `json:"accuracy"`
}

func NewDashboardModel(db *sql.DB) *DashboardModel {
    return &DashboardModel{db: db}
}

func (m *DashboardModel) GetStats() (*DashboardStats, error) {
    stats := &DashboardStats{}

    // Get total words and studied words
    err := m.db.QueryRow(`
        SELECT 
            (SELECT COUNT(*) FROM words) as total_words,
            COUNT(DISTINCT word_id) as words_studied
        FROM word_reviews`).Scan(&stats.TotalWords, &stats.WordsStudied)
    if err != nil {
        return nil, err
    }

    // Get session stats
    err = m.db.QueryRow(`
        SELECT 
            COUNT(*) as total_sessions,
            COALESCE(AVG(CASE WHEN correct THEN 100.0 ELSE 0.0 END), 0) as avg_accuracy
        FROM study_sessions s
        LEFT JOIN word_reviews wr ON wr.study_session_id = s.id`).
        Scan(&stats.TotalSessions, &stats.AverageAccuracy)
    if err != nil {
        return nil, err
    }

    // Get study streak
    err = m.db.QueryRow(`
        WITH RECURSIVE dates AS (
            SELECT date(MAX(start_time)) as date
            FROM study_sessions
            UNION ALL
            SELECT date(date, '-1 day')
            FROM dates
            WHERE EXISTS (
                SELECT 1 
                FROM study_sessions 
                WHERE date(start_time) = date(dates.date, '-1 day')
            )
        )
        SELECT COUNT(*) FROM dates`).Scan(&stats.StudyStreakDays)
    if err != nil {
        return nil, err
    }

    // Get active groups
    err = m.db.QueryRow(`
        SELECT COUNT(DISTINCT group_id) 
        FROM study_sessions 
        WHERE start_time >= datetime('now', '-30 days')`).
        Scan(&stats.ActiveGroups)
    if err != nil {
        return nil, err
    }

    return stats, nil
}

func (m *DashboardModel) GetRecentSessions(limit int) ([]RecentSession, error) {
    rows, err := m.db.Query(`
        SELECT 
            s.id,
            s.activity_name,
            g.name as group_name,
            s.start_time,
            s.end_time,
            COALESCE(AVG(CASE WHEN wr.correct THEN 100.0 ELSE 0.0 END), 0) as accuracy
        FROM study_sessions s
        LEFT JOIN groups g ON g.id = s.group_id
        LEFT JOIN word_reviews wr ON wr.study_session_id = s.id
        GROUP BY s.id
        ORDER BY s.start_time DESC
        LIMIT ?`,
        limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var sessions []RecentSession
    for rows.Next() {
        var s RecentSession
        var endTime sql.NullTime
        err := rows.Scan(
            &s.ID,
            &s.ActivityName,
            &s.GroupName,
            &s.StartTime,
            &endTime,
            &s.Accuracy,
        )
        if err != nil {
            return nil, err
        }
        if endTime.Valid {
            s.EndTime = &endTime.Time
        }
        sessions = append(sessions, s)
    }
    return sessions, nil
}

func (m *DashboardModel) GetLearningProgress() (map[string]interface{}, error) {
    var totalWords, studiedWords int
    var avgAccuracy float64

    err := m.db.QueryRow(`
        SELECT 
            (SELECT COUNT(*) FROM words) as total_words,
            COUNT(DISTINCT word_id) as studied_words,
            COALESCE(AVG(CASE WHEN correct THEN 100.0 ELSE 0.0 END), 0) as avg_accuracy
        FROM word_reviews`).
        Scan(&totalWords, &studiedWords, &avgAccuracy)
    if err != nil {
        return nil, err
    }

    return map[string]interface{}{
        "total_words":      totalWords,
        "studied_words":    studiedWords,
        "progress_percent": float64(studiedWords) / float64(totalWords) * 100,
        "avg_accuracy":     avgAccuracy,
    }, nil
}
