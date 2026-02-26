// simple-bank/gapi/rpc_login_user.go
package gapi

import (
	"context"

	db "github.com/thinhcompany/simple-bank/db/sqlc"
	pb "github.com/thinhcompany/simple-bank/pb/pb/v1"
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
	if err := util.CheckPassword(req.GetPassword(), user.HashedPassword); err != nil {
		return nil, status.Error(
			codes.Unauthenticated,
			"incorrect password",
		)
	}

	// 4. Create access token
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			"cannot create access token",
		)
	}

	// 5. Create refresh token
	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			"cannot create refresh token",
		)
	}

	// 6. Store session
	mtdt := s.extractMetatada(ctx)
	_, err = s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			"failed to create session",
		)
	}

	// 7. Build gRPC response
	rsp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		SessionId:             refreshPayload.ID.String(),
	}

	return rsp, nil
}
