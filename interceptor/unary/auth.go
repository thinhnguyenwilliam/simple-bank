package unary

import (
	"context"
	"strings"

	"github.com/thinhcompany/simple-bank/token"
	"github.com/thinhcompany/simple-bank/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var publicMethods = map[string]bool{
	"/pb.v1.LaptopService/CreateLaptop": false, // private
	"/pb.v1.AuthService/Login":          true,  // public
	"/pb.v1.AuthService/Register":       true,
}

func AuthInterceptor(jwtMaker token.Maker) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization header is missing")
		}

		fields := strings.Fields(values[0])
		if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization format")
		}

		accessToken := fields[1]

		payload, err := jwtMaker.VerifyToken(accessToken)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		ctx = util.ContextWithUser(ctx, payload.Username)
		return handler(ctx, req)
	}
}
