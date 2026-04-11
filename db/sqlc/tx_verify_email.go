// simple-bank/db/sqlc/tx_verify_email.go
package db

import (
	"context"
	"errors"
)

type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	User        Users
	VerifyEmail VerifyEmails
}

var (
	ErrInvalidVerifyEmail = errors.New("invalid verify email")
)

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// ✅ 1. get valid verify email
		verifyEmail, err := q.GetVerifyEmailBySecretCode(ctx, arg.SecretCode)
		if err != nil {
			return err
		}

		// ✅ 2. check email_id match
		if verifyEmail.ID != arg.EmailId {
			return ErrInvalidVerifyEmail
		}

		// ✅ 3. mark used
		err = q.MarkVerifyEmailUsed(ctx, verifyEmail.ID)
		if err != nil {
			return err
		}

		// ✅ 4. update user
		result.User, err = q.VerifyUserEmail(ctx, verifyEmail.Username)
		if err != nil {
			return err
		}

		result.VerifyEmail = verifyEmail

		return nil
	})

	return result, err
}
