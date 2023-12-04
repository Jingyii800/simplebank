package db

// store provides all functions to execute db queries and transaction
type Store struct{
	//*Queries is an embedded pointer to a struct that contains database query methods, 
	//enabling the Store struct to directly access these methods.
	*Queries
	db *sql.DB
}

//Newstore creates a new Store
func NewStore(db *sql.DB) *Store{
	return &Store{
		db:    db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
//In database operations, it's often used to set timeouts or deadlines for queries or to cancel them if they take too long or if the operation that required the database query is no longer needed.
//It's a common practice in Go to pass a context.Context as the first parameter of a function, especially when the function involves I/O operations like database access.
//ctx represents a context.Context, used to manage request lifecycles, particularly for controlling cancellations and timeouts.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error{
	tx, err:= store.db.BeginTx(ctx, nil)
	if err!= nil{
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil{
		if rbErr := tx.Rollback(); rbErr != nil{
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}


// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json: "from_account_id"`
	ToAccountID   int64 `json: "to_account_id"`
	Amount        int64 `json:"amount"`
}
//transferTxResult is the result of the transfer transcation
type transferTxResult struct {
	Transfer     Transfer `json:"transfer"`
	FromAccount  Account  `json:"from_account"`
	ToAccount    Account  `json:"to_account"`
	FromEntry    Entry    `json:"from_entry"`
	ToEntry      Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other
// It created a transfer record, add account entries, and update accounts' balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams)(transferTxResult, error){
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}
		// add account entries
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err!= nil{
			return err
		}

		// update accounts' balance
		// 1st version:get account => update account's balance
		// 2nd version: AddAccountBalance
		// 3rd version: solve deadlock & simplify code
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.amount, arg.ToAccountID, arg.amount)
		}else{
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.amount, arg.FromAccountID, -arg.amount)
		}
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
)(account1 Account, account2 Account, err error){
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountID1,
		Amount: amount1,
	})
	if err != nil{
		return err
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountID2,
		Amount: amount2,
	})
	return 

}