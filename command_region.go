package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/Mfishburn71/pokedex/internal/pokeapi"
)

var validRegions = buildValidRegions()

func buildValidRegions() map[string]struct{} {
	m := map[string]struct{}{}
	for _, region := range pokeapi.MethodToRegion {
		m[strings.ToLower(region)] = struct{}{}
	}
	return m
}

func commandRegionSearch(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("you need to give me a region name")
	}
	if len(args) > 2 {
		return errors.New("please only provide one region and one search term at a time")
	}

	regionQuery := strings.ToLower(args[0])
	query := ""

	if !matchesRegion(regionQuery) {
		return errors.New("unknown region")
	}
	if len(args) == 2 {
		query = strings.ToLower(args[1])
		fmt.Printf("Looking for %s in %s...\n", query, regionQuery)
	} else {
		fmt.Printf("Looking for %s...\n", regionQuery)
	}

	results := []string{}
	url := pokeapi.BaseURL + "/location-area"

	for {
		resp, err := cfg.PokeAPIClient.ListLocations(url)
		if err != nil {
			return err
		}

		pageResults, err := collectRegionMatchesForPage(cfg, resp.Results, regionQuery, query)
		if err != nil {
			return err
		}
		results = append(results, pageResults...)

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

func matchesRegion(target string) bool {
	_, ok := validRegions[strings.ToLower(target)]
	return ok
}

func locationInRegion(details pokeapi.LocationArea, regionQuery string) bool {
	for _, encounter := range details.PokemonEncounters {
		for _, vd := range encounter.VersionDetails {
			versionName := strings.ToLower(vd.Version.Name)
			regionName := strings.ToLower(pokeapi.MethodToRegion[versionName])
			if regionName == regionQuery {
				return true
			}
		}
	}
	return false
}

func collectRegionMatchesForPage(
	cfg *config,
	areas []pokeapi.LocationAreaSummary,
	regionQuery, query string,
) ([]string, error) {
	var wg sync.WaitGroup
	results := []string{}

	sem := make(chan struct{}, 10)
	resultsCh := make(chan string, 100)
	errCh := make(chan error, 1)

	for _, area := range areas {
		area := area
		wg.Add(1)

		go func() {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			details, err := cfg.PokeAPIClient.GetLocationArea(area.Name)
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}

			if locationInRegion(details, regionQuery) &&
				strings.Contains(strings.ToLower(area.Name), query) {
				resultsCh <- area.Name
			}
		}()
	}

	wg.Wait()
	close(resultsCh)

	for name := range resultsCh {
		results = append(results, name)
	}

	select {
	case err := <-errCh:
		return results, err
	default:
	}
	return results, nil
}
