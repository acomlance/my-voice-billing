package grpc

import (
	"net"

	"my-voice-billing/internal/config"
	"my-voice-billing/internal/container"
	"my-voice-billing/internal/transport/grpc/handlers"
	"my-voice-billing/internal/transport/grpc/middleware"
	"my-voice-billing/internal/transport/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewServer(c *container.Container) *grpc.Server {
	var isDev bool
	if c.Config != nil {
		isDev = c.Config.Environment == config.EnvDev
	}
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(middleware.UnaryRecovery, middleware.UnaryTiming(isDev)),
	}
	s := grpc.NewServer(opts...)
	pb.RegisterTaskServiceServer(s, handlers.NewTaskServer(c.TaskLogic))
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
