package models

import "time"

type Task struct {
	Id             int64      `db:"id"`
	Token          string     `db:"token"`
	AccountId      int64      `db:"account_id"`
	ReservedTokens int64      `db:"reserved_tokens"`
	DateCreate     *time.Time `db:"date_create"`
}
