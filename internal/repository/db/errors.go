package db

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ErrDuplicateEntry    = "duplicate key"
	ErrForeignKey        = "foreign key"
	PgCodeCheckViolation = "23514"
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

func IsCheckConstraintError(err error) bool {
	var pgErr *pgconn.PgError
	if err != nil && errors.As(err, &pgErr) {
		return pgErr.Code == PgCodeCheckViolation
	}
	return false
}
