package handlers

import (
	"context"

	"my-voice-billing/internal/models"
	"my-voice-billing/internal/transport/grpc/pb"
)

var _ pb.TaskServiceServer = (*TaskServer)(nil)

type TaskLogic interface {
	Create(ctx context.Context, t *models.Task) error
	Close(ctx context.Context, token string, closedTokens int64) error
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

func (s *TaskServer) CloseTask(ctx context.Context, req *pb.CloseTaskRequest) (*pb.CloseTaskResponse, error) {
	if err := s.taskLogic.Close(ctx, req.GetToken(), req.GetClosedTokens()); err != nil {
		return nil, handleErr(err, "CloseTask")
	}
	return &pb.CloseTaskResponse{}, nil
}
