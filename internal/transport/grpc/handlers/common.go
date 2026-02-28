package handlers

import (
	"my-voice-billing/internal/transport/grpc/errors"

	"github.com/rs/zerolog/log"
)

func handleErr(err error, method string) error {
	if !errors.IsNotFound(err) && !errors.IsConflict(err) && !errors.IsInvalid(err) && !errors.IsInsufficientBalance(err) {
		log.Error().Err(err).Str("method", method).Msg("internal error")
	}
	return errors.ToStatus(err).Err()
}
