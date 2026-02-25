// simple-bank/service/laptop_server.go
package service

import (
	"context"

	"github.com/google/uuid"

	pb "github.com/thinhcompany/simple-bank/pb/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LaptopServer struct {
	pb.UnimplementedLaptopServiceServer
	store LaptopStore
}

func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{
		store: store,
	}
}

func (s *LaptopServer) CreateLaptop(
	ctx context.Context,
	req *pb.CreateLaptopRequest,
) (*pb.CreateLaptopResponse, error) {

	laptop := req.GetLaptop()
	if laptop == nil {
		return nil, status.Error(codes.InvalidArgument, "laptop is required")
	}

	if laptop.Id == "" {
		laptop.Id = uuid.NewString()
	}

	if _, err := uuid.Parse(laptop.Id); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid laptop ID")
	}

	if laptop.PriceUsd <= 0 {
		return nil, status.Error(codes.InvalidArgument, "price must be positive")
	}

	laptop.UpdatedAt = timestamppb.Now()

	err := s.store.Save(laptop)
	if err != nil {
		if err == ErrAlreadyExists {
			return nil, status.Errorf(
				codes.AlreadyExists,
				"laptop %s already exists",
				laptop.Id,
			)
		}
		return nil, status.Error(codes.Internal, "cannot save laptop")
	}

	return &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}, nil
}
