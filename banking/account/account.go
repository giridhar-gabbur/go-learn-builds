package account

import (
    "fmt"
    "sync"
)

type BankAccount struct {
    mu      sync.RWMutex
    Owner   string
    balance float64
    history []string
    closed  bool
}

func NewBankAccount(name string, amount float64) (*BankAccount, error) {
    if name == "" {
        return nil, fmt.Errorf("owner name cannot be empty")
    }
    if amount < 0 {
        return nil, fmt.Errorf("initial balance cannot be negative")
    }
    acc := &BankAccount{Owner: name, balance: amount}
    acc.history = append(acc.history,
        fmt.Sprintf("Account opened. Initial Balance: %.2f", amount))
    return acc, nil
}

func (b *BankAccount) Deposit(amount float64) error {
    if amount <= 0 {
        return TransactionError{
            Type: "Deposit", Amount: amount, Reason: ErrInvalidAmount,
        }
    }
    b.mu.Lock()
    defer b.mu.Unlock()
    if b.closed {
        return fmt.Errorf("deposit: %w", ErrAccountClosed)
    }
    b.balance += amount
    b.history = append(b.history,
        fmt.Sprintf("Credit: +%.2f | Balance: %.2f", amount, b.balance))
    return nil
}

func (b *BankAccount) Withdraw(amount float64) error {
    if amount <= 0 {
        return TransactionError{
            Type: "Withdraw", Amount: amount, Reason: ErrInvalidAmount,
        }
    }
    b.mu.Lock()
    defer b.mu.Unlock()
    if b.closed {
        return fmt.Errorf("withdraw: %w", ErrAccountClosed)
    }
    if amount > b.balance {
        return TransactionError{
            Type: "Withdraw", Amount: amount, Reason: ErrInsufficientFunds,
        }
    }
    b.balance -= amount
    b.history = append(b.history,
        fmt.Sprintf("Debit: -%.2f | Balance: %.2f", amount, b.balance))
    return nil
}

func (b *BankAccount) Balance() float64 {
    b.mu.RLock()
    defer b.mu.RUnlock()
    return b.balance
}

func (b *BankAccount) Close() {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.closed = true
}

func (b *BankAccount) String() string {
    return fmt.Sprintf("Account[%s] Balance: $%.2f", b.Owner, b.Balance())
}