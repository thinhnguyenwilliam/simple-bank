package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/thinhcompany/simple-bank/token"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuthRoute(
	router *gin.Engine,
	tokenMaker token.Maker,
) {
	router.GET(
		"/auth",
		authMiddleware(tokenMaker),
		func(ctx *gin.Context) {
			payload := ctx.MustGet("authorization_payload").(*token.Payload)
			ctx.JSON(http.StatusOK, gin.H{
				"username": payload.Username,
			})
		},
	)
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// create token maker
	symmetricKey := "12345678901234567890123456789012"
	tokenMaker, err := token.NewPasetoMaker(symmetricKey)
	require.NoError(t, err)

	username := "thinh"
	duration := time.Minute

	accessToken, payload, err := tokenMaker.CreateToken(username, "test_role", duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	expiredToken, payload, err := tokenMaker.CreateToken(username, "test_role", -time.Minute)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, r *http.Request)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, r *http.Request) {
				r.Header.Set("Authorization", "Bearer "+accessToken)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			name:      "NoAuthorizationHeader",
			setupAuth: func(t *testing.T, r *http.Request) {},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name: "InvalidAuthorizationHeader",
			setupAuth: func(t *testing.T, r *http.Request) {
				r.Header.Set("Authorization", "invalid")
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name: "UnsupportedAuthorizationType",
			setupAuth: func(t *testing.T, r *http.Request) {
				r.Header.Set("Authorization", "Basic "+accessToken)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, r *http.Request) {
				r.Header.Set("Authorization", "Bearer "+expiredToken)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			addAuthRoute(router, tokenMaker)

			rec := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/auth", nil)
			require.NoError(t, err)

			tc.setupAuth(t, req)
			router.ServeHTTP(rec, req)
			tc.checkResponse(t, rec)
		})
	}
}
