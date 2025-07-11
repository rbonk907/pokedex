package main

import (
	"fmt"
	"math/rand"
)

func commandCatch(cfg *commandConfig, params ...string) error {
	if len(params) != 1 {
		return fmt.Errorf("no pokemon provided. Example usage: catch <pokemon_name>")
	}

	pokeName := params[0]

	pokemon, err := cfg.pokeapiClient.GetPokemon(pokeName)
	if err != nil {
		return err
	}

	baseExperience := pokemon.BaseExperience
	maxBaseExperience := 1_640_000
	catchChance := 1.0 - float64(baseExperience)/float64(maxBaseExperience)

	fmt.Printf("Throwing a Pokeball at %s...\n", pokeName)
	if rand.Float64() < catchChance {
		fmt.Printf("%s was caught!\n", pokeName)
	} else {
		fmt.Printf("%s escaped!\n", pokeName)
	}
	// add the pokemon to the pokedex
	cfg.caughtPokes[pokeName] = pokemon

	return nil
}
