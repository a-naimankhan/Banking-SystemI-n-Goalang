package app

import (
	"BankingSystem/internal/domain"
	"BankingSystem/internal/repository"
	"BankingSystem/internal/service"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func RunTestMode_1() {
	fmt.Println("----Run TestMode---")

	//Initializing Block
	repo := repository.NewInMemRepo()
	bank := service.NewBankingService(repo)

	//First Test is creating 10 acc with 1000 balance each
	//TotalCapital must be 10 000
	numAccounts := 10
	initialBalance := float32(1000)

	for i := 1; i <= numAccounts; i++ {
		repo.Create(&domain.Account{
			ID:      fmt.Sprintf("%d", i),
			Owner:   fmt.Sprintf("User%d", i),
			Balance: initialBalance,
		})
	}

	expectedTotal := initialBalance * float32(numAccounts)
	fmt.Printf("Initial Total Capital : %.2f\n", expectedTotal)

	//Starting Goroutines :
	var wg sync.WaitGroup
	numWorkers := 50
	transferPerWorker := 100

	start := time.Now()

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < transferPerWorker; i++ {
				fromIdx := rand.Intn(numAccounts) + 1
				toIdx := rand.Intn(numAccounts) + 1
				//preventing eventhough we have checker against that case
				if fromIdx == toIdx {
					continue
				}

				fromID := fmt.Sprintf("%d", fromIdx)
				toID := fmt.Sprintf("%d", toIdx)
				amount := float32(rand.Intn(10) + 1)

				bank.Transfer(fromID, toID, amount)
			}
		}(w)
	}

	wg.Wait()
	duration := time.Since(start)

	finalTotal := repo.GetTotalCapital()

	fmt.Println("\n=== Simulation Results ===")
	fmt.Printf("Final Total Capital : %.2f\n", finalTotal)
	fmt.Printf("Time Taken : %v\n ", duration)
	fmt.Printf("Difrence : %.2f\n", finalTotal-expectedTotal)

	if finalTotal == expectedTotal {
		fmt.Println("SUCCES : data integrity maintained !")
	} else {
		fmt.Println("Failed : Money Leaked ")
	}

}

func RunTestMode_2() {

	repo := repository.NewInMemRepo()

	// Создаем обычный счет
	repo.Create(&domain.Account{ID: "1", Owner: "Aibar", Balance: 1000})

	// Создаем накопительный счет (10% годовых / ставка 0.1)
	savings := &domain.SavingAccount{
		Account:      &domain.Account{ID: "2", Owner: "Investor", Balance: 5000},
		InterestRate: 0.1,
	}
	repo.Create(savings)

	// Запускаем фоновый воркер
	// Он будет работать в отдельной горутине
	service.StartInterestWorker(repo, 10*time.Second)

	// Дальше твоя обычная логика меню или переводов...
	fmt.Println("Bank is running. Wait 10s to see interest accrual.")
	select {} // Чтобы main не закрылся сразу
}
