package models

import "time"

type Token struct {
	Id         int64      `db:"id"`
	UserId     int64      `db:"user_id"`
	Token      string     `db:"token"`
	DateCreate *time.Time `db:"date_create"`
	DateUpdate *time.Time `db:"date_update"`
}
