package transaction

// Transaction represents a single banking transaction
type Transaction struct {
    ID     string
    Type   string      // "Deposit" or "Withdrawal"
    Amount float64
}

// Result represents the outcome of a transaction
type Result struct {
    TxID    string
    Success bool
    Err     error
}
