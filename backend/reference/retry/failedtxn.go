package retry

type Transaction struct {
	ID     string
	status string
}

// GetFailedTransactions simulates a database call to get failed transactions
func GetFailedTransactions() []Transaction {
	// mock
	txns := []Transaction{
		{ID: "1", status: "failed"},
		{ID: "2", status: "failed"},
		{ID: "3", status: "failed"},
	}

	return txns
}

// IsTransactionFailed simulates a database call to fetch a transaction and check if a transaction is failed
func IsTransactionFailed(txn Transaction) bool {
	// mock
	return txn.status == "failed"
}
