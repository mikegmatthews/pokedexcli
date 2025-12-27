package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mikegmatthews/pokedexcli/internal/pokedex"
)

var cliRegistry map[string]cliCommand
var pokeDex *pokedex.Pokedex

func init() {
	mapCmdConfig := cmdConfig{}
	pokeDex = pokedex.NewPokedex()

	cliRegistry = map[string]cliCommand{
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
			description: "Displays the next 20 map locations",
			config:      &mapCmdConfig,
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 map locations",
			config:      &mapCmdConfig,
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Lists all Pokemon encounters for a given area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a given Pokemon",
			callback:    commandCatch,
		},
	}
}

func main() {
	userInput := bufio.NewScanner(os.Stdin)

	for true {
		fmt.Print("Pokedex > ")

		if userInput.Scan() {
			text := userInput.Text()
			inParams := cleanInput(text)

			if len(inParams) > 0 {
				command, ok := cliRegistry[inParams[0]]
				if ok {
					err := command.callback(command.config, inParams[1:]...)
					if err != nil {
						fmt.Println(err)
					}
				} else {
					fmt.Println("Unknown command")
				}
			}
		}
	}
}
