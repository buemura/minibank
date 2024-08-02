package account

type AccountRepository interface {
	FindById(id string) (*Account, error)
}
