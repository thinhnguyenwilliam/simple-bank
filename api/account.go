// simple-bank/api/account.go
package api

import (
	"database/sql"
	"log"
	"net/http"

	db "github.com/thinhcompany/simple-bank/db/sqlc"
	"github.com/thinhcompany/simple-bank/token"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) listAccount(c *gin.Context) {
	var req listAccountRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := c.MustGet("authorization_payload").(*token.Payload)
	username := payload.Username

	arg := db.ListAccountsParams{
		Owner:  username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := s.store.ListAccounts(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, accounts)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getAccount(c *gin.Context) {
	var req getAccountRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := s.store.GetAccount(c, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 🔐 authorization check
	payload := c.MustGet("authorization_payload").(*token.Payload)
	if account.Owner != payload.Username {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "account does not belong to the authenticated user",
		})
		return
	}

	c.JSON(http.StatusOK, account)
}

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) createAccount(c *gin.Context) {
	var req createAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := c.MustGet("authorization_payload").(*token.Payload)
	username := payload.Username

	arg := db.CreateAccountParams{
		Owner:    username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := s.store.CreateAccount(c, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Println("pq error is:", pqErr.Code.Name())

			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			case "unique_violation":
				c.JSON(http.StatusConflict, errorResponse(err))
				return
			}
		}

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, account)
}
