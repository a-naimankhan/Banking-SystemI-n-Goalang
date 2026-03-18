package domain

import (
	"errors"
	"sync"
)

// THE main things about Account
type Account struct {
	Mu      sync.Mutex
	ID      string
	Owner   string
	Balance float32
}

type AccountRepository interface {
	GetByID(id string) (*Account, error)
	Update(anyAccount interface{}) error
	Create(anyAccount interface{}) error
}

func (a *Account) Deposit(amount float32) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	a.Balance += amount
	return nil
}

func (a *Account) Withdraw(amount float32) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if amount > a.Balance {
		return errors.New("amount must be less than or equal to balance")
	}

	a.Balance -= amount
	return nil
}
