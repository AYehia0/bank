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
}
