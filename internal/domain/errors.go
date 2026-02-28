package domain

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrConflict            = errors.New("conflict")
	ErrInvalid             = errors.New("invalid argument")
	ErrInsufficientBalance  = errors.New("insufficient balance")
)
