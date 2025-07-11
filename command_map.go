package main

import (
	"fmt"
)

func commandMap(cfg *commandConfig, params ...string) error {
	var url string
	if cfg.lastCommand == "map" && cfg.nextUrl != nil {
		url = *cfg.nextUrl
	} else {
		url = "https://pokeapi.co/api/v2/location-area"
	}

	locationList, err := cfg.pokeapiClient.GetLocationList(url)
	if err != nil {
		return err
	}

	cfg.nextUrl = locationList.Next
	cfg.prevUrl = locationList.Previous

	for _, loc := range locationList.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandMapB(cfg *commandConfig, params ...string) error {
	if cfg.lastCommand != "map" && cfg.lastCommand != "mapb" {
		return fmt.Errorf("can only call 'mapb' as a consecutive call to 'map' or 'mapb")
	}

	if cfg.prevUrl == nil {
		fmt.Println("You're on the first page")
		return nil
	}

	locationList, err := cfg.pokeapiClient.GetLocationList(*cfg.prevUrl)
	if err != nil {
		return err
	}

	cfg.nextUrl = locationList.Next
	cfg.prevUrl = locationList.Previous

	for _, loc := range locationList.Results {
		fmt.Println(loc.Name)
	}
	return nil
}
