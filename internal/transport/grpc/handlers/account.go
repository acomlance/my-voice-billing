package handlers

import (
	"context"

	"my-voice-billing/internal/container"
	"my-voice-billing/internal/models"
	"my-voice-billing/internal/transport/grpc/errors"
	"my-voice-billing/internal/transport/grpc/pb"

	"github.com/rs/zerolog/log"
)

type AccountServer struct {
	pb.UnimplementedAccountServiceServer
	c *container.Container
}

func NewAccountServer(c *container.Container) *AccountServer {
	return &AccountServer{c: c}
}

func (s *AccountServer) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	a, err := s.c.AccountLogic.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, handleErr(err, "GetAccount")
	}
	return accountToProto(a), nil
}

func (s *AccountServer) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	a := &models.Account{
		Id:      req.GetId(),
		Balance: req.GetBalance(),
		Reserve: req.GetReserve(),
		State:   int16(req.GetState()),
	}
	if err := s.c.AccountLogic.Create(ctx, a); err != nil {
		return nil, handleErr(err, "CreateAccount")
	}
	return &pb.CreateAccountResponse{Id: a.Id}, nil
}

func (s *AccountServer) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountResponse, error) {
	a, err := s.c.AccountLogic.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, handleErr(err, "UpdateAccount")
	}
	a.Balance = req.GetBalance()
	a.Reserve = req.GetReserve()
	a.State = int16(req.GetState())
	if err := s.c.AccountLogic.Update(ctx, a); err != nil {
		return nil, handleErr(err, "UpdateAccount")
	}
	return &pb.UpdateAccountResponse{}, nil
}

func handleErr(err error, method string) error {
	if !errors.IsNotFound(err) && !errors.IsConflict(err) {
		log.Error().Err(err).Str("method", method).Msg("internal error")
	}
	return errors.ToStatus(err).Err()
}

func accountToProto(a *models.Account) *pb.GetAccountResponse {
	return &pb.GetAccountResponse{
		Id:      a.Id,
		Balance: a.Balance,
		Reserve: a.Reserve,
		State:   int32(a.State),
	}
}
