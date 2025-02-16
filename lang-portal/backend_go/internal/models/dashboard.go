package models

import "database/sql"

type DashboardModel struct {
    DB           *sql.DB
    WordModel    *WordModel
    GroupModel   *GroupModel
    SessionModel *StudySessionModel
}

func NewDashboardModel(db *sql.DB, wm *WordModel, gm *GroupModel, sm *StudySessionModel) *DashboardModel {
    return &DashboardModel{
        DB:           db,
        WordModel:    wm,
        GroupModel:   gm,
        SessionModel: sm,
    }
}
