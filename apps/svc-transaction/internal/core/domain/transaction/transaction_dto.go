package transaction

type CreateTransactionIn struct {
	AccountID            string          `json:"account_id"`
	DestinationAccountID *string         `json:"destination_account_id,omitempty"`
	Amount               int             `json:"amount"`
	TransactionType      TransactionType `json:"transaction_type"`
}
