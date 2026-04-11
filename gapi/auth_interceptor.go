package gapi

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// 👇 map method → roles
		accessibleRoles := map[string][]string{
			"/pb.SimpleBankService/UpdateUser": {"admin", "user"},
		}

		roles, ok := accessibleRoles[info.FullMethod]
		if ok {
			roleMap := make(map[string]bool)
			for _, r := range roles {
				roleMap[r] = true
			}

			payload, err := s.authorizeUser(ctx, roleMap)
			if err != nil {
				return nil, status.Errorf(codes.Unauthenticated, "unauthorized: %v", err)
			}

			// 👇 attach payload vào context
			ctx = context.WithValue(ctx, "user", payload)
		}

		return handler(ctx, req)
	}
}
