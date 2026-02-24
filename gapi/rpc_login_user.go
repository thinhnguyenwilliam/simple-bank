// simple-bank/gapi/rpc_login_user.go
package gapi

import (
	"context"

	"github.com/thinhcompany/simple-bank/pb"
	"github.com/thinhcompany/simple-bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(
	ctx context.Context,
	req *pb.LoginUserRequest,
) (*pb.LoginUserResponse, error) {

	// 1. Validate input
	if req.GetUsername() == "" || req.GetPassword() == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"username and password are required",
		)
	}

	// 2. Get user from DB
	user, err := s.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		return nil, status.Error(
			codes.NotFound,
			"user not found",
		)
	}

	// 3. Check password
	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Error(
			codes.Unauthenticated,
			"incorrect password",
		)
	}

	// 4. Create access token
	accessToken, payload, err := s.tokenMaker.CreateToken(
		user.Username,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			"cannot create access token",
		)
	}

	// 5. Build response
	rsp := &pb.LoginUserResponse{
		User:                 convertUser(user),
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(payload.ExpiredAt),
	}

	return rsp, nil
}
