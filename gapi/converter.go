// simple-bank/gapi/converter.go
package gapi

import (
	db "github.com/thinhcompany/simple-bank/db/sqlc"
	pbv1 "github.com/thinhcompany/simple-bank/pb/pb/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Convert DB model → protobuf
func convertUser(user db.Users) *pbv1.User {
	return &pbv1.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
