package db

import (
	"context"
	"testing"

	"github.com/bookmark-manager/bookmark-manager/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTxUpdateNote(t *testing.T) {
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
			Url:          util.RandomURL(),
			Title:        util.RandomName(),
			ThumbnailUrl: util.RandomThumbnailURL(),
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
	createNoteResult, err := store.TxCreateNote(context.Background(), createNoteArg)
	require.NoError(t, err)

	updateNoteParams := UpdateNoteParams{
		ID:      createNoteResult.Note.ID,
		Title:   createNoteResult.Note.Title,
		Content: createNoteResult.Note.Content,
	}
	updateWebN := 3
	updateWebIDs := make([]uuid.UUID, updateWebN)
	for i := 0; i < updateWebN; i++ {
		arg := CreateWebParams{
			UserID:       user.ID,
			Url:          util.RandomURL(),
			Title:        util.RandomName(),
			ThumbnailUrl: util.RandomThumbnailURL(),
		}
		web, err := store.CreateWeb(context.Background(), arg)
		require.NoError(t, err)
		updateWebIDs[i] = web.ID
	}
	updateNoteArg := TxUpdateNoteParams{
		UpdateNoteParams: updateNoteParams,
		WebIds:           updateWebIDs,
	}
	result, err := store.TxUpdateNote(context.Background(), updateNoteArg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.NotEmpty(t, result.Note)
	require.Equal(t, result.Note.ID, createNoteResult.Note.ID)
	require.Equal(t, result.Note.Title, updateNoteParams.Title)
	require.Equal(t, result.Note.Content, updateNoteParams.Content)
	require.Equal(t, result.Note.UserID, createNoteResult.Note.UserID)

	updatedNote, err := store.GetNote(context.Background(), createNoteResult.Note.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedNote)

	require.NotEmpty(t, result.Note)
	require.Equal(t, updatedNote.ID, createNoteResult.Note.ID)
	require.Equal(t, updatedNote.Title, updateNoteParams.Title)
	require.Equal(t, updatedNote.Content, updateNoteParams.Content)
	require.Equal(t, updatedNote.UserID, createNoteResult.Note.UserID)
}
