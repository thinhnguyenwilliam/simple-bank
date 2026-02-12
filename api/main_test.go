// simple-bank/api/main_test.go
package api

import (
	"os"
	"testing"
	"time"

	db "github.com/thinhcompany/simple-bank/db/sqlc"

	"github.com/thinhcompany/simple-bank/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// NewTestServer creates a test server with mock store
func NewTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32), // 32 chars
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
