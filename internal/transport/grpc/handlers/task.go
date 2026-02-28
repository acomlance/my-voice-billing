package handlers

import (
	"context"

	"my-voice-billing/internal/models"
	"my-voice-billing/internal/transport/grpc/pb"
)

// TaskLogic — интерфейс логики задач для внедрения в хендлер
type TaskLogic interface {
	Create(ctx context.Context, t *models.Task) error
	DeleteByToken(ctx context.Context, token string, closedTokens int64) error
}

type TaskServer struct {
	pb.UnimplementedTaskServiceServer
	taskLogic TaskLogic
}

func NewTaskServer(taskLogic TaskLogic) *TaskServer {
	return &TaskServer{taskLogic: taskLogic}
}

func (s *TaskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	t := &models.Task{
		AccountId:      req.GetAccountId(),
		ReservedTokens: req.GetReservedTokens(),
	}
	if err := s.taskLogic.Create(ctx, t); err != nil {
		return nil, handleErr(err, "CreateTask")
	}
	return &pb.CreateTaskResponse{Token: t.Token}, nil
}

func (s *TaskServer) DeleteTaskByToken(ctx context.Context, req *pb.DeleteTaskByTokenRequest) (*pb.DeleteTaskByTokenResponse, error) {
	if err := s.taskLogic.DeleteByToken(ctx, req.GetToken(), req.GetClosedTokens()); err != nil {
		return nil, handleErr(err, "DeleteTaskByToken")
	}
	return &pb.DeleteTaskByTokenResponse{}, nil
}
