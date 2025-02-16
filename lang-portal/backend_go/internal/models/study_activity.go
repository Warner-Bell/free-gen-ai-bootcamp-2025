package models

import (
    "database/sql"
    "encoding/json"
    "time"
)

type StudyActivity struct {
    ID          int64           `json:"id"`
    Name        string          `json:"name"`
    Description string          `json:"description"`
    Type        string          `json:"type"`
    Settings    json.RawMessage `json:"settings,omitempty"`
    Active      bool            `json:"active"`
    CreatedAt   time.Time       `json:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at"`
}

type StudyActivityModel struct {
    DB *sql.DB
}

func NewStudyActivityModel(db *sql.DB) *StudyActivityModel {
    return &StudyActivityModel{DB: db}
}

func (m *StudyActivityModel) Create(name, description, activityType string, settings json.RawMessage) (*StudyActivity, error) {
    activity := &StudyActivity{
        Name:        name,
        Description: description,
        Type:        activityType,
        Settings:    settings,
        Active:      true,
    }

    query := `
        INSERT INTO study_activities (name, description, type, settings, active, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, DATETIME('now'), DATETIME('now'))
        RETURNING id, created_at, updated_at`

    err := m.DB.QueryRow(
        query,
        activity.Name,
        activity.Description,
        activity.Type,
        activity.Settings,
        activity.Active,
    ).Scan(&activity.ID, &activity.CreatedAt, &activity.UpdatedAt)

    if err != nil {
        return nil, err
    }

    return activity, nil
}

func (m *StudyActivityModel) GetByID(id int64) (*StudyActivity, error) {
    activity := &StudyActivity{}
    err := m.DB.QueryRow(`
        SELECT id, name, description, type, settings, active, created_at, updated_at
        FROM study_activities
        WHERE id = ?`,
        id,
    ).Scan(
        &activity.ID,
        &activity.Name,
        &activity.Description,
        &activity.Type,
        &activity.Settings,
        &activity.Active,
        &activity.CreatedAt,
        &activity.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return activity, nil
}

func (m *StudyActivityModel) Update(id int64, name, description, activityType string, settings json.RawMessage, active bool) (*StudyActivity, error) {
    query := `
        UPDATE study_activities
        SET name = ?, description = ?, type = ?, settings = ?, active = ?, updated_at = DATETIME('now')
        WHERE id = ?`

    result, err := m.DB.Exec(
        query,
        name,
        description,
        activityType,
        settings,
        active,
        id,
    )
    if err != nil {
        return nil, err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return nil, err
    }
    if rowsAffected == 0 {
        return nil, sql.ErrNoRows
    }

    return m.GetByID(id)
}

func (m *StudyActivityModel) Delete(id int64) error {
    result, err := m.DB.Exec("DELETE FROM study_activities WHERE id = ?", id)
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
    return nil
}

func (m *StudyActivityModel) GetAll() ([]StudyActivity, error) {
    rows, err := m.DB.Query(`
        SELECT id, name, description, type, settings, active, created_at, updated_at
        FROM study_activities
        ORDER BY created_at DESC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var activities []StudyActivity
    for rows.Next() {
        var a StudyActivity
        err := rows.Scan(
            &a.ID,
            &a.Name,
            &a.Description,
            &a.Type,
            &a.Settings,
            &a.Active,
            &a.CreatedAt,
            &a.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        activities = append(activities, a)
    }
    return activities, nil
}

func (m *StudyActivityModel) GetStats(id int64) (map[string]interface{}, error) {
    var stats = make(map[string]interface{})

    // Get total sessions using this activity
    var totalSessions int
    err := m.DB.QueryRow(`
        SELECT COUNT(*)
        FROM study_sessions
        WHERE activity_name = (
            SELECT name
            FROM study_activities
            WHERE id = ?
        )`,
        id,
    ).Scan(&totalSessions)
    if err != nil {
        return nil, err
    }

    // Get average completion rate
    var avgCompletionRate sql.NullFloat64
    err = m.DB.QueryRow(`
        SELECT AVG(CASE 
            WHEN end_time IS NOT NULL THEN 1
            ELSE 0
        END) * 100
        FROM study_sessions
        WHERE activity_name = (
            SELECT name
            FROM study_activities
            WHERE id = ?
        )`,
        id,
    ).Scan(&avgCompletionRate)
    if err != nil {
        return nil, err
    }

    stats["total_sessions"] = totalSessions
    if avgCompletionRate.Valid {
        stats["avg_completion_rate"] = avgCompletionRate.Float64
    } else {
        stats["avg_completion_rate"] = 0.0
    }

    return stats, nil
}
