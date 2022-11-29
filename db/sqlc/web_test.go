package db

import (
	"context"
	"testing"
	"time"

	"github.com/bookmark-manager/bookmark-manager/util"
	"github.com/stretchr/testify/require"
)

func TestCreateWeb(t *testing.T) {
	user := createRandomUser(t)
	createRandomWeb(t, user)
}

func TestGetWeb(t *testing.T) {
	user := createRandomUser(t)

	createdWeb := createRandomWeb(t, user)
	getWeb, err := testQueries.GetWeb(context.Background(), createdWeb.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getWeb)

	require.Equal(t, createdWeb.ID, getWeb.ID)
	require.Equal(t, createdWeb.UserID, getWeb.UserID)
	require.Equal(t, createdWeb.Url, getWeb.Url)
	require.Equal(t, createdWeb.Title, getWeb.Title)
	require.Equal(t, createdWeb.ThumbnailUrl, getWeb.ThumbnailUrl)
	require.WithinDuration(t, createdWeb.CreatedAt, getWeb.CreatedAt, time.Second)
}

func TestListWebs(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		createRandomWeb(t, user)
	}

	arg := ListWebsByUserIdParams{
		UserID: user.ID,
		Limit:  5,
		Offset: 5,
	}
	webs, err := testQueries.ListWebsByUserId(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, webs)

	for _, web := range webs {
		require.NotEmpty(t, web)
	}
}

func TestDeleteWeb(t *testing.T) {
	user := createRandomUser(t)
	web := createRandomWeb(t, user)

	beforeDeleteWeb, err := testQueries.GetWeb(context.Background(), web.ID)
	require.NoError(t, err)
	require.Equal(t, web.ID, beforeDeleteWeb.ID)

	testQueries.DeleteWeb(context.Background(), web.ID)
	_, err = testQueries.GetWeb(context.Background(), web.ID)
	require.Error(t, err)
}

func createRandomWeb(t *testing.T, user User) Web {
	ThumbnailURL := util.RandomThumbnailUrl()
	arg := CreateWebParams{
		UserID:       user.ID,
		Url:          util.RandomUrl(),
		Title:        util.RandomName(),
		ThumbnailUrl: ThumbnailURL,
	}

	web, err := testQueries.CreateWeb(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, web)

	require.Equal(t, arg.UserID, web.UserID)
	require.Equal(t, arg.Url, web.Url)
	require.Equal(t, arg.Title, web.Title)
	require.Equal(t, arg.ThumbnailUrl, web.ThumbnailUrl)

	require.NotZero(t, web.CreatedAt)
	return web
}
