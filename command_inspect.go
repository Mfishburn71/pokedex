package main

import (
	"errors"
	"fmt"
)

func commandInspect(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("you need to give me a pokemon name")
	}
	if len(args) > 1 {
		return errors.New("please only provide one pokemon at a time")
	}

	name := args[0]
	pokemon, ok := cfg.Pokedex[name]
	if !ok {
		_, err := cfg.PokeAPIClient.GetPokemon(name)
		if err != nil {
			return err
		}

		return errors.New("you have not caught that pokemon")
	}

	fmt.Printf("Name: %s\n", prettify(pokemon.Name))
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Base Experience: %d\n", pokemon.BaseExperience)

	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", prettify(t.Type.Name))
	}

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: %d\n", prettify(stat.Stat.Name), stat.BaseStat)
	}

	return nil
}
