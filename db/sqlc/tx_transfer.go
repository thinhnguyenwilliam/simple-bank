// simple-bank/db/sqlc/tx_transfer.go
package db

import (
	"context"
	"fmt"
)

// Transfer Transaction Logic
// TransferTxParams contains input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult contains result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfers `json:"transfer"`
	FromAccount Accounts  `json:"from_account"`
	ToAccount   Accounts  `json:"to_account"`
	FromEntry   Entries   `json:"from_entry"`
	ToEntry     Entries   `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to another
// It creates a transfer record, adds entries, and updates account balances
func (store *SQLStore) TransferTx(
	ctx context.Context,
	arg TransferTxParams,
) (TransferTxResult, error) {

	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)
		fmt.Println(txName, "[1---] Create transfer")

		// 1. Create transfer
		// Because the fields match exactly, you can convert the struct instead of rebuilding it
		result.Transfer, err = q.CreateTransfer(
			ctx,
			CreateTransferParams(arg),
		)
		if err != nil {
			return err
		}

		// 2. Create entries
		fmt.Println(txName, "[2---] Create entries 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "[3---] Create entries 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// 3. Update balances (deadlock-safe order)
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err =
				addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err =
				addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return err
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Accounts, account2 Accounts, err error) {

	account1, err = q.UpdateAccount(ctx, UpdateAccountParams{
		ID:      accountID1,
		Balance: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.UpdateAccount(ctx, UpdateAccountParams{
		ID:      accountID2,
		Balance: amount2,
	})
	return
}
