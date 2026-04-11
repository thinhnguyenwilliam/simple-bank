// simple-bank\api\token.go
package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (s *Server) renewAccessToken(ctx *gin.Context) {

	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 1️⃣ Verify refresh token
	payload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	log.Println("refresh ID:", payload.ID)
	// 2️⃣ Get session from DB
	session, err := s.store.GetSession(ctx, payload.ID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// 3️⃣ Check session validity
	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, errorResponse(
			fmt.Errorf("session is blocked"),
		))
		return
	}

	if session.Username != payload.Username {
		ctx.JSON(http.StatusUnauthorized, errorResponse(
			fmt.Errorf("incorrect session user"),
		))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, errorResponse(
			fmt.Errorf("token reuse detected"),
		))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(
			fmt.Errorf("session expired"),
		))
		return
	}

	// 4️⃣ Create new access token
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		payload.Username,
		"test_role",
		s.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, resp)
}
