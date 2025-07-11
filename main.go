package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rbonk907/pokedex/internal/pokeapi"
	"github.com/rbonk907/pokedex/internal/pokecache"
)

type commandConfig struct {
	lastCommand string
	cache       *pokecache.Cache
	nextUrl     *string
	prevUrl     *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*commandConfig, ...string) error
}

var commandRegistry map[string]cliCommand

func getLocationList(url string, cfg *commandConfig) (pokeapi.LocationAreaList, error) {
	bytes, found := cfg.cache.Get(url)
	if !found {
		res, err := http.Get(url)
		if err != nil {
			return pokeapi.LocationAreaList{}, fmt.Errorf("error making request: %w", err)
		}
		defer res.Body.Close()

		bytes, err = io.ReadAll(res.Body)
		if err != nil {
			return pokeapi.LocationAreaList{}, fmt.Errorf("failed to read data into []byte from response: %w", err)
		}

		// add the byte array to the map
		cfg.cache.Add(url, bytes)
	}

	var locationList pokeapi.LocationAreaList
	if err := json.Unmarshal(bytes, &locationList); err != nil {
		return pokeapi.LocationAreaList{}, fmt.Errorf("error unmarshaling json: %w", err)
	}

	return locationList, nil
}

func commandMap(cfg *commandConfig, params ...string) error {
	var url string
	if cfg.lastCommand == "map" && cfg.nextUrl != nil {
		url = *cfg.nextUrl
	} else {
		url = "https://pokeapi.co/api/v2/location-area"
	}

	locationList, err := getLocationList(url, cfg)
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

	locationList, err := getLocationList(*cfg.prevUrl, cfg)
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

func commandExplore(cfg *commandConfig, params ...string) error {
	if len(params) < 2 {
		return fmt.Errorf("no location provided. Example usage: explore <area_name>")
	}

	url := "https://pokeapi.co/api/v2/location-area/" + params[1]

	bytes, found := cfg.cache.Get(url)
	if !found {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer res.Body.Close()

		bytes, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read data into []byte from response: %w", err)
		}

		// add the byte array to the map
		cfg.cache.Add(url, bytes)
	}

	var location pokeapi.LocationAreaDetail
	if err := json.Unmarshal(bytes, &location); err != nil {
		return fmt.Errorf("error unmarshaling json: %w", err)
	}

	fmt.Printf("Exploring %s...\n", params[1])
	if len(location.PokemonEncounters) < 1 {
		return fmt.Errorf("no pokemon found")
	}

	fmt.Println("Found Pokemon:")
	for _, pokeEncounter := range location.PokemonEncounters {
		fmt.Printf("- %s\n", pokeEncounter.Pokemon.Name)
	}

	return nil
}

func main() {
	// init
	commandRegistry = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays a paginated list of locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous page of locations",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Lists possible pokemon encounters at a given location",
			callback:    commandExplore,
		},
	}

	cache := pokecache.NewCache(5 * time.Second)
	commandCfg := commandConfig{
		cache: cache,
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		text := scanner.Text()
		if text == "" {
			continue
		}
		fields := cleanInput(text)

		command, ok := commandRegistry[fields[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(&commandCfg, fields...)
		if err != nil {
			fmt.Println(err)
		}

		commandCfg.lastCommand = command.name
	}
}

func cleanInput(text string) []string {
	words := strings.Fields(text)

	for i, word := range words {
		words[i] = strings.ToLower(word)
	}

	return words
}
