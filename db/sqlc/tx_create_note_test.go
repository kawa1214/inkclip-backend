package db

import (
	"context"
	"testing"

	"github.com/bookmark-manager/bookmark-manager/util"
	"github.com/google/uuid"
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
	webs := make([]Web, n)
	websArg := make([]uuid.UUID, n)
	for i := 0; i < n; i++ {
		arg := CreateWebParams{
			UserID:       user.ID,
			Url:          util.RandomUrl(),
			Title:        util.RandomName(),
			ThumbnailUrl: util.RandomThumbnailUrl(),
		}
		web, err := store.CreateWeb(context.Background(), arg)
		require.NoError(t, err)

		webs[i] = web
		websArg[i] = web.ID
	}

	arg := TxCreateNoteParams{
		CreateNoteParams: noteArg,
		WebIds:           websArg,
	}
	result, err := store.TxCreateNote(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.NotEmpty(t, result.Note)
	require.Equal(t, len(result.Webs), n)

	for _, web := range result.Webs {
		require.NotEmpty(t, web)
		require.Equal(t, web.UserID, user.ID)
	}
}
