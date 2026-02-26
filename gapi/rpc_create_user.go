package gapi

import (
	"context"
	"errors"
	"time"

	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"

	db "github.com/thinhcompany/simple-bank/db/sqlc"
	pb "github.com/thinhcompany/simple-bank/pb/pb/v1"
	"github.com/thinhcompany/simple-bank/util"
	"github.com/thinhcompany/simple-bank/worker"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (*pb.CreateUserResponse, error) {

	// 1. Validate input
	if err := validateCreateUserRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// 2. Hash password
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot hash password")
	}

	// 3. Create user (DB ONLY)
	txResult, err := s.store.TxCreateUser(ctx, db.TxCreateUserParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPassword,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, status.Error(
				codes.AlreadyExists,
				"username or email already exists",
			)
		}
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	// 4. Enqueue verify-email task (AFTER COMMIT)
	taskPayload := &worker.PayloadSendVerifyEmail{
		Username: txResult.User.Username,
		Email:    txResult.User.Email,
	}

	opts := []asynq.Option{
		asynq.Queue("critical"),
		asynq.MaxRetry(10),
		asynq.Timeout(30 * time.Second),
		asynq.ProcessIn(10 * time.Second),
		asynq.Deadline(time.Now().Add(1 * time.Hour)),
	}

	err = s.taskDistributor.DistributeTaskSendVerifyEmail(
		context.Background(), // IMPORTANT
		taskPayload,
		opts...,
	)
	if err != nil {
		log.Error().
			Err(err).
			Str("username", txResult.User.Username).
			Msg("failed to enqueue verify email task")
		// ❗ do NOT return error — user already created
	}

	// 5. Build response
	rsp := &pb.CreateUserResponse{
		User: convertUser(txResult.User),
	}

	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) error {
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
