package transaction

type TransactionRepository interface {
	FindById(id string) (*Transaction, error)
	FindByAccountId(in *GetTransactionListIn) (*GetTransactionListOut, error)
	Create(acc *Transaction) (*Transaction, error)
	Update(acc *Transaction) (*Transaction, error)
}
