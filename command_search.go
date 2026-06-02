package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/Mfishburn71/pokedex/internal/pokeapi"
)

func commandSearch(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("you need to give me a location name")
	}
	if len(args) > 1 {
		return errors.New("please only provide one search term at a time")
	}
	query := strings.ToLower(args[0])
	fmt.Printf("Looking for %s...\n", query)

	results := []string{}
	url := pokeapi.BaseURL + "/location-area"

	for {
		resp, err := cfg.PokeAPIClient.ListLocations(url)
		if err != nil {
			return err
		}

		for _, area := range resp.Results {
			if strings.Contains(strings.ToLower(area.Name), query) {
				results = append(results, area.Name)
			}
		}

		if resp.Next == nil {
			break
		}
		url = *resp.Next
	}

	if len(results) == 0 {
		fmt.Println("No matches found")
		return nil
	}

	sort.Strings(results)

	fmt.Printf("Results found: %d\n", len(results))
	for _, result := range results {
		fmt.Println(result)
	}

	return nil
}

func commandDebugArea(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("provide one area name")
	}

	data, err := cfg.PokeAPIClient.GetLocationAreaRaw(args[0])
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
