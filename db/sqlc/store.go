package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store providers all functions to execute db queries and transactions
type Store interface {
	Querier
	TxCreateNote(ctx context.Context, arg TxCreateNoteParams) (TxCreateNoteResult, error)
	TxDeleteNote(ctx context.Context, arg TxDeleteNoteParams) error
	TxUpdateNote(ctx context.Context, arg TxUpdateNoteParams) (TxUpdateNoteResult, error)
}

// SQLStore providers all functions to execute SQL queries and transactions
type SQLStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
