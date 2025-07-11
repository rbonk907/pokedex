package main

import (
	"fmt"
)

func commandHelp(cfg *commandConfig, params ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")

	for _, value := range commandRegistry {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}

	return nil
}
