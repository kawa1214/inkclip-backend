package db

import (
	"context"

	"github.com/google/uuid"
)

type TxDeleteNoteParams struct {
	NoteID uuid.UUID
}

func (store *SQLStore) TxDeleteNote(ctx context.Context, arg TxDeleteNoteParams) error {

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		noteWebs, err := q.ListNoteWebsByNoteId(ctx, arg.NoteID)
		if err != nil {
			return err
		}
		for i := 0; i < len(noteWebs); i++ {
			noteWeb := noteWebs[i]
			deleteNoteWebArg := DeleteNoteWebParams{
				NoteID: noteWeb.WebID,
				WebID:  noteWeb.WebID,
			}
			err = q.DeleteNoteWeb(ctx, deleteNoteWebArg)
			if err != nil {
				return err
			}
		}

		err = q.DeleteNote(ctx, arg.NoteID)
		if err != nil {
			return err
		}

		return err
	})

	return err
}
