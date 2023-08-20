package token

import (
	"testing"
	"time"

	"github.com/AYehia0/go-bk-mst/utils"
	"github.com/stretchr/testify/require"
)

func TestValidateToken(t *testing.T) {
	secret := utils.RandomString(32)
	creator, err := NewJWTCreator(secret)

	require.NoError(t, err)

	username := utils.GetRandomOwnerName()
	token, err := creator.Create(
		username,
		time.Minute,
	)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	createdAt := time.Now()
	expiredAt := createdAt.Add(time.Minute)

	payload, err := creator.Verify(token)
	require.NotEmpty(t, payload)

	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
	require.WithinDuration(t, createdAt, payload.CreatedAt, time.Second)
	
	// check if the UUID was given
	require.NotZero(t, payload.Id)
}

func TestInvalidToken(t *testing.T) {
	t.Run("TokenExpired", func(t *testing.T) {
		secret := utils.RandomString(32)
		creator, err := NewJWTCreator(secret)

		require.NoError(t, err)

		username := utils.GetRandomOwnerName()
		token, err := creator.Create(
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

	t.Run("InvalidAlgo", func(t *testing.T) {

		payload, err := NewPayload(utils.GetRandomOwnerName(), time.Minute)
		require.NoError(t, err)

		// create a new token
		token := jwt.NewWithClaims(jwt.SigningMethodNone, payload)

		// sign/create the token
		fToken, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
		require.NoError(t, err)

		// trying to verify that token
		creator, err := NewJWTCreator(utils.RandomString(32))
		payload, err = creator.Verify(fToken)
		require.Error(t, err)
		require.ErrorIs(t, err, TokenInvalidError)
	})
}
