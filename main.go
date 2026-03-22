package main

import (
	"BankingSystem/cmd/app"
	"fmt"
)

func main() {
	var cmd string
	fmt.Println("If you want to start press 1. \nIf you want to join test mode press 2. ")
	fmt.Scanf("%s", &cmd)
	for {
		switch cmd {
		case "1":
			app.RunApp()
		case "2":
			app.RunTestMode_1()
		case "3":
			app.RunTestMode_2()
		case "4":
			app.RunTestMode_3()
		default:
			fmt.Println("Unknow Command")
		}

	}

}
