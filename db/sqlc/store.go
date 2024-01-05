package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
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
	connPool *pgxpool.Pool
}

// Newstore creates a new Store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
