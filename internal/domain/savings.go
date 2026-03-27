package domain

import "time"

type SavingAccount struct {
	*Account
	InterestRate float32
	LastAccrual  time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (s *SavingAccount) AccrueInterest() float32 {
	interest := s.Balance * s.InterestRate
	s.Balance += interest
	s.LastAccrual = time.Now()
	return interest
}
