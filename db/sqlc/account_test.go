package db

import (
	"context"
	"testing"
	"time"

	"github.com/AYehia0/go-bk-mst/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		OwnerName: utils.GetRandomOwnerName(),
		Balance:   utils.GetRandomAmount(),
		Currency:  utils.GetRandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	// account checking
	require.Equal(t, arg.OwnerName, account.OwnerName)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	// database specific
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccountById(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccountById(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.OwnerName, account2.OwnerName)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)

	// check if the timestamps within some duration
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestGetAccounts(t *testing.T) {
	// create n test accounts
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	// get 10 accounts
	args := GetAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.GetAccounts(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
