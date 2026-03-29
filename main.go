package main

import (
	"BankingSystem/cmd/app"
	"fmt"
	"os"
)

func main() {
	for {
		fmt.Println("\n--- Banking System Menu ---")
		fmt.Println("1. Run Standard App")
		fmt.Println("2. Run Test Mode 1 (In-Memory Concurrency)")
		fmt.Println("3. Run Test Mode 2 (Interest Accrual Basics)")
		fmt.Println("4. Run Test Mode 3 (Interest Worker Lifecycle)")
		fmt.Println("5. Run Test Mode 4 (Postgres High Concurrency)")
		fmt.Println("q. Exit")
		fmt.Print("Choose an option: ")

		var cmd string
		fmt.Scanln(&cmd)

		switch cmd {
		case "1":
			// Обычный запуск приложения
			app.RunApp()

		case "2":
			// Стресс-тест в памяти: 50 воркеров делают по 100 рандомных переводов.
			// Проверка Thread-Safety (мьютексов) и целостности итогового капитала.
			app.RunTestMode_1()

		case "3":
			// Простая проверка логики начисления процентов на SavingAccount.
			// Создает счета и имитирует ожидание для проверки формулы профита.
			app.RunTestMode_2()

		case "4":
			// Тест жизненного цикла фонового воркера.
			// Проверяет корректный запуск (Start) и мягкую остановку (Stop) тикера процентов.
			app.RunTestMode_3()

		case "5":
			// Хардкорный тест на ACID в Postgres.
			// Запускает 20 конкурентных транзакций через БД.
			// Проверяет атомарность: чтобы ни один цент не пропал при UPDATE в разных горутинах.
			app.RunTestMode_4()

		case "q":
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Unknown Command, please try again.")
		}
	}
}
