package db

import (
	"context"
	"testing"
	"time"

	"github.com/bookmark-manager/bookmark-manager/util"
	"github.com/google/uuid"
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

func TestListWebByNoteIds(t *testing.T) {
	user := createRandomUser(t)
	note1 := createRandomNote(t, user)
	note2 := createRandomNote(t, user)
	note3 := createRandomNote(t, user)
	n := 5
	for i := 0; i < n; i++ {
		web1 := createRandomWeb(t, user)
		createRandomNoteWeb(t, note1, web1)
		web2 := createRandomWeb(t, user)
		createRandomNoteWeb(t, note2, web2)
		web3 := createRandomWeb(t, user)
		createRandomNoteWeb(t, note3, web3)
	}

	webs, err := testQueries.ListWebByNoteIds(context.Background(), []uuid.UUID{
		note1.ID,
		note2.ID,
	})
	require.NoError(t, err)
	require.Len(t, webs, n*2)

	for _, web := range webs {
		require.NotEmpty(t, web)
		require.Equal(t, user.ID, web.UserID)
		require.NotEmpty(t, web.NoteID)
	}
}

func TestListWebByNoteId(t *testing.T) {
	user := createRandomUser(t)
	note1 := createRandomNote(t, user)
	note2 := createRandomNote(t, user)
	n := 5
	for i := 0; i < n; i++ {
		web1 := createRandomWeb(t, user)
		createRandomNoteWeb(t, note1, web1)
		web2 := createRandomWeb(t, user)
		createRandomNoteWeb(t, note2, web2)
	}

	webs, err := testQueries.ListWebByNoteId(context.Background(), note1.ID)
	require.NoError(t, err)
	require.Len(t, webs, n)

	for _, web := range webs {
		require.NotEmpty(t, web)
		require.Equal(t, user.ID, web.UserID)
	}
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
