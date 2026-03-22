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

func (s *BankingService) Transfer(fromID, toID string, amount float32) error {

	//Validation to check if the user are trying to transfer to itself
	if fromID == toID {
		return errors.New("Can't Transfer to urself!")
	}

	//Getting Account Block
	//For now Getting from map
	//In future from DB but it will not change the logic here!
	fromAcc, err := s.repo.GetByID(fromID)
	if err != nil {
		return err
	}

	toAcc, err := s.repo.GetByID(toID)
	if err != nil {
		return err
	}

	//Transaction level
	//
	if fromAcc.ID < toAcc.ID {
		fromAcc.Mu.Lock()
		toAcc.Mu.Lock()
	} else {
		toAcc.Mu.Lock()
		fromAcc.Mu.Lock()
	}

	defer fromAcc.Mu.Unlock()
	defer toAcc.Mu.Unlock()

	if err := fromAcc.Withdraw(amount); err != nil {
		return err
	}

	if err := toAcc.Deposit(amount); err != nil {
		fromAcc.Deposit(amount) //rollback if something went wrong
		return err
	}

	//Saving Memories
	if err := s.repo.Update(fromAcc); err != nil {
		return err
	}
	if err := s.repo.Update(toAcc); err != nil {
		return err
	}

	return nil
}
