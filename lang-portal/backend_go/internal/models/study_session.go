// internal/models/study_session.go
package models

import (
    "database/sql"
    "time"
)

type StudySession struct {
    ID           int64      `json:"id"`
    ActivityName string     `json:"activity_name"`
    GroupID      int64      `json:"group_id"`
    StartTime    time.Time  `json:"start_time"`
    EndTime      *time.Time `json:"end_time,omitempty"`
}

type WordReview struct {
    ID             int64     `json:"id"`
    WordID         int64     `json:"word_id"`
    StudySessionID int64     `json:"study_session_id"`
    Correct        bool      `json:"correct"`
    CreatedAt      time.Time `json:"created_at"`
}

type StudySessionModel struct {
    db *sql.DB
}

func NewStudySessionModel(db *sql.DB) *StudySessionModel {
    return &StudySessionModel{db: db}
}

func (m *StudySessionModel) Create(groupID int64, activityName string) (*StudySession, error) {
    result, err := m.db.Exec(`
        INSERT INTO study_sessions (group_id, activity_name, start_time)
        VALUES (?, ?, CURRENT_TIMESTAMP)`,
        groupID, activityName)
    if err != nil {
        return nil, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return nil, err
    }

    return m.GetByID(id)
}

func (m *StudySessionModel) GetByID(id int64) (*StudySession, error) {
    var session StudySession
    var endTime sql.NullTime

    err := m.db.QueryRow(`
        SELECT id, activity_name, group_id, start_time, end_time
        FROM study_sessions
        WHERE id = ?`,
        id).Scan(&session.ID, &session.ActivityName, &session.GroupID, &session.StartTime, &endTime)
    if err != nil {
        return nil, err
    }

    if endTime.Valid {
        session.EndTime = &endTime.Time
    }

    return &session, nil
}

func (m *StudySessionModel) End(id int64) error {
    _, err := m.db.Exec(`
        UPDATE study_sessions
        SET end_time = CURRENT_TIMESTAMP
        WHERE id = ? AND end_time IS NULL`,
        id)
    return err
}

func (m *StudySessionModel) AddReview(sessionID, wordID int64, correct bool) error {
    _, err := m.db.Exec(`
        INSERT INTO word_reviews (word_id, study_session_id, correct)
        VALUES (?, ?, ?)`,
        wordID, sessionID, correct)
    return err
}

func (m *StudySessionModel) GetSessionStats(sessionID int64) (map[string]interface{}, error) {
    var total, correct int
    err := m.db.QueryRow(`
        SELECT 
            COUNT(*) as total,
            SUM(CASE WHEN correct = 1 THEN 1 ELSE 0 END) as correct
        FROM word_reviews
        WHERE study_session_id = ?`,
        sessionID).Scan(&total, &correct)
    if err != nil {
        return nil, err
    }

    return map[string]interface{}{
        "total_reviews": total,
        "correct":       correct,
        "accuracy":      float64(correct) / float64(total) * 100,
    }, nil
}
