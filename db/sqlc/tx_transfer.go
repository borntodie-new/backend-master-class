package sqlc

import (
	"context"
	"fmt"
)

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer      Transfer `json:"transfer"`
	FromAccountId int64    `json:"from_account_id"`
	ToAccountId   int64    `json:"to_account_id"`
	FromEntry     Entry    `json:"from_entry"`
	ToEntry       Entry    `json:"to_entry"`
	FromAccount   Account  `json:"from_account"`
	ToAccount     Account  `json:"to_account"`
}

// TransferTx performs money transfer from one account to the other
// It creates a transfer records, add account entries, and update accounts' balance within a single database transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		txName := ctx.Value(txKey)

		var err error
		// perform create a transfer information
		fmt.Println(txName, "create transfer")
		fmt.Printf("account%d transfer %d$ to account%d\n", arg.FromAccountId, arg.Amount, arg.ToAccountId)
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})

		// perform create a entry about fromAccountId
		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})

		// perform create a entry about toAccountId
		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})

		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountId, -arg.Amount, arg.ToAccountId, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)
		}

		if err != nil {
			return err
		}
		//// update account's balance
		//// get account information
		//fmt.Println(txName, "get account 2")
		//account1, err := store.GetAccountForUpdate(context.Background(), arg.FromAccountId)
		//if err != nil {
		//	return err
		//}
		//// update account information
		//fmt.Println(txName, "update account 1")
		//result.FromAccount, err = store.UpdateAccount(context.Background(), UpdateAccountParams{
		//	ID:      account1.ID,
		//	Balance: account1.Balance - arg.Amount,
		//})
		//if err != nil {
		//	return err
		//}
		//// get account information
		//fmt.Println(txName, "get account 2")
		//account2, err := store.GetAccountForUpdate(context.Background(), arg.ToAccountId)
		//if err != nil {
		//	return err
		//}
		//// update account information
		//fmt.Println(txName, "update account 2")
		//result.ToAccount, err = store.UpdateAccount(context.Background(), UpdateAccountParams{
		//	ID:      account2.ID,
		//	Balance: account2.Balance + arg.Amount,
		//})
		//if err != nil {
		//	return err
		//}
		return nil
	})
	return result, err
}

func addMoney(ctx context.Context, q *Queries, accountID1, amount1, accountID2, amount2 int64) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID:     accountID2,
	})
	if err != nil {
		return
	}
	return
}
