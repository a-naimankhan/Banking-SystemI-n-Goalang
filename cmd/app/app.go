package app

import (
	"BankingSystem/internal/domain"
	"BankingSystem/internal/repository"
	"BankingSystem/internal/service"
	"fmt"
)

func RunApp() {
	//Starting Running Bank Application

	// 1. Инициализируем хранилище
	repo := repository.NewInMemRepo()

	// 2. Создаем аккаунты
	acc1 := &domain.Account{ID: "1", Owner: "Aibar", Balance: 1000}
	acc2 := &domain.Account{ID: "2", Owner: "Brother", Balance: 500}

	repo.Create(acc1)
	repo.Create(acc2)

	// 3. Dependency Injection через интерфейс
	bank := service.NewBankingService(repo)

	// 4. Перевод
	fmt.Println("Processing transfer...")
	if err := bank.Transfer("1", "2", 300); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Success!")
	}

	// 5. Итог
	res1, _ := repo.GetByID("1")
	res2, _ := repo.GetByID("2")
	fmt.Printf("%s: %.2f | %s: %.2f\n", res1.Owner, res1.Balance, res2.Owner, res2.Balance)
}
