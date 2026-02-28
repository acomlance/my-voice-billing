package handlers

import (
	"context"

	"my-voice-billing/internal/container"
	"my-voice-billing/internal/models"
	"my-voice-billing/internal/transport/grpc/errors"
	"my-voice-billing/internal/transport/grpc/pb"

	"github.com/rs/zerolog/log"
)

type TokenServer struct {
	pb.UnimplementedTokenServiceServer
	c *container.Container
}

func NewTokenServer(c *container.Container) *TokenServer {
	return &TokenServer{c: c}
}

func (s *TokenServer) GetToken(ctx context.Context, req *pb.GetTokenRequest) (*pb.GetTokenResponse, error) {
	t, err := s.c.TokenLogic.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, handleErr(err, "GetToken")
	}
	return &pb.GetTokenResponse{
		Id:     t.Id,
		UserId: t.UserId,
		Token:  t.Token,
	}, nil
}

func (s *TokenServer) CreateToken(ctx context.Context, req *pb.CreateTokenRequest) (*pb.CreateTokenResponse, error) {
	t := &models.Token{UserId: req.GetUserId(), Token: req.GetToken()}
	if err := s.c.TokenLogic.Create(ctx, t); err != nil {
		return nil, handleErr(err, "CreateToken")
	}
	return &pb.CreateTokenResponse{Id: t.Id}, nil
}

func handleErr(err error, method string) error {
	if !errors.IsNotFound(err) && !errors.IsConflict(err) {
		log.Error().Err(err).Str("method", method).Msg("internal error")
	}
	return errors.ToStatus(err).Err()
}
