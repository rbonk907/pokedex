package main

import "fmt"

func commandPokedex(cfg *commandConfig, params ...string) error {
	if len(cfg.caughtPokes) == 0 {
		fmt.Println("You have not caught any pokemon")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for _, pokemon := range cfg.caughtPokes {
		fmt.Printf(" - %s\n", pokemon.Name)
	}

	return nil
}
