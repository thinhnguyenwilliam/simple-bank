// simple-bank/api/server.go
package api

import (
	db "simple-bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}

	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts", server.listAccount)
	router.GET("/accounts/:id", server.getAccount)

	server.router = router
	return server
}

// Run starts the HTTP server
func (s *Server) Run(address string) error {
	return s.router.Run(address)
}
