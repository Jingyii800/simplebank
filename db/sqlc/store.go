package db

import (
	"context"
	"database/sql"
	"fmt"
)

// for mock: Store
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

// store provides all functions to execute db queries and transaction
type SQLStore struct {
	//*Queries is an embedded pointer to a struct that contains database query methods,
	//enabling the Store struct to directly access these methods.
	*Queries
	db *sql.DB
}

// Newstore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
// In database operations, it's often used to set timeouts or deadlines for queries or to cancel them if they take too long or if the operation that required the database query is no longer needed.
// It's a common practice in Go to pass a context.Context as the first parameter of a function, especially when the function involves I/O operations like database access.
// ctx represents a context.Context, used to manage request lifecycles, particularly for controlling cancellations and timeouts.
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
