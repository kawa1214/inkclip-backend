package db

import (
	"context"
	"testing"

	"github.com/bookmark-manager/bookmark-manager/util"
	"github.com/stretchr/testify/require"
)

func TestTxCreateNote(t *testing.T) {
	store := NewStore(testDB)
	user := createRandomUser(t)

	noteArg := CreateNoteParams{
		Title:   util.RandomString(6),
		Content: util.RandomString(6),
		UserID:  user.ID,
	}
	n := 5
	websArg := make([]CreateWebParams, n)
	for i := 0; i < n; i++ {
		websArg[i] = CreateWebParams{
			UserID:       user.ID,
			Url:          util.RandomUrl(),
			Title:        util.RandomName(),
			ThumbnailUrl: util.RandomThumbnailUrl(),
		}
	}

	arg := TxCreateNoteParams{
		createNoteParams:    noteArg,
		CreateWebParamsList: websArg,
	}
	result, err := store.TxCreateNote(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.NotEmpty(t, result.note)
	require.Equal(t, len(result.webs), n)

	for _, web := range result.webs {
		require.NotEmpty(t, web)
		require.Equal(t, web.UserID, user.ID)
	}
}
