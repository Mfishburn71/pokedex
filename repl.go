package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/Mfishburn71/pokedex/internal/combat"
	"github.com/Mfishburn71/pokedex/internal/pokeapi"
)

type config struct {
	Next          *string
	Previous      *string
	PokeAPIClient pokeapi.Client
	Pokedex       map[string]combat.CaughtPokemon
	CurrentBall   string
	TrainerName   string
	Party         []string
}

func startRepl(commands map[string]cliCommand, client pokeapi.Client) {

	scanner := bufio.NewScanner(os.Stdin)

	cfg := config{
		PokeAPIClient: client,
		Pokedex:       make(map[string]combat.CaughtPokemon),
		CurrentBall:   "pokeball",
	}

	////Auto Load Save file if it exists
	data, err := loadFromFile("save.json")
	if err == nil {
		cfg.TrainerName = data.TrainerName
		cfg.Pokedex = data.Pokedex
		if cfg.Pokedex == nil { //Watch for nil / corrupted pokedex
			cfg.Pokedex = make(map[string]combat.CaughtPokemon)
		}
		cfg.Party = data.Party
		if cfg.Party == nil { //Watch for nil / corrupted party
			cfg.Party = []string{}
		}
		roll := rand.Intn(5)
		switch roll {
		case 1:
			fmt.Printf("Welcome back, %s!\n", cfg.TrainerName)
		case 2:
			fmt.Printf("Good to see you again, %s!\n", cfg.TrainerName)
		case 3:
			fmt.Printf("How have you been, %s?\n", cfg.TrainerName)
		case 4:
			fmt.Printf("Discovered anything new lately, %s?\n", cfg.TrainerName)
		default:
			fmt.Printf("Been busy lately, %s?\n", cfg.TrainerName)
		}

	} else {
		cfg.Pokedex = make(map[string]combat.CaughtPokemon)
		cfg.Party = []string{}

		fmt.Print("Welcome, trainer! What is your name? ")
		scanner.Scan()
		cfg.TrainerName = strings.TrimSpace(scanner.Text())

		fmt.Printf("Nice to meet you, %s!\n", cfg.TrainerName)
	}
	/////Continue as normal
	for {

		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			reply := cleanInput(scanner.Text())
			if len(reply) == 0 {
				continue
			}
			commandName := reply[0]
			args := reply[1:]
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

func displayName(caught combat.CaughtPokemon) string {
	if caught.Nickname != "" {
		return fmt.Sprintf("%s (%s)", caught.Nickname, prettify(caught.Pokemon.Name))
	}
	return prettify(caught.Pokemon.Name)
}
