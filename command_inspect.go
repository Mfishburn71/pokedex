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

	fmt.Printf("Name: %s\n", prettify(displayName(pokemon)))
	fmt.Printf("Height: %d\n", pokemon.Pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Pokemon.Weight)
	fmt.Printf("Base Experience: %d\n", pokemon.Pokemon.BaseExperience)

	fmt.Println("Types:")
	for _, t := range pokemon.Pokemon.Types {
		fmt.Printf("  - %s\n", prettify(t.Type.Name))
	}

	fmt.Println("Stats:")
	for _, stat := range pokemon.Pokemon.Stats {
		fmt.Printf("  - %s: %d\n", prettify(stat.Stat.Name), stat.BaseStat)
	}

	return nil
}
