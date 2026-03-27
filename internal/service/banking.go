package service

import (
	"BankingSystem/internal/domain"
	"errors"
)

type BankingService struct {
	repo domain.AccountRepository
}

func NewBankingService(r domain.AccountRepository) *BankingService {
	return &BankingService{repo: r}
}

func (r *BankingService) Transfer(fromID, toID string, amount float32) error {
	if fromID == toID {
		return errors.New("can't transfer to self")
	}
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	txRepo, err := r.repo.BeginTx()
	if err != nil {
		return err
	}

	committed := false
	defer func() {
		if !committed {
			_ = txRepo.RollbackTx()
		}
	}()

	fromAcc, err := txRepo.GetByID(fromID)
	if err != nil {
		return err
	}

	toAcc, err := txRepo.GetByID(toID)
	if err != nil {
		return err
	}

	if fromAcc.ID < toAcc.ID {
		fromAcc.Mu.Lock()
		toAcc.Mu.Lock()
		defer toAcc.Mu.Unlock()
		defer fromAcc.Mu.Unlock()
	} else {
		toAcc.Mu.Lock()
		fromAcc.Mu.Lock()
		defer fromAcc.Mu.Unlock()
		defer toAcc.Mu.Unlock()
	}

	if err := fromAcc.Withdraw(amount); err != nil {
		return err
	}

	if err := toAcc.Deposit(amount); err != nil {
		return err
	}

	if err := txRepo.Update(fromAcc); err != nil {
		return err
	}

	if err := txRepo.Update(toAcc); err != nil {
		return err
	}

	if err := txRepo.CommitTx(); err != nil {
		return err
	}
	committed = true

	return nil
}
