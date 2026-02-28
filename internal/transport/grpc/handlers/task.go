package handlers

import (
	"context"

	"my-voice-billing/internal/container"
	"my-voice-billing/internal/models"
	"my-voice-billing/internal/transport/grpc/pb"
)

type TaskServer struct {
	pb.UnimplementedTaskServiceServer
	c *container.Container
}

func NewTaskServer(c *container.Container) *TaskServer {
	return &TaskServer{c: c}
}

func (s *TaskServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	t, err := s.c.TaskLogic.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, handleErr(err, "GetTask")
	}
	return &pb.GetTaskResponse{
		Id:      t.Id,
		TokenId: t.TokenId,
		Task:    t.Task,
	}, nil
}

func (s *TaskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	t := &models.Task{TokenId: req.GetTokenId(), Task: req.GetTask()}
	if err := s.c.TaskLogic.Create(ctx, t); err != nil {
		return nil, handleErr(err, "CreateTask")
	}
	return &pb.CreateTaskResponse{Id: t.Id}, nil
}
