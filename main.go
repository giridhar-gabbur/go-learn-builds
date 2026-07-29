package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"log"
)

type BankAccount struct {
	owner string
	balance float64
	history []string
	closed bool
	mu sync.RWMutex
}


var (
	ErrInsufficientFunds = errors.New("Insufficient Funds in account")
	ErrInvalidAmount = errors.New("Invalid amount requested")
	ErrAccountClosed = errors.New("Sorry! Account is closed")
)

type TransactionError struct {
	Type string
	Amount float64
	Reason error
}

func (e TransactionError) Error() string{
	return fmt.Sprintf("Trasaction error of type: %s | Amount: %0.2f | Error: %v", e.Type, e.Amount, e.Reason)
}

func (e TransactionError) Unwrap() error {
    return e.Reason
}

func NewBankAccount (name string, amount float64) (*BankAccount, error){
	if name == "" {
		return nil, fmt.Errorf("Bank account name cannot be empty")
	}
	if amount < 0 {
		return nil, fmt.Errorf("Balance cannot be negative")
	}
	acc := &BankAccount{owner: name, balance: amount}
	acc.history = append(acc.history, fmt.Sprintf("Account opened. Initial Balance: %.2f",amount))
	return acc, nil
}

func FindAccount (owner string, accounts []*BankAccount) *BankAccount{
	for _,acc := range accounts {
		if acc.owner == owner {
			return acc
		}
	}
	return nil
}

func (b *BankAccount) DepositAmount(amount float64) (e error){
	
	if amount <= 0 {
		return TransactionError{Type: "Deposit", Amount: amount, Reason: ErrInvalidAmount}
	}
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed{
		return fmt.Errorf("deposit: %w", ErrAccountClosed)
	}
	
	b.balance += amount
	b.history = append(b.history, fmt.Sprintf("Credit: +%.2f | Balance: %.2f",amount, b.balance))
	return nil
}

func (b *BankAccount) Withdraw(amount float64) (e error){
	if amount <= 0 {
		return TransactionError{Type: "Withdraw", Amount: amount, Reason: ErrInvalidAmount}
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed{
		return ErrAccountClosed
	}
	if amount > b.balance{
		return TransactionError{
			Type:   "Withdraw",
			Amount: amount,
			Reason: ErrInsufficientFunds,
    	}
	}
	
	b.balance -= amount
	b.history = append(b.history, fmt.Sprintf("Debit: -%.2f | Balance: %.2f",amount, b.balance))
	return nil
}

func (b *BankAccount) CurrBalance() (float64){
	b.mu.RLock()
    defer b.mu.RUnlock()
	return b.balance
}

func (b *BankAccount) Transfer(amount float64, target *BankAccount) error {
	if err := b.Withdraw(amount); err != nil {
		return fmt.Errorf("Transfer Failed: %w", err)
	}
	if err := target.DepositAmount(amount); err != nil{
		b.DepositAmount(amount)
		return fmt.Errorf("Transfer Failed: %w", err)
	}
	return nil
}

func TransferBetween(from *BankAccount, to *BankAccount, amount float64) error {
	if from == to {
		return fmt.Errorf("Transfer failed. Accounts must be distinct!")
	}
	err := from.Withdraw(amount)
	if err != nil {
		return fmt.Errorf("Transfer Failed: %w", err)
	}
	if err := to.DepositAmount(amount); err != nil{
		from.DepositAmount(amount)
		return fmt.Errorf("Transfer Failed: %w", err)
	}
	return nil
}

func ProcessTransactions(b *BankAccount, transactions []float64, wg *sync.WaitGroup) {
	for _,amount := range transactions {
		wg.Add(1)
		go func(amt float64){
			defer wg.Done()
			if amt > 0 {
				b.DepositAmount(amt)
			} else if amt< 0 {
				b.Withdraw(-amt)
			}
		} (amount)
	}
}


func (a *BankAccount) String() (string){
	return fmt.Sprintf("Name: %s: Balance: %.2f", a.owner, a.balance)
}

func (b *BankAccount) PrintHistory(){
	b.mu.RLock()
    defer b.mu.RUnlock()
	fmt.Printf("---------Transaction History: %s----------\n",b.owner)
	for _,entry := range b.history {
		fmt.Println(" ", entry)
	}
}

type Transaction struct {
	ID string
	Type string
	Amount float64
}

type TxResult struct {
	TxID string
	Success bool 
	Err error
}

func txWorker (id int, transactions <- chan Transaction, results chan <- TxResult, b *BankAccount, wg *sync.WaitGroup){
	defer wg.Done()
	for tx := range transactions {
		fmt.Printf("[Worker %d] processing the transaction: %s of type %s of %.2f\n", id, tx.ID, tx.Type, tx.Amount)
		var err error
		if tx.Type == "Deposit"{
			err = b.DepositAmount(tx.Amount)
		} else if tx.Type == "Withdrawal"{
			err = b.Withdraw(tx.Amount)
		}

		results <- TxResult{
			TxID: tx.ID,
			Success: err == nil,
			Err: err,
		}
	}
}

func monitorAccount(b *BankAccount, alerts chan<- string, done <-chan bool) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()  
	for {
		select {
		case <- done:
			fmt.Println("Monitor stopped!")
			close(alerts)
			return
		case <- ticker.C:
			if b.CurrBalance() < 1000 {
				alerts <- "Low Balance!"
			}
		}
	}
}

func main() {
	b, err := NewBankAccount("Giridhar", 1000)
	if err != nil {
		log.Fatal("Account creation failed:", err)
	}
	var wg sync.WaitGroup
	var wgMonitor sync.WaitGroup

	transactions := make(chan Transaction, 20)
	results:= make(chan TxResult, 20)

	alerts := make(chan string, 20)
	done := make(chan bool)

	wgMonitor.Add(1)
	go func() {
		defer wgMonitor.Done()
		monitorAccount(b, alerts, done)
	}()

	for w := 1; w < 3; w++ {
		wg.Add(1)
		go txWorker(w, transactions, results, b, &wg)
	}

	txs := []Transaction{
		{ID: "1A", Type: "Deposit", Amount: 1300},
		{ID: "1B", Type: "Withdrawal", Amount: 1300},
		{ID: "1C", Type: "Deposit", Amount: 53},
		{ID: "1D", Type: "Withdrawal", Amount: 100.23},
		{ID: "1E", Type: "Deposit", Amount: 13},
		{ID: "1F", Type: "Deposit", Amount: 130},
	}

	for _,tx := range txs {
		transactions <- tx
	}
	close(transactions)

	go func() {
		wg.Wait()
		close(results)
	}()

	successes, failures := 0,0
	for res := range results {
		if res.Success {
			fmt.Printf("Transation %s: Approved\n", res.TxID)
			successes++
		} else {
			fmt.Printf("Transaction %s: Declined -> %v\n", res.TxID, res.Err)
			failures++
		}
	}
	fmt.Printf("\n=== Summary ===\n")
    fmt.Printf("✅ Succeeded: %d\n", successes)
    fmt.Printf("❌ Failed:    %d\n", failures)
	fmt.Printf("\nFinal Audited Account Balance: $%.2f\n", b.CurrBalance())

	done <- true
	wgMonitor.Wait()
	for alert := range alerts {
		fmt.Println(alert)
	}
}
