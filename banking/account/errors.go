package account

import (
    "errors"
    "fmt"       // ← added!
)

var (
    ErrInsufficientFunds = errors.New("insufficient funds")
    ErrInvalidAmount     = errors.New("invalid amount")
    ErrAccountClosed     = errors.New("account is closed")
)

type TransactionError struct {
    Type   string
    Amount float64
    Reason error
}

func (e TransactionError) Error() string {
    return fmt.Sprintf("transaction '%s' | amount: %.2f | reason: %v",
        e.Type, e.Amount, e.Reason)
}

func (e TransactionError) Unwrap() error {
    return e.Reason
}