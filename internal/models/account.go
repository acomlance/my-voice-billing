package models

import "time"

type Account struct {
	Id         int64      `db:"id"`
	Balance    int64      `db:"balance"`
	Reserve    int64      `db:"reserve"`
	State      int16      `db:"state"`
	DateCreate *time.Time `db:"date_create"`
	DateUpdate *time.Time `db:"date_update"`
}
