package main

import (
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var cliRegistry map[string]cliCommand

func cleanInput(text string) []string {
	cleaned := make([]string, 0)

	for f := range strings.FieldsSeq(text) {
		cleaned = append(cleaned, strings.ToLower(f))
	}

	return cleaned
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cli := range cliRegistry {
		fmt.Printf("%s: %s\n", cli.name, cli.description)
	}

	return nil
}
