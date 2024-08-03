package transaction

type TransactionRepository interface {
	FindById(id string) (*Transaction, error)
	FindByAccountId(accountId string) ([]*Transaction, error)
	Create(acc *Transaction) (*Transaction, error)
	Update(acc *Transaction) (*Transaction, error)
}
