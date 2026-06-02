package main

import (
	"fmt"
	"os"
)

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	data := SaveData{
		TrainerName: cfg.TrainerName,
		Pokedex:     cfg.Pokedex,
		Party:       cfg.Party,
	}
	if err := saveToFile("save.json", data); err != nil {
		return err
	}
	fmt.Println("Pokedex saved to save.json")
	cfg.PokeAPIClient.Stop()
	os.Exit(0)
	return nil // unreachable, satisfies error return signature
}
