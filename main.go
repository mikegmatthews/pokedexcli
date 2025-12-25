package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
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
	}

	userInput := bufio.NewScanner(os.Stdin)

	for true {
		fmt.Print("Pokedex > ")

		if userInput.Scan() {
			text := userInput.Text()
			inParams := cleanInput(text)

			if len(inParams) > 0 {
				command, ok := cliRegistry[inParams[0]]
				if ok {
					err := command.callback()
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
