package postgres_test

import (
	"BankingSystem/internal/domain"
	"BankingSystem/internal/repository/postgres"
	"fmt"
	"sync"
	"testing"

	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {

	dsn := "host=localhost user=postgres password=Postgres123 dbname=banking port=5432 sslmode=disable"
	//if you want to test u have to change here as well
	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Очищаем таблицы перед тестом, чтобы старые данные не мешали
	db.Migrator().DropTable(&domain.Account{}, &domain.SavingAccount{})
	db.AutoMigrate(&domain.Account{}, &domain.SavingAccount{})

	return db
}
func TestPostgresTransfer(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewAccountRepo(db)

	// 1. Создаем тестовые данные
	accA := &domain.Account{ID: "101", Owner: "TestUser1", Balance: 1000}
	accB := &domain.Account{ID: "102", Owner: "TestUser2", Balance: 1000}

	repo.Create(accA)
	repo.Create(accB)

	initialTotal := accA.Balance + accB.Balance // 2000

	// 2. Функция для выполнения трансфера в транзакции
	transfer := func(amount float32) error {
		txRepo, err := repo.BeginTx()
		if err != nil {
			return err
		}
		committed := false
		defer func() {
			if !committed {
				_ = txRepo.RollbackTx()
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
			return fmt.Errorf("insufficient funds")
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

	// 3. Запускаем 20 конкурентных трансферов по 100
	var wg sync.WaitGroup
	numGoroutines := 20
	amount := float32(100)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := transfer(amount); err != nil {
				t.Errorf("Transfer failed: %v", err)
			}
		}()
	}

	wg.Wait()

	// 4. Проверка результата: общий капитал должен остаться 2000
	finalA, _ := repo.GetByID("101")
	finalB, _ := repo.GetByID("102")
	finalTotal := finalA.Balance + finalB.Balance

	if finalTotal != initialTotal {
		t.Errorf("Money leak detected! Initial total: %.2f, Final total: %.2f", initialTotal, finalTotal)
	}

	// Дополнительная проверка: A должно уменьшиться на 2000, B увеличиться на 2000
	expectedA := 1000 - float32(numGoroutines)*amount
	expectedB := 1000 + float32(numGoroutines)*amount

	if finalA.Balance != expectedA {
		t.Errorf("Account A balance: expected %.2f, got %.2f", expectedA, finalA.Balance)
	}
	if finalB.Balance != expectedB {
		t.Errorf("Account B balance: expected %.2f, got %.2f", expectedB, finalB.Balance)
	}
}
