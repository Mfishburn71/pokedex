package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Mfishburn71/pokedex/internal/pokeapi"
)

func startRepl(commands map[string]cliCommand) {

	scanner := bufio.NewScanner(os.Stdin)

	cfg := config{
		PokeAPIClient: pokeapi.NewClient(),
		Pokedex:       map[string]pokeapi.PokemonInfo{},
	}

	for {

		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			reply := cleanInput(scanner.Text())
			if len(reply) == 0 {
				continue
			}
			commandName := reply[0]
			args := []string{}
			if len(reply) > 1 {
				args = reply[1:]
			}
			if len(reply) > 0 {
				//fmt.Printf("Your command was: %v\n", reply[0])
				cmd, ok := commands[commandName]
				if ok {
					if err := cmd.callback(&cfg, args...); err != nil {
						fmt.Println("error:", err)
					}
				} else {
					fmt.Println("Unknown Command")
				}
			}
		}
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, args ...string) error
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
