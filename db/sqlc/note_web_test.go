package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateNoteWeb(t *testing.T) {
	user := createRandomUser(t)
	note := createRandomNote(t, user)
	web := createRandomWeb(t, user)
	createRandomNoteWeb(t, note, web)
}

func createRandomNoteWeb(t *testing.T, note Note, web Web) NoteWeb {
	arg := CreateNoteWebParams{
		NoteID: note.ID,
		WebID:  web.ID,
	}

	noteWeb, err := testQueries.CreateNoteWeb(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, noteWeb)

	require.Equal(t, arg.NoteID, noteWeb.NoteID)
	require.Equal(t, arg.WebID, noteWeb.WebID)

	return noteWeb
}

func TestDeleteNoteWeb(t *testing.T) {
	user := createRandomUser(t)
	note := createRandomNote(t, user)
	web := createRandomWeb(t, user)
	noteWeb := createRandomNoteWeb(t, note, web)
	require.NotEmpty(t, noteWeb)

	deleteArg := DeleteNoteWebParams{
		NoteID: note.ID,
		WebID:  web.ID,
	}
	err := testQueries.DeleteNoteWeb(context.Background(), deleteArg)
	require.NoError(t, err)

	getArg := GetNoteWebParams{
		NoteID: note.ID,
		WebID:  web.ID,
	}
	_, err = testQueries.GetNoteWeb(context.Background(), getArg)
	require.Error(t, err)
}

func TestListNoteWebByNoteId(t *testing.T) {
	user := createRandomUser(t)
	note := createRandomNote(t, user)
	n := 5
	for i := 0; i < n; i++ {
		web := createRandomWeb(t, user)
		createRandomNoteWeb(t, note, web)
	}

	noteWebs, err := testQueries.ListNoteWebsByNoteId(context.Background(), note.ID)
	require.NoError(t, err)
	require.NotEmpty(t, noteWebs)

	require.Equal(t, n, len(noteWebs))

	for _, noteWeb := range noteWebs {
		require.NotEmpty(t, noteWeb)
	}
}
