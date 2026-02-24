// simple-bank/gapi/server.go
package gapi

import (
	db "github.com/thinhcompany/simple-bank/db/sqlc"

	pbv1 "github.com/thinhcompany/simple-bank/pb/pb/v1"
	"github.com/thinhcompany/simple-bank/token"

	"github.com/thinhcompany/simple-bank/util"
)

type Server struct {
	pbv1.UnimplementedSimpleBankServiceServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
