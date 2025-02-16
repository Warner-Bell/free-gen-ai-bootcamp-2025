package models

import (
    "database/sql"
    "time"
)

type Group struct {
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}

type GroupModel struct {
    db *sql.DB
}

type GroupWord struct {
    WordID    int64     `json:"word_id"`
    GroupID   int64     `json:"group_id"`
    CreatedAt time.Time `json:"created_at"`
    Word      Word      `json:"word"`
}

func NewGroupModel(db *sql.DB) *GroupModel {
    return &GroupModel{db: db}
}

func (m *GroupModel) GetAll() ([]Group, error) {
    rows, err := m.db.Query(`
        SELECT id, name, description, created_at, updated_at 
        FROM groups 
        ORDER BY created_at DESC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var groups []Group
    for rows.Next() {
        var g Group
        err := rows.Scan(
            &g.ID,
            &g.Name,
            &g.Description,
            &g.CreatedAt,
            &g.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        groups = append(groups, g)
    }
    return groups, nil
}

func (m *GroupModel) Create(group *Group) error {
    query := `
        INSERT INTO groups (name, description, created_at, updated_at)
        VALUES (?, ?, datetime('now'), datetime('now'))
        RETURNING id, created_at, updated_at`
    
    return m.db.QueryRow(
        query,
        group.Name,
        group.Description,
    ).Scan(&group.ID, &group.CreatedAt, &group.UpdatedAt)
}

func (m *GroupModel) GetByID(id int64) (*Group, error) {
    var group Group
    err := m.db.QueryRow(`
        SELECT id, name, description, created_at, updated_at 
        FROM groups 
        WHERE id = ?`, id,
    ).Scan(
        &group.ID,
        &group.Name,
        &group.Description,
        &group.CreatedAt,
        &group.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &group, nil
}

func (m *GroupModel) AddWord(groupID, wordID int64) error {
    _, err := m.db.Exec(`
        INSERT INTO word_groups (word_id, group_id)
        VALUES (?, ?)`,
        wordID, groupID)
    return err
}

func (m *GroupModel) RemoveWord(groupID, wordID int64) error {
    _, err := m.db.Exec(`
        DELETE FROM word_groups
        WHERE group_id = ? AND word_id = ?`,
        groupID, wordID)
    return err
}

func (m *GroupModel) GetGroupWords(groupID int64) ([]GroupWord, error) {
    rows, err := m.db.Query(`
        SELECT 
            wg.word_id,
            wg.group_id,
            wg.created_at,
            w.word,
            w.translation,
            w.notes,
            w.created_at,
            w.updated_at
        FROM word_groups wg
        JOIN words w ON w.id = wg.word_id
        WHERE wg.group_id = ?
        ORDER BY wg.created_at DESC`,
        groupID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var groupWords []GroupWord
    for rows.Next() {
        var gw GroupWord
        var w Word
        err := rows.Scan(
            &gw.WordID,
            &gw.GroupID,
            &gw.CreatedAt,
            &w.Word,
            &w.Translation,
            &w.Notes,
            &w.CreatedAt,
            &w.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        w.ID = gw.WordID
        gw.Word = w
        groupWords = append(groupWords, gw)
    }
    return groupWords, nil
}
