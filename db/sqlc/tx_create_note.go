package db

import "context"

type TxCreateNoteParams struct {
	createNoteParams    CreateNoteParams
	CreateWebParamsList []CreateWebParams
}

type TxCreateNoteResult struct {
	note Note
	webs []Web
}

func (store *SQLStore) TxCreateNote(ctx context.Context, arg TxCreateNoteParams) (TxCreateNoteResult, error) {
	var result TxCreateNoteResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.note, err = q.CreateNote(ctx, arg.createNoteParams)
		if err != nil {
			return err
		}

		result.webs = make([]Web, len(arg.CreateWebParamsList))
		for i := 0; i < len(arg.CreateWebParamsList); i++ {
			result.webs[i], err = q.CreateWeb(ctx, arg.CreateWebParamsList[i])
			if err != nil {
				return err
			}
			_, err = q.CreateNoteWeb(ctx, CreateNoteWebParams{
				NoteID: result.note.ID,
				WebID:  result.webs[i].ID,
			})
			if err != nil {
				return err
			}
		}

		return err
	})

	return result, err
}
