package db

import (
	"context"
	"testing"

	"github.com/AYehia0/go-bk-mst/utils"
	"github.com/stretchr/testify/require"
)

// acc1: FromAccount and acc2: ToAccount
func createRandomTransfer(t *testing.T, acc1 Account, acc2 Account) Transfer {
	args := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        utils.GetRandomAmount(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, acc1.ID)
	require.Equal(t, transfer.ToAccountID, acc2.ID)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	createRandomTransfer(t, acc1, acc2)
}

func TestGetTransferById(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	tr := createRandomTransfer(t, acc1, acc2)

	transfer, err := testQueries.GetTransferById(context.Background(), tr.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, acc1.ID)
	require.Equal(t, transfer.ToAccountID, acc2.ID)
}

func TestGetTransfers(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	for i := 0; i < 5; i++ {
		createRandomTransfer(t, acc1, acc2)
	}
	args := GetTransfersParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Limit:         5,
		Offset:        0,
	}
	transfers, err := testQueries.GetTransfers(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.Equal(t, transfer.FromAccountID, acc1.ID)
		require.Equal(t, transfer.ToAccountID, acc2.ID)
	}
}
