package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/Mfishburn71/pokedex/internal/combat"
)

func commandCatch(cfg *config, args ...string) error {

	if len(args) == 0 {
		return errors.New("You need to give me a pokemon name")
	}
	if len(args) > 1 {
		return errors.New("Please only provide one pokemon at a time")
	}
	pokemonName := args[0]

	resp, err := cfg.PokeAPIClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	if _, ok := cfg.Pokedex[resp.Name]; ok {
		fmt.Printf("%s is already in your Pokedex\n", prettify(resp.Name))
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	////Bonus for pokeball type
	bonus := 0
	switch cfg.CurrentBall {
	case "greatball":
		bonus = 15
	case "ultraball":
		bonus = 30
	case "masterball":
		bonus = 1000
	default:
		bonus = 0
	}

	chance := 100 - resp.BaseExperience + bonus
	if chance < 15 {
		chance = 15
	}
	roll := rand.Intn(100)

	if roll < chance {
		cfg.Pokedex[resp.Name] = combat.CaughtPokemon{
			Pokemon:  resp,
			BallType: cfg.CurrentBall,
		}
		fmt.Printf("%s was caught!\n", pokemonName)
		//fmt.Println(cfg.Pokedex)
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

func commandPokedex(cfg *config, args ...string) error {
	names := make([]string, 0, len(cfg.Pokedex))
	for name := range cfg.Pokedex {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Println("Your Pokedex:")
	for _, name := range names {
		fmt.Printf(" - %s\n", prettify(name))
	}
	return nil
}

func commandBall(cfg *config, args ...string) error {
	if len(args) == 0 {
		cfg.CurrentBall = "pokeball"
		fmt.Println("No ball specified. Defaulting to pokeball")
		return nil
	}
	if len(args) > 1 {
		return errors.New("You can only throw one ball at a time!")
	}

	ball := strings.ToLower(strings.ReplaceAll(args[0], " ", ""))

	validBalls := map[string]bool{
		"pokeball": true, "greatball": true,
		"ultraball": true, "masterball": true,
	}
	if !validBalls[ball] {
		return errors.New("I don't have any data for that kind of pokeball")
	}
	if ball == "masterball" {
		fmt.Println("You found a Master Ball!")
	}
	cfg.CurrentBall = ball

	fmt.Printf("Equipped %s\n", prettify(cfg.CurrentBall))
	return nil
}
