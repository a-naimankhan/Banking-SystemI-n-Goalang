package repository

import (
	"BankingSystem/internal/domain"
	"errors"
	"sync"
)

type InMemoRepo struct {
	mu      sync.RWMutex
	storage map[string]interface{}
}

func NewInMemRepo() *InMemoRepo {
	return &InMemoRepo{
		storage: make(map[string]interface{}),
	}
}

func (r *InMemoRepo) Create(acc interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var id string

	switch a := acc.(type) {
	case *domain.Account:
		id = a.ID
	case *domain.SavingAccount:
		id = a.ID
	}
	r.storage[id] = acc
	return nil
}

func (r *InMemoRepo) GetByID(id string) (*domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	acc, ok := r.storage[id]
	if !ok {
		return nil, errors.New("Account doesn't exists")
	}
	switch a := acc.(type) {
	case *domain.Account:
		return a, nil
	case *domain.SavingAccount:
		return a.Account, nil
	default:
		return nil, errors.New("unkown account type")
	}
}

func (r *InMemoRepo) GetAllSavings() []*domain.SavingAccount {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var savings []*domain.SavingAccount
	for _, val := range r.storage {
		if s, ok := val.(*domain.SavingAccount); ok {
			savings = append(savings, s)
		}
	}
	return savings
}

func (r *InMemoRepo) Update(acc interface{}) error {
	//since now we didn't add DB and everything already saving in memory by it's own methods for now we leave this part as free
	//But to satisfy the interface we have to add this method !
	return nil

}

func (r *InMemoRepo) Delete(id string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.storage[id] == nil {
		return errors.New("Account doesn't exists")
	}

	r.storage[id] = nil
	return nil
}

func (r *InMemoRepo) GetTotalCapital() float32 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var total float32

	for _, acc := range r.storage {
		switch a := acc.(type) {
		case *domain.Account:
			total += a.Balance
		case *domain.SavingAccount:
			total += a.Balance
		}
	}

	return total

}
