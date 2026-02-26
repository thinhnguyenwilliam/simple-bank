// simple-bank/db/sqlc/tx_create_user.go
package db

import "context"

type TxCreateUserParams struct {
	CreateUserParams
	AfterCreate func(user Users) error
}

type TxCreateUserResult struct {
	User Users
}

func (store *SQLStore) TxCreateUser(
	ctx context.Context,
	arg TxCreateUserParams,
) (TxCreateUserResult, error) {
	var result TxCreateUserResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		// 👇 optional callback
		if arg.AfterCreate != nil {
			return arg.AfterCreate(result.User)
		}

		return nil
	})

	return result, err
}
