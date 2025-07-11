package main

import (
	"fmt"
	"os"
)

func commandExit(cfg *commandConfig, params ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil // this will never get call
}
