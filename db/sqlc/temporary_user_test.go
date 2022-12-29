package db

import (
	"context"
	"testing"
	"time"

	"github.com/inkclip/backend/util"
	"github.com/stretchr/testify/require"
)

func TestCreateTemporaryUser(t *testing.T) {
	createRandomTemporaryUser(t)
}

func TestGetTemporaryUserByEmailAndToken(t *testing.T) {
	tmpUser := createRandomTemporaryUser(t)
	arg := GetTemporaryUserByEmailAndTokenParams{
		Email: tmpUser.Email,
		Token: tmpUser.Token,
	}
	gotTmpuser, err := testQueries.GetTemporaryUserByEmailAndToken(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, gotTmpuser)

	require.Equal(t, tmpUser.Email, gotTmpuser.Email)
	require.Equal(t, tmpUser.HashedPassword, gotTmpuser.HashedPassword)
	require.Equal(t, tmpUser.Token, gotTmpuser.Token)
	require.Equal(t, tmpUser.ExpiresAt, gotTmpuser.ExpiresAt)
	require.WithinDuration(t, tmpUser.CreatedAt, gotTmpuser.CreatedAt, time.Second)
}

func createRandomTemporaryUser(t *testing.T) TemporaryUser {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateTemporaryUserParams{
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
		Token:          util.RandomString(6),
		ExpiresAt:      time.Now().Add(time.Second),
	}

	tmpUser, err := testQueries.CreateTemporaryUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tmpUser)

	require.Equal(t, arg.Email, tmpUser.Email)
	require.Equal(t, arg.HashedPassword, tmpUser.HashedPassword)
	require.Equal(t, arg.Token, tmpUser.Token)
	require.Equal(t, arg.ExpiresAt.UTC(), tmpUser.ExpiresAt.UTC())

	require.NotZero(t, tmpUser.CreatedAt)
	return tmpUser
}
