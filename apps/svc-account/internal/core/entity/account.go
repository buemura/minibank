package entity

import (
	"crypto/rand"
	"errors"
	"time"

	"github.com/buemura/minibank/svc-account/internal/core/dto"
	"github.com/lucsky/cuid"
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

func NewAccount(in dto.CreateAccountIn) (*Account, error) {
	if err := validate(in); err != nil {
		return nil, err
	}

	cuid, err := cuid.NewCrypto(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &Account{
		ID:            cuid,
		Balance:       0,
		OwnerName:     in.OwnerName,
		OwnerDocument: in.OwnerDocument,
		Status:        "ACTIVE",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

func validate(in dto.CreateAccountIn) error {
	if in.OwnerName == "" {
		return errors.New("invalid owner name")
	}

	if in.OwnerDocument == "" {
		return errors.New("invalid owner document")
	}

	return nil
}
