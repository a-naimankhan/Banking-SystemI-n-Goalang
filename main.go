package main

import (
	"BankingSystem/cmd/app"
	"fmt"
)

func main() {
	var cmd string
	fmt.Println("If you want to start press 1. \nIf you want to join test mode press 2. ")
	fmt.Scanf("%s", &cmd)

	switch cmd {
	case "1":
		app.RunApp()
	case "2":
		app.RunTestMode()
	default:
		fmt.Println("Unknow Command")
	}

}
