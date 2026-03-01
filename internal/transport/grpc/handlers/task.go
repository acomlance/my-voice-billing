package handlers

import (
	"context"

	"my-voice-billing/internal/models"
	"my-voice-billing/internal/transport/grpc/pb"
)

type TaskLogic interface {
	Create(ctx context.Context, t *models.Task) error
	Delete(ctx context.Context, token string, closedTokens int64) error
}

type TaskServer struct {
	pb.UnimplementedTaskServiceServer
	taskLogic TaskLogic
}

func NewTaskServer(taskLogic TaskLogic) *TaskServer {
	return &TaskServer{taskLogic: taskLogic}
}

func (s *TaskServer) Create(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	t := &models.Task{
		AccountId:      req.GetAccountId(),
		ReservedTokens: req.GetReservedTokens(),
	}
	if err := s.taskLogic.Create(ctx, t); err != nil {
		return nil, handleErr(err, "Create")
	}
	return &pb.CreateTaskResponse{Token: t.Token}, nil
}

func (s *TaskServer) Delete(ctx context.Context, req *pb.DeleteTaskByTokenRequest) (*pb.DeleteTaskByTokenResponse, error) {
	if err := s.taskLogic.Delete(ctx, req.GetToken(), req.GetClosedTokens()); err != nil {
		return nil, handleErr(err, "Delete")
	}
	return &pb.DeleteTaskByTokenResponse{}, nil
}
