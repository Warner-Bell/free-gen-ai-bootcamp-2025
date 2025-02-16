package models

import (
    "database/sql"
    "time"
)

type StudySession struct {
    ID           int64      `json:"id"`
    GroupID      int64      `json:"group_id"`
    ActivityName string     `json:"activity_name"`
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
    DB *sql.DB
}

func NewStudySessionModel(db *sql.DB) *StudySessionModel {
    return &StudySessionModel{DB: db}
}

func (m *StudySessionModel) Create(groupID int64, activityName string) (*StudySession, error) {
    result, err := m.DB.Exec(`
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

    err := m.DB.QueryRow(`
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

func (m *StudySessionModel) GetAll() ([]StudySession, error) {
    rows, err := m.DB.Query(`
        SELECT id, group_id, activity_name, start_time, end_time
        FROM study_sessions
        ORDER BY start_time DESC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var sessions []StudySession
    for rows.Next() {
        var s StudySession
        var endTime sql.NullTime
        err := rows.Scan(&s.ID, &s.GroupID, &s.ActivityName, &s.StartTime, &endTime)
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

func (m *StudySessionModel) GetByGroupID(groupID int64) ([]StudySession, error) {
    rows, err := m.DB.Query(`
        SELECT id, group_id, activity_name, start_time, end_time
        FROM study_sessions
        WHERE group_id = ?
        ORDER BY start_time DESC`, groupID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var sessions []StudySession
    for rows.Next() {
        var s StudySession
        var endTime sql.NullTime
        err := rows.Scan(&s.ID, &s.GroupID, &s.ActivityName, &s.StartTime, &endTime)
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

func (m *StudySessionModel) End(id int64) error {
    _, err := m.DB.Exec(`
        UPDATE study_sessions
        SET end_time = CURRENT_TIMESTAMP
        WHERE id = ? AND end_time IS NULL`,
        id)
    return err
}

func (m *StudySessionModel) AddReview(sessionID, wordID int64, correct bool) error {
    _, err := m.DB.Exec(`
        INSERT INTO word_reviews (word_id, study_session_id, correct, created_at)
        VALUES (?, ?, ?, CURRENT_TIMESTAMP)`,
        wordID, sessionID, correct)
    return err
}

func (m *StudySessionModel) GetSessionStats(sessionID int64) (map[string]interface{}, error) {
    var total, correct int
    err := m.DB.QueryRow(`
        SELECT 
            COUNT(*) as total,
            SUM(CASE WHEN correct = 1 THEN 1 ELSE 0 END) as correct
        FROM word_reviews
        WHERE study_session_id = ?`,
        sessionID).Scan(&total, &correct)
    if err != nil {
        return nil, err
    }

    var accuracy float64
    if total > 0 {
        accuracy = float64(correct) / float64(total) * 100
    }

    return map[string]interface{}{
        "total_reviews": total,
        "correct":       correct,
        "accuracy":      accuracy,
    }, nil
}

func (m *StudySessionModel) GetSessionReviews(sessionID int64) ([]WordReview, error) {
    rows, err := m.DB.Query(`
        SELECT id, word_id, study_session_id, correct, created_at
        FROM word_reviews
        WHERE study_session_id = ?
        ORDER BY created_at`,
        sessionID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var reviews []WordReview
    for rows.Next() {
        var r WordReview
        err := rows.Scan(&r.ID, &r.WordID, &r.StudySessionID, &r.Correct, &r.CreatedAt)
        if err != nil {
            return nil, err
        }
        reviews = append(reviews, r)
    }
    return reviews, nil
}

func (m *StudySessionModel) Delete(id int64) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Delete associated reviews first
    _, err = tx.Exec("DELETE FROM word_reviews WHERE study_session_id = ?", id)
    if err != nil {
        return err
    }

    // Delete the session
    result, err := tx.Exec("DELETE FROM study_sessions WHERE id = ?", id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }

    return tx.Commit()
}
