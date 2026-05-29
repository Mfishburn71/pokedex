package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
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

	chance := 100 - resp.BaseExperience
	if chance < 5 {
		chance = 5
	}
	roll := rand.Intn(100)

	if roll < chance {
		cfg.Pokedex[resp.Name] = resp
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
