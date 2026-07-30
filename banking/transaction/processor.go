package transaction

import (
	"context"
	"fmt"
	"sync"
	"time"

	"banking/account"
	"banking/middleware"
)

func ProcessRequest(
	b         *account.BankAccount,
	requestID string,
	txs       []Transaction,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ctx = context.WithValue(ctx, middleware.ReqIDKey, requestID)

	var wg sync.WaitGroup
	transactions := make(chan Transaction, len(txs))
	results      := make(chan Result, len(txs))

	for w := 1; w <= 5; w++ {
		wg.Add(1)
		go Worker(ctx, w, transactions, results, b, &wg)
	}

	for _, tx := range txs {
		transactions <- tx
	}
	close(transactions)

	go func() {
		wg.Wait()
		close(results)
	}()

	successes, failures := 0, 0
	for res := range results {
		if res.Success {
			fmt.Printf("[%s] ✅ %s: Approved\n", requestID, res.TxID)
			successes++
		} else {
			fmt.Printf("[%s] ❌ %s: Declined → %v\n", requestID, res.TxID, res.Err)
			failures++
		}
	}

	fmt.Printf("\n[%s] Summary: ✅ %d | ❌ %d\n", requestID, successes, failures)
	fmt.Printf("[%s] Balance: $%.2f\n", requestID, b.Balance())

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Printf("[%s] ⏰ Batch timed out! %d/%d completed\n",
			requestID, successes, len(txs))
	}
}