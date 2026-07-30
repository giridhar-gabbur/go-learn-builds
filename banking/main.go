package main

import (
	"fmt"
	"log"

	"banking/account"
	"banking/transaction"
)

func main() {
	b, err := account.NewBankAccount("Giridhar", 1000)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(b)

	fmt.Println("\n====== Processing REQ-001 ======")
	transaction.ProcessRequest(b, "REQ-001", []transaction.Transaction{
		{ID: "1A", Type: "Deposit",    Amount: 1300},
		{ID: "1B", Type: "Withdrawal", Amount: 1300},
		{ID: "1C", Type: "Deposit",    Amount: 53},
	})

	fmt.Println("\n====== Processing REQ-002 ======")
	transaction.ProcessRequest(b, "REQ-002", []transaction.Transaction{
		{ID: "2A", Type: "Deposit",    Amount: 500},
		{ID: "2B", Type: "Withdrawal", Amount: 200},
		{ID: "2C", Type: "Deposit",    Amount: 100},
	})

	fmt.Printf("\nFinal Balance: $%.2f\n", b.Balance())
}
