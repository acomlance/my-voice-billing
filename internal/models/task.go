package models

import "time"

type Task struct {
	Id         int64      `db:"id"`
	TokenId    int64      `db:"token_id"`
	Task       string     `db:"task"`
	DateCreate *time.Time `db:"date_create"`
	DateUpdate *time.Time `db:"date_update"`
}
