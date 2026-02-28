package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"my-voice-billing/internal/config"
	"my-voice-billing/internal/container"
	"my-voice-billing/internal/tools"
	"my-voice-billing/internal/transport/grpc"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.Load("config/app.yml")
	if err != nil {
		fmt.Fprintln(os.Stderr, "load config:", err)
		os.Exit(1)
	}
	tools.InitFromConfig(cfg.Log.Level, cfg.Log.Dir)

	ctx := context.Background()
	if err := container.Init(ctx, cfg); err != nil {
		log.Fatal().Err(err).Msg("init container")
	}
	defer container.Shutdown()

	grpcServer := grpc.NewServer(container.Get())
	port := cfg.Grpc.Port
	if port <= 0 {
		port = 50051
	}
	addr := fmt.Sprintf(":%d", port)
	go func() {
		if err := grpc.Serve(grpcServer, addr); err != nil {
			log.Fatal().Err(err).Str("addr", addr).Msg("grpc serve")
		}
	}()

	log.Info().Str("addr", addr).Msg("gRPC server started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Info().Msg("shutdown signal received")
	grpcServer.GracefulStop()
}
