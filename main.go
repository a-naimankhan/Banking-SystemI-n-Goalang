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
			app.RunApp()
		case "2":
			app.RunTestMode_1()
		case "3":
			app.RunTestMode_2()
		case "4":
			app.RunTestMode_3()
		case "5":
			app.RunTestMode_4()
		case "q":
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Unknown Command, please try again.")
		}
	}
}
