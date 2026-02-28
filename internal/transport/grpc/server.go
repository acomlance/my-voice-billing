package grpc

import (
	"context"
	"errors"
	"net"

	"my-voice-billing/internal/container"
	"my-voice-billing/internal/domain"
	"my-voice-billing/internal/models"
	"my-voice-billing/internal/transport/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type tokenServer struct {
	pb.UnimplementedTokenServiceServer
	c *container.Container
}

func (s *tokenServer) GetToken(ctx context.Context, req *pb.GetTokenRequest) (*pb.GetTokenResponse, error) {
	t, err := s.c.TokenLogic.GetByID(ctx, req.GetId())
	if err != nil {
		if isNotFound(err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetTokenResponse{
		Id:     t.Id,
		UserId: t.UserId,
		Token:  t.Token,
	}, nil
}

func (s *tokenServer) CreateToken(ctx context.Context, req *pb.CreateTokenRequest) (*pb.CreateTokenResponse, error) {
	t := &models.Token{UserId: req.GetUserId(), Token: req.GetToken()}
	if err := s.c.TokenLogic.Create(ctx, t); err != nil {
		if isConflict(err) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CreateTokenResponse{Id: t.Id}, nil
}

type taskServer struct {
	pb.UnimplementedTaskServiceServer
	c *container.Container
}

func (s *taskServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	t, err := s.c.TaskLogic.GetByID(ctx, req.GetId())
	if err != nil {
		if isNotFound(err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetTaskResponse{
		Id:      t.Id,
		TokenId: t.TokenId,
		Task:    t.Task,
	}, nil
}

func (s *taskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	t := &models.Task{TokenId: req.GetTokenId(), Task: req.GetTask()}
	if err := s.c.TaskLogic.Create(ctx, t); err != nil {
		if isNotFound(err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if isConflict(err) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CreateTaskResponse{Id: t.Id}, nil
}

func isNotFound(err error) bool {
	return errors.Is(err, domain.ErrNotFound)
}

func isConflict(err error) bool {
	return errors.Is(err, domain.ErrConflict)
}

func NewServer(c *container.Container) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterTokenServiceServer(s, &tokenServer{c: c})
	pb.RegisterTaskServiceServer(s, &taskServer{c: c})
	reflection.Register(s)
	return s
}

func Serve(s *grpc.Server, addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(lis)
}
