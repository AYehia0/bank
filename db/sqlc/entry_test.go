package db

import (
	"context"
	"testing"

	"github.com/AYehia0/go-bk-mst/utils"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, acc Account) Entry {
	randAmount := utils.GetRandomAmount()
	args := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    randAmount,
	}
	entry, err := testQueries.CreateEntry(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, acc.ID)
	require.Equal(t, entry.Amount, randAmount)

	return entry
}

func TestCreateEntry(t *testing.T) {
	acc := createRandomAccount(t)
	createRandomEntry(t, acc)
}

func TestGetEntry(t *testing.T) {
	acc := createRandomAccount(t)
	e := createRandomEntry(t, acc)
	entry, err := testQueries.GetEntryById(context.Background(), e.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, e.ID, entry.ID)
	require.Equal(t, e.AccountID, entry.AccountID)
}

func TestGetEntries(t *testing.T) {
	acc := createRandomAccount(t)
	for i := 0; i < 5; i++ {
		createRandomEntry(t, acc)
	}

	entryParams := GetEntriesParams{
		AccountID: acc.ID,
		Limit:     5,
		Offset:    0,
	}
	entries, err := testQueries.GetEntries(context.Background(), entryParams)

	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, acc.ID, entry.AccountID)
	}
}
