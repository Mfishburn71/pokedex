package main

import "fmt"

func commandHelp(cfg *config, args ...string) error {
	fmt.Println(`Welcome to the Pokedex!
Usage:

help: Display a help message
map: Display the next 20 map locations
mapb: Display the previous 20 map locations
search: Search for a specific map location. Requires one exact name.
region: Search inside a specific region, with up to one search term.
explore: Show Pokémon in a given area and basic environment details. Requires one exact name.
catch: Throw a Pokeball at a Pokémon. Requires a pokemon name.
pokedex: List the Pokemon stored in your Pokedex
inspect: List detailed information on a Pokemon stored in your Pokedex
ball: choose a type of pokeball to throw. Great Ball and Ultra Ball make things easier to catch
battle: Simulate a basic battle between two pokemon. You must give it two different names from your pokedex.
save: Save trainer data to local storage
load: Load trainer data from local storage
party: List your current party
partyadd: Add a pokemon to your party (Maximum 6 at a time)
partyremove: Remove a pokemon from your party
nickname: Give a nickname to a pokemon. Type Pokemon name first, then your chosen nickname.
delete: Delete trainer data entirely
exit: Exit the Pokedex`)
	return nil
}
