package account

import (
	"time"
)

type Account struct {
	ID            string
	Balance       int
	OwnerName     string
	OwnerDocument string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
