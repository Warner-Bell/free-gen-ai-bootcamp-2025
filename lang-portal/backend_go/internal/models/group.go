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
    DB *sql.DB
}

type GroupWord struct {
    WordID    int64     `json:"word_id"`
    GroupID   int64     `json:"group_id"`
    CreatedAt time.Time `json:"created_at"`
    Word      Word      `json:"word"`
}

func NewGroupModel(DB *sql.DB) *GroupModel {
    return &GroupModel{DB: DB} // Fixed: Changed db to DB
}

func (m *GroupModel) GetAll() ([]Group, error) {
    rows, err := m.DB.Query(` // Fixed: Changed m.db to m.DB
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
    
    return m.DB.QueryRow( // Fixed: Changed m.db to m.DB
        query,
        group.Name,
        group.Description,
    ).Scan(&group.ID, &group.CreatedAt, &group.UpdatedAt)
}

func (m *GroupModel) Update(group *Group) error {
    query := `
        UPDATE groups
        SET name = ?, description = ?
        WHERE id = ?`
    
    result, err := m.DB.Exec(query, group.Name, group.Description, group.ID) // Fixed: Changed m.db to m.DB
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

func (m *GroupModel) Delete(id int64) error {
    query := "DELETE FROM groups WHERE id = ?"
    result, err := m.DB.Exec(query, id) // Fixed: Changed m.db to m.DB
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

func (m *GroupModel) GetByID(id int64) (*Group, error) {
    var group Group
    err := m.DB.QueryRow(` // Fixed: Changed m.db to m.DB
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
    _, err := m.DB.Exec(` // Fixed: Changed m.db to m.DB
        INSERT INTO word_groups (word_id, group_id)
        VALUES (?, ?)`,
        wordID, groupID)
    return err
}

func (m *GroupModel) RemoveWord(groupID, wordID int64) error {
    _, err := m.DB.Exec(` // Fixed: Changed m.db to m.DB
        DELETE FROM word_groups
        WHERE group_id = ? AND word_id = ?`,
        groupID, wordID)
    return err
}

func (m *GroupModel) GetGroupWords(groupID int64) ([]GroupWord, error) {
    rows, err := m.DB.Query(` // Fixed: Changed m.db to m.DB
        SELECT 
            w.id, w.word, w.translation, w.notes, w.japanese, w.romaji, w.english, w.created_at, w.updated_at,
            wg.created_at as group_created_at
        FROM words w
        JOIN word_groups wg ON w.id = wg.word_id
        WHERE wg.group_id = ?`, groupID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var groupWords []GroupWord
    for rows.Next() {
        var gw GroupWord
        err := rows.Scan(
            &gw.Word.ID,
            &gw.Word.Word,
            &gw.Word.Translation,
            &gw.Word.Notes,
            &gw.Word.Japanese,
            &gw.Word.Romaji,
            &gw.Word.English,
            &gw.Word.CreatedAt,
            &gw.Word.UpdatedAt,
            &gw.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        gw.WordID = gw.Word.ID
        gw.GroupID = groupID
        groupWords = append(groupWords, gw)
    }
    return groupWords, nil
}
