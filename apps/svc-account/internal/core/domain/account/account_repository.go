package account

type AccountRepository interface {
	FindById(id string) (*Account, error)
	FindByOwnerDocument(document string) (*Account, error)
	Create(acc *Account) (*Account, error)
	Update(acc *Account) (*Account, error)
}
