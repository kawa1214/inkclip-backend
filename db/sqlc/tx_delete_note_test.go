package db

import (
	"context"
	"testing"

	"github.com/bookmark-manager/bookmark-manager/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTxDeleteNote(t *testing.T) {
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

	createNoteArg := TxCreateNoteParams{
		CreateNoteParams: noteArg,
		WebIds:           websArg,
	}
	result, err := store.TxCreateNote(context.Background(), createNoteArg)
	require.NoError(t, err)

	deleteNoteArg := TxDeleteNoteParams{
		NoteID: result.Note.ID,
	}
	err = store.TxDeleteNote(context.Background(), deleteNoteArg)
	require.NoError(t, err)

	for _, web := range result.Webs {
		arg := GetNoteWebParams{
			NoteID: result.Note.ID,
			WebID:  web.ID,
		}
		_, err := store.GetNoteWeb(context.Background(), arg)
		require.Error(t, err)
	}
	_, err = store.GetNote(context.Background(), result.Note.ID)
	require.Error(t, err)
}
