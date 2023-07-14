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
			transferArgs := TransferTxParams{
				FromAccountId: acc1.ID,
				ToAccountId:   acc2.ID,
				Amount:        amount,
			}
			transferRes, err := store.TransferTransaction(context.Background(), transferArgs)
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

	require.Equal(t, updatedAcc1.Balance, acc1.Balance-int64(numConcurrent)*amount)
	require.Equal(t, updatedAcc2.Balance, acc2.Balance-int64(numConcurrent)*amount)

}
