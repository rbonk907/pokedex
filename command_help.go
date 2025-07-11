package main

import (
	"fmt"
)

func commandHelp(cfg *commandConfig, params ...string) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")

	for _, value := range getCommands() {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	fmt.Println()

	return nil
}
