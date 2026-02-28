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

// ToStatus маппит domain-ошибки в gRPC status.
func ToStatus(err error) *status.Status {
	if err == nil {
		return nil
	}
	if IsNotFound(err) {
		return status.New(codes.NotFound, err.Error())
	}
	if IsConflict(err) {
		return status.New(codes.AlreadyExists, err.Error())
	}
	return status.New(codes.Internal, err.Error())
}
