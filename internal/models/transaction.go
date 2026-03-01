package models

import "time"

type Transaction struct {
	Id                     int64      `db:"id"`
	AccountId              int64      `db:"account_id"`
	Status                 int16      `db:"status"`
	Amount                 int64      `db:"amount"`
	Tokens                 int64      `db:"tokens"`
	TokensAfter            *int64     `db:"tokens_after"`
	PaymentType            int16      `db:"payment_type"`
	PaymentMethod          int16      `db:"payment_method"`
	PaymentData            string     `db:"payment_data"`
	Description            string     `db:"description"`
	DateCreate             *time.Time `db:"date_create"`
	DateUpdate             *time.Time `db:"date_update"`
}
