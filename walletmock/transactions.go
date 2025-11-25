package main

type Transaction struct {
	RefID        string `json:"refId"`
	PlayerID     string `json:"playerId"`
	Type         string `json:"type"`
	Amount       int64  `json:"amount"`
	Currency     string `json:"currency"`
	BalanceAfter int64  `json:"balanceAfter"`
}

type TransactionLog struct {
	Items []Transaction
}

func NewTransactionLog() *TransactionLog {
	return &TransactionLog{Items: []Transaction{}}
}

func (l *TransactionLog) Add(tx Transaction) {
	l.Items = append(l.Items, tx)
}
