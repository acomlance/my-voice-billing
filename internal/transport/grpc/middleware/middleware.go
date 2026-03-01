package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryRecovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = status.Error(codes.Internal, fmt.Sprintf("panic: %v", r))
		}
	}()
	return handler(ctx, req)
}

func UnaryTiming(isDev bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		dur := time.Since(start)
		log.Debug().Str("method", info.FullMethod).Dur("duration", dur).Msg("grpc request")
		if isDev {
			_ = grpc.SetHeader(ctx, metadata.Pairs("x-response-time", dur.String()))
		}
		return resp, err
	}
}
