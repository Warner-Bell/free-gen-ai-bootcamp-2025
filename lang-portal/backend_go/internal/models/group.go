type GroupWord struct {
    WordID    int64     `json:"word_id"`
    GroupID   int64     `json:"group_id"`
    CreatedAt time.Time `json:"created_at"`
    Word      Word      `json:"word"`
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
            w.japanese,
            w.romaji,
            w.english,
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
            &w.Japanese,
            &w.Romaji,
            &w.English,
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