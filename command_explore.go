package main

import (
	"fmt"
)

func commandExplore(cfg *commandConfig, params ...string) error {
	if len(params) != 1 {
		return fmt.Errorf("no location provided. Example usage: explore <area_name>")
	}

	location, err := cfg.pokeapiClient.GetLocation(params[0])
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", params[0])
	if len(location.PokemonEncounters) < 1 {
		return fmt.Errorf("no pokemon found")
	}

	fmt.Println("Found Pokemon:")
	for _, pokeEncounter := range location.PokemonEncounters {
		fmt.Printf("- %s\n", pokeEncounter.Pokemon.Name)
	}

	return nil
}
