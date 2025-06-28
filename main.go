package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/rbonk907/pokedex/internal/pokeapi"
)

type commandConfig struct {
	lastCommand string
	nextUrl     *string
	prevUrl     *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*commandConfig) error
}

var commandRegistry map[string]cliCommand

func commandExit(cfg *commandConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil // this will never get call
}

func commandHelp(cfg *commandConfig) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")

	for _, value := range commandRegistry {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}

	return nil
}

func getLocationList(url string) (pokeapi.LocationAreaList, error) {
	res, err := http.Get(url)
	if err != nil {
		return pokeapi.LocationAreaList{}, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	var locationList pokeapi.LocationAreaList
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return pokeapi.LocationAreaList{}, fmt.Errorf("failed to read data into []byte from response: %w", err)
	}

	if err := json.Unmarshal(body, &locationList); err != nil {
		return pokeapi.LocationAreaList{}, fmt.Errorf("error unmarshaling json: %w", err)
	}

	return locationList, nil
}

func commandMap(cfg *commandConfig) error {
	var url string
	if cfg.lastCommand == "map" && cfg.nextUrl != nil {
		url = *cfg.nextUrl
	} else {
		url = "https://pokeapi.co/api/v2/location-area"
	}

	locationList, err := getLocationList(url)
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

func commandMapB(cfg *commandConfig) error {
	if cfg.lastCommand != "map" && cfg.lastCommand != "mapb" {
		return fmt.Errorf("can only call 'mapb' as a consecutive call to 'map' or 'mapb")
	}

	if cfg.prevUrl == nil {
		fmt.Println("You're on the first page")
		return nil
	}

	locationList, err := getLocationList(*cfg.prevUrl)
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
	}

	commandCfg := commandConfig{}

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

		err := command.callback(&commandCfg)
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
