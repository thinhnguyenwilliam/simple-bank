package unary

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		start := time.Now()
		log.Printf("➡️ unary call: %s", info.FullMethod)

		resp, err := handler(ctx, req)

		log.Printf("⬅️ unary done: %s (%v)", info.FullMethod, time.Since(start))
		return resp, err
	}
}
