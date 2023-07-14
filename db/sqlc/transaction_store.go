// This is the store to create transactions on the db
// The store provides all the db functionality that require ACID to be true.
package db

import (
	"context"
	"database/sql"
	"fmt"
)

// provides all the functions to execute db queries and transactions
type Store struct {
	*Queries //composition over inhertance
	db       *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execute the transaction
// callback on the same function
func (store *Store) execTransaction(ctx context.Context, fn func(*Queries) error) error {
	// store.db.Begin() uses the background context
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	// get query object
	query := New(tx)
	err = fn(query)

	//rollback on any error
	if err != nil {
		if rbError := tx.Rollback(); rbError != nil {
			return fmt.Errorf("Transaction Error: %v, Rollback Error: %v", err, rbError)
		}
		return err
	}

	// commit if all operations were successful
	return tx.Commit()
}

// contains the input params for a successful transaction
type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// contains the output result for a successful transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	ToEntry     Entry    `json:"to_entry"`
	FromEntry   Entry    `json:"from_entry"`
}

var txKey = struct{}{}

func (store *Store) TransferTransaction(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var res TransferTxResult

	// go Closures : https://betterprogramming.pub/closures-made-simple-with-golang-69db3017cd7b?gi=48e0b91f624a
	err := store.execTransaction(ctx, func(q *Queries) error {
		var err error

		// 1. create a transfer
		res.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		// 2. create entry to the account who received the amount with negative amount
		res.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		// 3. create entry from the account who sent the amount with positive amount
		res.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// get the accounts from the database, then add/subtract from their balance (need proper locking mechanism)
		senderAccount, err := q.GetAccountByIdForUpdate(ctx, arg.FromAccountId)
		if err != nil {
			return err
		}

		res.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      senderAccount.ID,
			Balance: senderAccount.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		receiverAccount, err := q.GetAccountByIdForUpdate(ctx, arg.ToAccountId)
		if err != nil {
			return err
		}

		res.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      receiverAccount.ID,
			Balance: receiverAccount.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})
	return res, err
}
