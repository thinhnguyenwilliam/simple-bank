package interceptor

import (
	s "github.com/thinhcompany/simple-bank/interceptor/stream"
	u "github.com/thinhcompany/simple-bank/interceptor/unary"
	"github.com/thinhcompany/simple-bank/token"

	"google.golang.org/grpc"
)

func Unary() grpc.ServerOption {
	jwtMaker, _ := token.NewJWTMaker("super-secret-key")
	return grpc.ChainUnaryInterceptor(
		// u.RecoveryInterceptor(),
		u.AuthInterceptor(jwtMaker),
		u.LoggingInterceptor(),
	)
}

func Stream() grpc.ServerOption {
	return grpc.ChainStreamInterceptor(
		s.LoggingInterceptor(),
		// s.AuthInterceptor(),
		// s.RecoveryInterceptor(),
	)
}
