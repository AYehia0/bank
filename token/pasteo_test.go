package token

import (
	"testing"
	"time"

	"github.com/AYehia0/go-bk-mst/utils"
	"github.com/stretchr/testify/require"
)

func TestValidatePasteoToken(t *testing.T) {
	secret := utils.RandomString(32)
	creator, err := NewPasteoToken(secret)

	require.NoError(t, err)

	username := utils.GetRandomOwnerName()
	token, payload, err := creator.Create(
		username,
		time.Minute,
	)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, token)

	createdAt := time.Now()
	expiredAt := createdAt.Add(time.Minute)

	payload, err = creator.Verify(token)
	require.NotEmpty(t, payload)

	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
	require.WithinDuration(t, createdAt, payload.CreatedAt, time.Second)

	// check if the UUID was given
	require.NotZero(t, payload.Id)
}

func TestInvalidPasteoToken(t *testing.T) {
	t.Run("TokenExpired", func(t *testing.T) {
		secret := utils.RandomString(32)
		creator, err := NewPasteoToken(secret)

		require.NoError(t, err)

		username := utils.GetRandomOwnerName()
		token, _, err := creator.Create(
			username,
			-time.Minute,
		)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		payload, err := creator.Verify(token)
		require.Error(t, err)
		require.ErrorIs(t, err, TokenExpiredError)
		require.Nil(t, payload)
	})
}
