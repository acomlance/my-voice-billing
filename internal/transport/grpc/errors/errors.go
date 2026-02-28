package errors

import (
	"errors"

	"my-voice-billing/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func IsNotFound(err error) bool {
	return errors.Is(err, domain.ErrNotFound)
}

func IsConflict(err error) bool {
	return errors.Is(err, domain.ErrConflict)
}

func IsInvalid(err error) bool {
	return errors.Is(err, domain.ErrInvalid)
}

func IsInsufficientBalance(err error) bool {
	return errors.Is(err, domain.ErrInsufficientBalance)
}

func ToStatus(err error) *status.Status {
	if err == nil {
		return nil
	}
	if IsNotFound(err) {
		return status.New(codes.NotFound, "")
	}
	if IsConflict(err) {
		return status.New(codes.AlreadyExists, "")
	}
	if IsInvalid(err) {
		return status.New(codes.InvalidArgument, "")
	}
	if IsInsufficientBalance(err) {
		return status.New(codes.FailedPrecondition, "")
	}
	return status.New(codes.Internal, "")
}
