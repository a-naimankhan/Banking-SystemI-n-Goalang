package service

import (
	"BankingSystem/internal/domain"
	"fmt"
	"time"
)

type SavingsRepository interface {
	GetAllSavings() []*domain.SavingAccount
}

type InterestWorker struct {
	repo     SavingsRepository
	interval time.Duration
	stopChan chan struct{}
}

func NewInterestWorker(r SavingsRepository, interval time.Duration) *InterestWorker {
	return &InterestWorker{
		repo:     r,
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

func (w *InterestWorker) Start() {
	ticker := time.NewTicker(w.interval)

	go func() {
		fmt.Printf("[InterestWorker] Started with interlval %s\n", w.interval)
		for {
			select {
			case <-ticker.C:
				w.process()
			case <-w.stopChan:
				ticker.Stop()
				fmt.Println("[Interest worker] Stopped")
				return
			}
		}
	}()
}

func (w *InterestWorker) process() {
	accounts := w.repo.GetAllSavings()
	if len(accounts) == 0 {
		return
	}

	fmt.Printf("[InterestWorker] Checking %d savings accounts...\n", len(accounts))
	for _, acc := range accounts {
		acc.Mu.Lock()
		interest := acc.AccrueInterest()
		acc.Mu.Unlock()

		if interest > 0 {
			fmt.Printf("[InterestWorker] Accrued %.2f to Account %s (Owner: %s)\n",
				interest, acc.ID, acc.Owner)
		}
	}
}

func (w *InterestWorker) Stop() {
	close(w.stopChan)
}
