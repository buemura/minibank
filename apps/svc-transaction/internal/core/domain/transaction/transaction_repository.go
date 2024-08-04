package transaction

type TransactionRepository interface {
	FindById(id string) (*Transaction, error)
	FindByAccountId(in *GetTransactionListIn) (*GetTransactionListOut, error)
	Create(trx *Transaction) (*Transaction, error)
	Update(trx *Transaction) (*Transaction, error)
}
