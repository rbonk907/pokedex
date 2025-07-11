package main

import (
	"fmt"
)

func commandInspect(cfg *commandConfig, params ...string) error {
	if len(params) != 1 {
		return fmt.Errorf("no pokemon provided. Example usage: inspect <pokemon_name>")
	}

	pokemon, found := cfg.caughtPokes[params[0]]
	if !found {
		return fmt.Errorf("you have not caught that pokemon")
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" -%v:%v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, pokeType := range pokemon.Types {
		fmt.Printf(" - %s\n", pokeType.Type.Name)
	}

	return nil
}
