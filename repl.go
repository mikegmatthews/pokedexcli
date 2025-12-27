package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/mikegmatthews/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*cmdConfig, ...string) error
	config      *cmdConfig
}

type cmdConfig struct {
	next     *string
	previous *string
}

func cleanInput(text string) []string {
	cleaned := make([]string, 0)

	for f := range strings.FieldsSeq(text) {
		cleaned = append(cleaned, strings.ToLower(f))
	}

	return cleaned
}

func commandExit(_ *cmdConfig, _ ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *cmdConfig, _ ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cli := range cliRegistry {
		fmt.Printf("%s: %s\n", cli.name, cli.description)
	}

	return nil
}

func commandMap(cfg *cmdConfig, _ ...string) error {
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

func commandMapB(cfg *cmdConfig, _ ...string) error {
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

func commandExplore(_ *cmdConfig, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("explore requires a location name")
	}

	areaName := args[0]
	pokemon, err := pokeapi.GetPokemonEncounters(areaName)
	if err != nil {
		return err
	}

	for _, p := range pokemon {
		fmt.Println(p)
	}

	return nil
}

func commandCatch(_ *cmdConfig, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("catch requires a Pokemon name")
	}

	pokeName := args[0]
	pokemon, err := pokeapi.GetPokemon(pokeName)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokeName)
	baseExp := pokemon.BaseExperience
	if rand.Intn(baseExp) > (baseExp / 2) {
		fmt.Printf("%s was caught!\n", pokeName)
		pokeDex.Caught(pokemon)
	} else {
		fmt.Printf("%s escaped!\n", pokeName)
	}

	return nil
}
