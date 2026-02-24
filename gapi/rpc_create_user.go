// simple-bank/gapi/rpc_create_user.go
package gapi

import (
	"context"
	"errors"

	"github.com/lib/pq"
	db "github.com/thinhcompany/simple-bank/db/sqlc"
	pbv1 "github.com/thinhcompany/simple-bank/pb/pb/v1"
	"github.com/thinhcompany/simple-bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(
	ctx context.Context,
	req *pbv1.CreateUserRequest,
) (*pbv1.CreateUserResponse, error) {

	// 1. Validate input
	if err := validateCreateUserRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// 2. Hash password
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot hash password")
	}

	// 3. Create user in DB
	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// PostgreSQL unique violation error code
			if pqErr.Code == "23505" {
				return nil, status.Error(
					codes.AlreadyExists,
					"username or email already exists",
				)
			}
		}

		return nil, status.Error(codes.Internal, "failed to create user")
	}

	// 4. Build response
	rsp := &pbv1.CreateUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}

func validateCreateUserRequest(req *pbv1.CreateUserRequest) error {
	if len(req.GetUsername()) < 3 {
		return errors.New("username too short")
	}
	if len(req.GetPassword()) < 6 {
		return errors.New("password too short")
	}
	if !util.IsValidEmail(req.GetEmail()) {
		return errors.New("invalid email")
	}
	return nil
}
