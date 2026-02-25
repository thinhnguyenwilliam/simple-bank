package stream

import (
	"log"
	"time"

	"google.golang.org/grpc"
)

func LoggingInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {

		start := time.Now()
		log.Printf("➡️ stream start: %s", info.FullMethod)

		err := handler(srv, ss)

		log.Printf("⬅️ stream end: %s (%v)", info.FullMethod, time.Since(start))
		return err
	}
}
