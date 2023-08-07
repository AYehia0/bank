package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// moving money from account1 to account2
// it should create a transfer record and 2 entries
// TODO: update the account balance
func TestTransferTransaction(t *testing.T) {
	store := NewStore(testDb)
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	amount := int64(20)
	numConcurrent := 5
	for i := 0; i < numConcurrent; i++ {
		go func() {
			ctx := context.Background()
			transferRes, err := store.TransferTransaction(ctx, TransferTxParams{
				FromAccountId: acc1.ID,
				ToAccountId:   acc2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- transferRes
		}()
	}

	for i := 0; i < numConcurrent; i++ {
		err := <-errs
		require.NoError(t, err)

		transferRes := <-results
		require.NotEmpty(t, transferRes)

		// result checking
		require.Equal(t, transferRes.Transfer.FromAccountID, acc1.ID)
		require.Equal(t, transferRes.Transfer.ToAccountID, acc2.ID)
		require.Equal(t, transferRes.Transfer.Amount, amount)

		// entry
		require.Equal(t, transferRes.FromEntry.AccountID, acc1.ID)
		require.Equal(t, transferRes.ToEntry.AccountID, acc2.ID)

		_, err = store.GetEntryById(context.Background(), transferRes.ToEntry.ID)
		require.NoError(t, err)

		require.NotEmpty(t, transferRes.FromEntry)
		require.NotEmpty(t, transferRes.ToEntry)

		// transfer
		_, err = store.GetTransferById(context.Background(), transferRes.Transfer.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := transferRes.FromAccount
		toAccount := transferRes.ToAccount
		require.NotEmpty(t, fromAccount)
		require.NotEmpty(t, toAccount)
		require.Equal(t, acc1.ID, fromAccount.ID)
		require.Equal(t, acc2.ID, toAccount.ID)

		// check account balance
		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // account1 diff should decrease by nX each time
	}
	updatedAcc1, err := store.GetAccountById(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedAcc2, err := store.GetAccountById(context.Background(), acc2.ID)
	require.NoError(t, err)

	require.Equal(t, acc1.Balance-int64(numConcurrent)*amount, updatedAcc1.Balance)
	require.Equal(t, acc2.Balance+int64(numConcurrent)*amount, updatedAcc2.Balance)

}

// deadlock happens when 2 transactions happen at the same time
// transfer 10$ from account 1 to account 2
// AND transfer 10$ from account 2 to account 1
// run 5 concurrent transactions from 1 to 2 and 2 to 1
// we expect the balance from acc1 to equal balance in acc2
func TestTransferTransactionDeadLock(t *testing.T) {
	store := NewStore(testDb)
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	errs := make(chan error)

	amount := int64(20)
	numConcurrent := 10
	for i := 0; i < numConcurrent; i++ {

		fromAccountId := acc1.ID
		toAccountId := acc2.ID

		if i%2 == 1 {
			fromAccountId = acc2.ID
			toAccountId = acc1.ID
		}
		go func() {
			ctx := context.Background()
			_, err := store.TransferTransaction(ctx, TransferTxParams{
				FromAccountId: fromAccountId,
				ToAccountId:   toAccountId,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < numConcurrent; i++ {
		err := <-errs
		require.NoError(t, err)
	}
	updatedAcc1, err := store.GetAccountById(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedAcc2, err := store.GetAccountById(context.Background(), acc2.ID)
	require.NoError(t, err)

	require.Equal(t, acc1.Balance, updatedAcc1.Balance)
	require.Equal(t, acc2.Balance, updatedAcc2.Balance)
}
