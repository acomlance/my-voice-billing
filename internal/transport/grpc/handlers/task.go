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
	return taskToProto(t), nil
}

func (s *TaskServer) GetTaskByToken(ctx context.Context, req *pb.GetTaskByTokenRequest) (*pb.GetTaskResponse, error) {
	t, err := s.c.TaskLogic.GetByToken(ctx, req.GetToken())
	if err != nil {
		return nil, handleErr(err, "GetTaskByToken")
	}
	return taskToProto(t), nil
}

func (s *TaskServer) ListTasksByAccountID(ctx context.Context, req *pb.ListTasksByAccountIDRequest) (*pb.ListTasksByAccountIDResponse, error) {
	list, err := s.c.TaskLogic.ListByAccountID(ctx, req.GetAccountId())
	if err != nil {
		return nil, handleErr(err, "ListTasksByAccountID")
	}
	out := make([]*pb.GetTaskResponse, 0, len(list))
	for i := range list {
		out = append(out, taskToProto(&list[i]))
	}
	return &pb.ListTasksByAccountIDResponse{Tasks: out}, nil
}

func (s *TaskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	t := &models.Task{
		Token:          req.GetToken(),
		AccountId:      req.GetAccountId(),
		ReservedTokens: req.GetReservedTokens(),
	}
	if err := s.c.TaskLogic.Create(ctx, t); err != nil {
		return nil, handleErr(err, "CreateTask")
	}
	return &pb.CreateTaskResponse{Id: t.Id}, nil
}

func (s *TaskServer) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	if err := s.c.TaskLogic.Delete(ctx, req.GetId()); err != nil {
		return nil, handleErr(err, "DeleteTask")
	}
	return &pb.DeleteTaskResponse{}, nil
}

func taskToProto(t *models.Task) *pb.GetTaskResponse {
	return &pb.GetTaskResponse{
		Id:             t.Id,
		Token:          t.Token,
		AccountId:      t.AccountId,
		ReservedTokens: t.ReservedTokens,
	}
}
