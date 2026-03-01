// simple-bank/gapi/error.go
package gapi

import (
	pb "github.com/thinhcompany/simple-bank/pb/pb/v1"
	"github.com/thinhcompany/simple-bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// invalidArgumentError creates a gRPC InvalidArgument error with field violations
func invalidArgumentError(
	violations []*errdetails.BadRequest_FieldViolation,
) error {
	st := status.New(codes.InvalidArgument, "invalid parameters")

	badRequest := &errdetails.BadRequest{
		FieldViolations: violations,
	}

	st, err := st.WithDetails(badRequest)
	if err != nil {
		return status.Errorf(codes.Internal, "cannot attach error details")
	}

	return st.Err()
}

func validateCreateUserRequest(
	req *pb.CreateUserRequest,
) []*errdetails.BadRequest_FieldViolation {

	var violations []*errdetails.BadRequest_FieldViolation

	if err := val.ValidateUsername(req.Username); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidateFullName(req.FullName); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := val.ValidateEmail(req.Email); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := val.ValidatePassword(req.Password); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}

func fieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
		Reason:      "INVALID_FORMAT",
		// LocalizedMessage: &errdetails.LocalizedMessage{
		// 	Locale:  "vi-VN",
		// 	Message: "My custom hihi",
		// },
	}
}
