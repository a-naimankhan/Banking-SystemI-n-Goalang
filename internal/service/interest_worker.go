package service

import (
	"BankingSystem/internal/domain"
	"fmt"
	"time"
)

type SavingsReposity interface {
	GetAllSavings() []*domain.SavingAccount
}

func StartInterestWorker(repo SavingsReposity, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			fmt.Println("[Worker] Checking savings accounts for interest...")
			accounts := repo.GetAllSavings()
			for _, acc := range accounts {
				amount := acc.AccrueInterest()
				fmt.Printf("[Worker] Accured %.2f to Account %s (Owner : %s)\n", amount, acc.ID, acc.Owner)
			}
		}
	}()
}
