package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rbonk907/pokedex/internal/pokeapi"
)

type commandConfig struct {
	lastCommand   string
	pokeapiClient pokeapi.Client
	caughtPokes   map[string]pokeapi.Pokemon
	nextUrl       *string
	prevUrl       *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*commandConfig, ...string) error
}

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	commandCfg := &commandConfig{
		caughtPokes:   map[string]pokeapi.Pokemon{},
		pokeapiClient: pokeClient,
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
		args := []string{}
		if len(fields) > 1 {
			args = fields[1:]
		}

		command, ok := getCommands()[fields[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(commandCfg, args...)
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

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
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
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Attempts to catch a pokemon and adds it to the pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon_name>",
			description: "Provides details about caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays a list of caught pokemon",
			callback:    commandPokedex,
		},
	}
}
