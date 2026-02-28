package db

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
)

const (
	ErrDuplicateEntry = "duplicate key"
	ErrForeignKey     = "foreign key"
)

func IsNoRowsError(err error) bool {
	return err != nil && (errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows))
}

func IsDuplicateError(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), ErrDuplicateEntry)
}

func IsForeignKeyError(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), ErrForeignKey)
}
