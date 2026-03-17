package repository

import (
	"BankingSystem/internal/domain"
	"errors"
	"sync"
)

type InMemoRepo struct {
	mu      sync.RWMutex
	storage map[string]*domain.Account
}

func NewInMemRepo() *InMemoRepo {
	return &InMemoRepo{
		storage: make(map[string]*domain.Account),
	}
}

func (r *InMemoRepo) Create(acc *domain.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.storage[acc.ID] = acc
	return nil
}
func (r *InMemoRepo) GetByID(id string) (*domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	acc, ok := r.storage[id]
	if !ok {
		return nil, errors.New("Account doesn't exists")
	}
	return acc, nil
}

func (r *InMemoRepo) Update(acc *domain.Account) error {
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
		total += acc.Balance
	}

	return total
}
