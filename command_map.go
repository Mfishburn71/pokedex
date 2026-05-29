package main

import (
	"fmt"

	"github.com/Mfishburn71/pokedex/internal/pokeapi"
)

//const defaultURL = "https://pokeapi.co/api/v2/location-area"
//Moved to internal/pokeapi/location_get.go to keep a single variable

type config struct {
	Next          *string
	Previous      *string
	PokeAPIClient pokeapi.Client
	PokemonInfo   pokeapi.PokemonInfo
	Pokedex       map[string]pokeapi.PokemonInfo
}

func commandMap(cfg *config, args ...string) error {
	url := pokeapi.BaseURL + "/location-area"
	if cfg.Next != nil {
		url = *cfg.Next
	}
	return fetchAndDisplayLocations(url, cfg)
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	return fetchAndDisplayLocations(*cfg.Previous, cfg)
}

func fetchAndDisplayLocations(url string, cfg *config) error {
	//resp, err := pokeapi.ListLocations(url)``
	resp, err := cfg.PokeAPIClient.ListLocations(url)
	if err != nil {
		return err
	}
	cfg.Next = resp.Next
	cfg.Previous = resp.Previous
	for _, area := range resp.Results {
		fmt.Println(area.Name)
	}
	return nil
}
