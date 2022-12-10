package db

import (
	"context"

	"github.com/google/uuid"
)

type TxUpdateNoteParams struct {
	UpdateNoteParams UpdateNoteParams
	WebIds           []uuid.UUID
}

type TxUpdateNoteResult struct {
	Note Note
	Webs []Web
}

func (store *SQLStore) TxUpdateNote(ctx context.Context, arg TxUpdateNoteParams) (TxUpdateNoteResult, error) {
	var result TxUpdateNoteResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Note, err = q.UpdateNote(ctx, arg.UpdateNoteParams)
		if err != nil {
			return err
		}

		currentNoteWebs, err := q.ListNoteWebsByNoteId(ctx, result.Note.ID)
		if err != nil {
			return err
		}
		for i := 0; i < len(currentNoteWebs); i++ {
			err := q.DeleteNoteWeb(ctx, DeleteNoteWebParams{
				NoteID: currentNoteWebs[i].NoteID,
				WebID:  currentNoteWebs[i].WebID,
			})
			if err != nil {
				return err
			}
		}

		result.Webs = make([]Web, len(arg.WebIds))
		for i := 0; i < len(arg.WebIds); i++ {
			webID := arg.WebIds[i]
			_, err := q.CreateNoteWeb(ctx, CreateNoteWebParams{
				NoteID: result.Note.ID,
				WebID:  webID,
			})
			if err != nil {
				return err
			}
			web, err := q.GetWeb(ctx, webID)
			result.Webs[i] = web
			if err != nil {
				return err
			}
		}

		return err
	})

	return result, err
}
