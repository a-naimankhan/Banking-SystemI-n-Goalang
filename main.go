package main

import "fmt"

func main() {
	var cmd string
	fmt.Println("If you want to start press 1. \nIf you want to join test mode press 2. ")
	fmt.Scanf("%s", &cmd)

	if cmd == "1" {
		RunApp()
	} else if cmd == "2" {
		//To Do test mode !
	}

}
