package main

import "fmt"

func commandHelp(cfg *config, args ...string) error {
	fmt.Println(`Welcome to the Pokedex!
Usage:

help: Displays a help message
map: Displays a list of 20 map locations
mapb: Displays the previous list of map locations
explore: Returns a list of pokemon inside a given area and basic environent details
catch: Throws a pokeball at a Pokemon
pokedex: List the Pokemon stored in your Pokedex
inspect: List detailed information on the Pokemon stored in your Pokedex
exit: Exit the Pokedex`)
	return nil
}
