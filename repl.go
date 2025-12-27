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
		fmt.Println("You may now inspect it with the inspect command.")
		pokeDex.Caught(pokemon)
	} else {
		fmt.Printf("%s escaped!\n", pokeName)
	}

	return nil
}

func commandInspect(_ *cmdConfig, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("inspect requies a Pokemon name")
	}

	pokeName := args[0]
	pokemon := pokeDex.Inspect(pokeName)

	if pokemon == nil {
		return fmt.Errorf("you have not caught that pokemon")
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("\t- %s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("\t- %s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(_ *cmdConfig, _ ...string) error {
	pokemon := pokeDex.GetAllPokemonNames()
	if len(pokemon) == 0 {
		fmt.Println("No Pokemon have been caught yet")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for _, name := range pokemon {
		fmt.Printf("\t- %s\n", name)
	}

	return nil
}
