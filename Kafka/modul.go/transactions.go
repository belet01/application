package modul

import "time"

type Transactions struct {
	Id              int        `json:"id"`
	AccountId       int        `json:"account_id,omitempty"`
	Amount          string     `json:"amount"`
	TransactionType string     `json:"transaction_type"`
	CreatedAt       time.Time  `json:"created_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

type TransactionREsponse struct {
	YourId          int     `json:"your_id"`
	AccountId       int     `json:"account_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
}
type TransactionResponse struct {
	AccountId   int            `json:"account_id"`
	Transaction []Transactions `json:"transactions"`
}
