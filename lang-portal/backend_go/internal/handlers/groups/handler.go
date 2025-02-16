// internal/models/group.go
package models

import (
    "database/sql"
    "time"
)

type Group struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type GroupModel struct {
    db *sql.DB
}

func NewGroupModel(db *sql.DB) *GroupModel {
    return &GroupModel{db: db}
}

func (m *GroupModel) GetAll() ([]Group, error) {
    query := `SELECT id, name, created_at, updated_at FROM groups ORDER BY id`
    rows, err := m.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var groups []Group
    for rows.Next() {
        var group Group
        err := rows.Scan(
            &group.ID,
            &group.Name,
            &group.CreatedAt,
            &group.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        groups = append(groups, group)
    }
    return groups, nil
}

func (m *GroupModel) Create(name string) (*Group, error) {
    group := &Group{
        Name:      name,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    query := `
        INSERT INTO groups (name, created_at, updated_at)
        VALUES (?, ?, ?)
        RETURNING id`
    
    err := m.db.QueryRow(
        query,
        group.Name,
        group.CreatedAt,
        group.UpdatedAt,
    ).Scan(&group.ID)

    if err != nil {
        return nil, err
    }

    return group, nil
}

func (m *GroupModel) GetByID(id int64) (*Group, error) {
    group := &Group{}
    query := `SELECT id, name, created_at, updated_at FROM groups WHERE id = ?`
    err := m.db.QueryRow(query, id).Scan(
        &group.ID,
        &group.Name,
        &group.CreatedAt,
        &group.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return group, nil
}

func (m *GroupModel) AddWord(groupID, wordID int64) error {
    query := `INSERT INTO word_groups (group_id, word_id) VALUES (?, ?)`
    _, err := m.db.Exec(query, groupID, wordID)
    return err
}

func (m *GroupModel) RemoveWord(groupID, wordID int64) error {
    query := `DELETE FROM word_groups WHERE group_id = ? AND word_id = ?`
    _, err := m.db.Exec(query, groupID, wordID)
    return err
}

func (m *GroupModel) GetGroupWords(groupID int64) ([]Word, error) {
    query := `
        SELECT w.id, w.japanese, w.romaji, w.english, w.created_at, w.updated_at
        FROM words w
        JOIN word_groups wg ON w.id = wg.word_id
        WHERE wg.group_id = ?
        ORDER BY w.id`
    
    rows, err := m.db.Query(query, groupID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var words []Word
    for rows.Next() {
        var word Word
        err := rows.Scan(
            &word.ID,
            &word.Japanese,
            &word.Romaji,
            &word.English,
            &word.CreatedAt,
            &word.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        words = append(words, word)
    }
    return words, nil
}
