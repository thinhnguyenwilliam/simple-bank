// simple-bank/gapi/rpc_update_user.go
package gapi

import (
	"context"
	"database/sql"

	db "github.com/thinhcompany/simple-bank/db/sqlc"
	pbv1 "github.com/thinhcompany/simple-bank/pb/pb/v1"
	"github.com/thinhcompany/simple-bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateUser(ctx context.Context, req *pbv1.UpdateUserRequest) (*pbv1.UpdateUserResponse, error) {

	arg := db.UpdateUserParams{
		Username: req.Username,
	}

	if req.FullName != nil {
		arg.FullName = sql.NullString{
			String: *req.FullName,
			Valid:  true,
		}
	}

	if req.Email != nil {
		arg.Email = sql.NullString{
			String: *req.Email,
			Valid:  true,
		}
	}

	if req.Password != nil {
		if *req.Password == "" {
			return nil, status.Errorf(codes.InvalidArgument, "password cannot be empty")
		}

		hashedPassword, err := util.HashPassword(*req.Password)
		if err != nil {
			return nil, err
		}

		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}
	}

	user, err := s.store.UpdateUser(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &pbv1.UpdateUserResponse{
		User: convertUser(user),
	}, nil
}
