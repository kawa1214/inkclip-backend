package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/bookmark-manager/bookmark-manager/util"
	"github.com/stretchr/testify/require"
)

func TestCreateNote(t *testing.T) {
	user := createRandomUser(t)
	createRandomNote(t, user)
}

func TestDeleteNote(t *testing.T) {
	user := createRandomUser(t)
	note := createRandomNote(t, user)

	err := testQueries.DeleteNote(context.Background(), note.ID)
	require.NoError(t, err)

	note2, err := testQueries.GetNote(context.Background(), note.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, note2)
}

func TestGetNote(t *testing.T) {
	user := createRandomUser(t)
	note := createRandomNote(t, user)
	gotNote, err := testQueries.GetNote(context.Background(), note.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gotNote)

	require.Equal(t, note.ID, gotNote.ID)
	require.Equal(t, note.Title, gotNote.Title)
	require.Equal(t, note.Content, gotNote.Content)
	require.Equal(t, note.UserID, gotNote.UserID)
	require.Equal(t, note.CreatedAt, gotNote.CreatedAt)
}

func TestUpdateNote(t *testing.T) {
	user := createRandomUser(t)
	note := createRandomNote(t, user)

	arg := UpdateNoteParams{
		ID:      note.ID,
		Title:   util.RandomString(12),
		Content: util.RandomString(12),
	}
	updatedNote, err := testQueries.UpdateNote(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedNote)

	require.Equal(t, arg.ID, updatedNote.ID)
	require.Equal(t, arg.Title, updatedNote.Title)
	require.Equal(t, arg.Content, updatedNote.Content)
	require.Equal(t, note.UserID, updatedNote.UserID)
	require.Equal(t, note.CreatedAt, updatedNote.CreatedAt)
}

func TestListNotesByUserId(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		createRandomNote(t, user)
	}

	arg := ListNotesByUserIdParams{
		UserID: user.ID,
		Limit:  5,
		Offset: 5,
	}
	notes, err := testQueries.ListNotesByUserId(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, notes)

	for _, note := range notes {
		require.NotEmpty(t, note)
	}
}

func createRandomNote(t *testing.T, user User) Note {
	arg := CreateNoteParams{
		Title:   util.RandomString(6),
		Content: util.RandomString(6),
		UserID:  user.ID,
	}
	note, err := testQueries.CreateNote(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, note)

	require.Equal(t, arg.Title, note.Title)
	require.Equal(t, arg.Content, note.Content)
	require.Equal(t, arg.UserID, note.UserID)

	require.NotZero(t, note.CreatedAt)

	return note
}
