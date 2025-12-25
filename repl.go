package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mikegmatthews/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*cmdConfig) error
	config      *cmdConfig
}

type cmdConfig struct {
	next     *string
	previous *string
}

var cliRegistry map[string]cliCommand

func cleanInput(text string) []string {
	cleaned := make([]string, 0)

	for f := range strings.FieldsSeq(text) {
		cleaned = append(cleaned, strings.ToLower(f))
	}

	return cleaned
}

func commandExit(_ *cmdConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *cmdConfig) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cli := range cliRegistry {
		fmt.Printf("%s: %s\n", cli.name, cli.description)
	}

	return nil
}

func commandMap(cfg *cmdConfig) error {
	var nextUrl string
	if cfg.next != nil {
		nextUrl = *cfg.next
	} else {
		nextUrl = "https://pokeapi.co/api/v2/location-area/"
	}

	areas, err := pokeapi.GetLocationAreas(nextUrl)
	if err != nil {
		return err
	}

	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}

	cfg.next = areas.Next
	cfg.previous = areas.Previous

	return nil
}

func commandMapB(cfg *cmdConfig) error {
	if cfg.previous == nil {
		fmt.Println("You're on the first page")
		return nil
	}

	areas, err := pokeapi.GetLocationAreas(*cfg.previous)
	if err != nil {
		return err
	}

	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}

	cfg.next = areas.Next
	cfg.previous = areas.Previous

	return nil
}
