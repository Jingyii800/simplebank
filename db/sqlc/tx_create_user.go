package db

import "context"

// CreateUserTxParams contains the input parameters of the transfer transaction
type CreateUserTxParams struct {
	CreateUserParams
	// callback func after user in inserted ,
	// error is to decide whether to commit or rollback the transaction
	AfterCreate func(user User) error
}

// CreateUserTxResult is the result of the transfer transcation
type CreateUserTxResult struct {
	User User
}

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return nil // rollback
		}

		return arg.AfterCreate(result.User)

	})
	return result, err
}
