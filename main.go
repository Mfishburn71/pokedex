package main

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
	}

	startRepl(commands)

}
