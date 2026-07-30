package transaction

import (
	"context"
	"fmt"
	"time"

	"banking/account"
	"banking/middleware"
)

func processTransaction(ctx context.Context, b *account.BankAccount, tx Transaction) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	subCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		var err error
		switch tx.Type {
		case "Deposit":
			err = b.Deposit(tx.Amount)
		case "Withdrawal":
			err = b.Withdraw(tx.Amount)
		default:
			err = fmt.Errorf("unknown transaction type: %s", tx.Type)
		}
		done <- err
	}()

	select {
	case err := <-done:
		return err
	case <-subCtx.Done():
		return fmt.Errorf("transaction timed out: %w", subCtx.Err())
	}
}

func Worker(
	ctx          context.Context,
	id           int,
	transactions <-chan Transaction,
	results      chan<- Result,
	b            *account.BankAccount,
	wg           interface{ Done() },
) {
	defer wg.Done()

	reqID, ok := ctx.Value(middleware.ReqIDKey).(string)
	if !ok {
		reqID = "UNKNOWN"
	}

	for tx := range transactions {
		select {
		case <-ctx.Done():
			results <- Result{
				TxID:    tx.ID,
				Success: false,
				Err:     ctx.Err(),
			}
			continue
		default:
		}

		fmt.Printf("[Worker %d][%s] processing: %s | %s | $%.2f\n",
			id, reqID, tx.ID, tx.Type, tx.Amount)

		err := processTransaction(ctx, b, tx)
		results <- Result{
			TxID:    tx.ID,
			Success: err == nil,
			Err:     err,
		}
	}
}
