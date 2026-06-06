//go:build !(js && wasm)

package main

import (
	"context"

	"github.com/Mfishburn71/pokedex/internal/pokeapi"
)

func main() {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays a list of map locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous list of map locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Returns a list of pokemon inside a given area and basic environent details",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Throws a pokeball at a Pokemon",
			callback:    commandCatch,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List the Pokemon stored in your Pokedex",
			callback:    commandPokedex,
		},
		"inspect": {
			name:        "inspect",
			description: "List detailed information on the Pokemon stored in your Pokedex",
			callback:    commandInspect,
		},
		"ball": {
			name:        "ball",
			description: "choose a type of pokeball to throw. Great Ball and Ultra Ball make things easier to catch",
			callback:    commandBall,
		},
		"battle": {
			name:        "battle",
			description: "battle: Simulate a basic battle between two pokemon. You must give it two different names",
			callback:    commandBattle,
		},
		"save": {
			name:        "save",
			description: "Save trainer data to local storage",
			callback:    commandSave,
		},
		"load": {
			name:        "load",
			description: "Load trainer data from local storage",
			callback:    commandLoad,
		},
		"delete": {
			name:        "delete",
			description: "Delete trainer data entirely",
			callback:    commandDelete,
		},
		"party": {
			name:        "party",
			description: "List your current party",
			callback:    commandParty,
		},
		"partyadd": {
			name:        "partyadd",
			description: "Add a pokemon to your party (Maximum 6 at a time)",
			callback:    commandPartyAdd,
		},
		"partyremove": {
			name:        "partyremove",
			description: "Remove a pokemom from your party",
			callback:    commandPartyRemove,
		},
		"nickname": {
			name:        "nickname",
			description: "Give a nickname to a pokemon. Type Pokemon name first, then your chosen nickname.",
			callback:    commandNickname,
		},
		"search": {
			name:        "search",
			description: "Search for a specific map location",
			callback:    commandSearch,
		},
		"debug": {
			name:        "debug",
			description: "Read RAW information only",
			callback:    commandDebugArea,
		},
		"region": {
			name:        "region",
			description: "Search inside a specific region, with up to one search term",
			callback:    commandRegionSearch,
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := pokeapi.NewClient(ctx)
	startRepl(commands, client)

}
