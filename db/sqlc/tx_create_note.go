package db

import (
	"context"

	"github.com/google/uuid"
)

type TxCreateNoteParams struct {
	CreateNoteParams CreateNoteParams
	WebIds           []uuid.UUID
}

type TxCreateNoteResult struct {
	Note Note
	Webs []Web
}

func (store *SQLStore) TxCreateNote(ctx context.Context, arg TxCreateNoteParams) (TxCreateNoteResult, error) {
	var result TxCreateNoteResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Note, err = q.CreateNote(ctx, arg.CreateNoteParams)
		if err != nil {
			return err
		}

		result.Webs = make([]Web, len(arg.WebIds))
		for i := 0; i < len(arg.WebIds); i++ {
			result.Webs[i], err = q.GetWeb(ctx, arg.WebIds[i])
			if err != nil {
				return err
			}
			_, err = q.CreateNoteWeb(ctx, CreateNoteWebParams{
				NoteID: result.Note.ID,
				WebID:  result.Webs[i].ID,
			})
			if err != nil {
				return err
			}
		}

		return err
	})

	return result, err
}
