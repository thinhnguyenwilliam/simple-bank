// simple-bank/api/transfer.go
package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/thinhcompany/simple-bank/db/sqlc"
	"github.com/thinhcompany/simple-bank/token"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(c *gin.Context) {
	var req transferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, ok := server.validAccount(c, req.FromAccountID, req.Currency)
	if !ok {
		return
	}

	payload := c.MustGet("authorization_payload").(*token.Payload)
	username := payload.Username
	if fromAccount.Owner != username {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "from account does not belong to the authenticated user",
		})
		return
	}

	if _, ok := server.validAccount(c, req.ToAccountID, req.Currency); !ok {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, result)
}

// What this function should do:
// Get account by ID
// If not found → 404
// If currency mismatch → 400
// Otherwise → return the account
// Returning just bool is limiting, so the best practice is to return the db.Accounts and a bool.
func (server *Server) validAccount(
	c *gin.Context,
	accountID int64,
	currency string,
) (db.Accounts, bool) {

	account, err := server.store.GetAccount(c, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "account not found",
			})
			return db.Accounts{}, false
		}

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return db.Accounts{}, false
	}

	if account.Currency != currency {
		c.JSON(http.StatusBadRequest, gin.H{
			"account_id":       account.ID,
			"account_currency": account.Currency,
			"request_currency": currency,
			"error": fmt.Sprintf(
				"account[%d] currency mismatch: %s vs %s",
				account.ID,
				account.Currency,
				currency,
			),
		})
		return db.Accounts{}, false
	}

	return account, true
}
