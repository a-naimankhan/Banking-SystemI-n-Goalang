package app

import (
	"BankingSystem/internal/domain"
	"BankingSystem/internal/repository"
	"BankingSystem/internal/repository/postgres"
	"BankingSystem/internal/service"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	overAllTimeStart := time.Now()

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < transferPerWorker; i++ {
				fromIdx := rand.Intn(numAccounts) + 1
				toIdx := rand.Intn(numAccounts) + 1
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

	overAllTimeDuration := time.Since(overAllTimeStart)

	fmt.Println("\n=== Overall Time ===")
	fmt.Printf("Overall Time Capital : %s\n", overAllTimeDuration)
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
	//service.StartInterestWorker(repo, 10*time.Second)

	// Дальше твоя обычная логика меню или переводов...
	fmt.Println("Bank is running. Wait 10s to see interest accrual.")
	select {} // Чтобы main не закрылся сразу
}

func RunTestMode_3() {
	repo := repository.NewInMemRepo()

	// 1. Создаем аккаунты (обычный и накопительный)
	repo.Create(&domain.Account{ID: "1", Owner: "Aibar", Balance: 1000})
	repo.Create(&domain.SavingAccount{
		Account:      &domain.Account{ID: "2", Owner: "Investor", Balance: 5000},
		InterestRate: 0.01, // 1% за каждый тик
	})

	// 2. Инициализируем и запускаем воркер (каждые 10 секунд для теста)
	worker := service.NewInterestWorker(repo, 10*time.Second)
	worker.Start()

	// 3. Твоя логика (меню или просто ожидание)
	fmt.Println("Bank is running. Press Enter to exit...")
	fmt.Scanln() // Ждем нажатия Enter, чтобы приложение не закрылось

	worker.Stop() // Мягко останавливаем воркер перед выходом
}

func RunTestMode_4() {
	fmt.Println("--- Starting High Concurrency Postgres Test (Mode 4) ---")

	// 1. Настройка БД (Хардкод как в твоем примере, либо тяни из .env)
	dsn := "host=localhost user=postgres password=Postgres123 dbname=banking port=5432 sslmode=disable"
	//if you want to test u have change here as well
	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Очистка и миграция
	db.Migrator().DropTable(&domain.Account{}, &domain.SavingAccount{})
	db.AutoMigrate(&domain.Account{}, &domain.SavingAccount{})

	repo := postgres.NewAccountRepo(db)

	// 2. Создаем тестовые данные
	accA := &domain.Account{ID: "101", Owner: "TestUser1", Balance: 2000} // Дадим больше денег для 20 тестов
	accB := &domain.Account{ID: "102", Owner: "TestUser2", Balance: 1000}

	repo.Create(accA)
	repo.Create(accB)

	initialTotal := accA.Balance + accB.Balance
	fmt.Printf("Initial Total Capital: %.2f\n", initialTotal)

	// 3. Функция трансфера
	transferFunc := func(amount float32) error {
		txRepo, err := repo.BeginTx()
		if err != nil {
			return err
		}

		committed := false
		defer func() {
			if !committed {
				txRepo.RollbackTx()
			}
		}()

		a, err := txRepo.GetByID("101")
		if err != nil {
			return err
		}
		b, err := txRepo.GetByID("102")
		if err != nil {
			return err
		}

		if a.Balance < amount {
			return fmt.Errorf("insufficient funds in acc 101")
		}

		a.Balance -= amount
		b.Balance += amount

		if err := txRepo.Update(a); err != nil {
			return err
		}
		if err := txRepo.Update(b); err != nil {
			return err
		}

		if err := txRepo.CommitTx(); err != nil {
			return err
		}
		committed = true
		return nil
	}

	// 4. Запуск 20 горутин
	var wg sync.WaitGroup
	numGoroutines := 20
	amount := float32(100)

	fmt.Printf("Running %d concurrent transfers of %.2f...\n", numGoroutines, amount)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if err := transferFunc(amount); err != nil {
				fmt.Printf("[Goroutine %d] Error: %v\n", id, err)
			}
		}(i)
	}

	wg.Wait()

	// 5. Финальная проверка
	finalA, _ := repo.GetByID("101")
	finalB, _ := repo.GetByID("102")
	finalTotal := finalA.Balance + finalB.Balance

	fmt.Println("\n=== Simulation Results ===")
	fmt.Printf("Final Account A: %.2f\n", finalA.Balance)
	fmt.Printf("Final Account B: %.2f\n", finalB.Balance)
	fmt.Printf("Final Total Capital: %.2f\n", finalTotal)

	if finalTotal != initialTotal {
		fmt.Printf("FAILED: Money Leaked! Difference: %.2f\n", finalTotal-initialTotal)
	} else {
		fmt.Println("SUCCESS: No money leaked. Transactions are atomic.")
	}
	fmt.Println("-------------------------------------------------------")
}
