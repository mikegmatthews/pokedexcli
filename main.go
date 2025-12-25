package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	userInput := bufio.NewScanner(os.Stdin)

	for true {
		fmt.Print("Pokedex > ")

		if userInput.Scan() {
			text := userInput.Text()
			inParams := cleanInput(text)

			if len(inParams) > 0 {
				fmt.Printf("Your command was: %s\n", inParams[0])
			}
		}
	}
}
