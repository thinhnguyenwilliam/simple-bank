package api

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"simple-bank/token"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "Bearer"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if authorizationHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(
				ErrMissingAuthorizationHeader,
			))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(
				ErrInvalidAuthorizationHeader,
			))
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != "bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(
				ErrUnsupportedAuthorizationType,
			))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// save payload to context
		ctx.Set("authorization_payload", payload)

		ctx.Next()
	}
} // addAuthorization adds Authorization header to request
func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	accessToken, payload, err := tokenMaker.CreateToken(
		username,
		duration,
	)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)
	require.NotEmpty(t, payload)

	authorizationHeaderValue := fmt.Sprintf(
		"%s %s",
		authorizationType,
		accessToken,
	)
	request.Header.Set(authorizationHeaderKey, authorizationHeaderValue)
}
