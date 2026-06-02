package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/Mfishburn71/pokedex/internal/pokeapi"
)

func commandExplore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("You need to give me a location name")
	}
	if len(args) > 1 {
		return errors.New("Please only provide one location at a time")
	}
	areaName := args[0]
	fmt.Printf("Exploring %s...\n", areaName)

	//Check if location exists
	location, err := cfg.PokeAPIClient.GetLocationArea(areaName)
	if err != nil {
		if strings.Contains(err.Error(), "bad status code") {
			return errors.New("I'm sorry, that's not in my records. Please check your spelling")
		}
		return err
	}

	fmt.Println("Found Pokemon:")
	///////Collects information for area summary
	habitats := make(map[string]struct{}) //Collect Habitat summary
	region := make(map[string]struct{})   //Collect region data
	levelRanges := []int{}
	//for _, rate := range location.EncounterMethodRates {
	for _, rate := range location.PokemonEncounters {
		for _, vd := range rate.VersionDetails {
			region[pokeapi.RegionFor(vd.Version.Name)] = struct{}{}
		}
	}
	for _, encounter := range location.PokemonEncounters {
		for _, vd := range encounter.VersionDetails {
			for _, ed := range vd.EncounterDetails {
				habitats[pokeapi.HabitatFor(ed.Method.Name)] = struct{}{}
				levelRanges = append(levelRanges, ed.MaxLevel)
				levelRanges = append(levelRanges, ed.MinLevel)
			}
		}
	}
	//Listing / Sorting Habitat names
	habitatList := make([]string, 0, len(habitats))
	for h := range habitats {
		habitatList = append(habitatList, h)
	}
	sort.Strings(habitatList)
	//Listing / Sorting Region names
	regionList := make([]string, 0, len(region))
	for r := range region {
		regionList = append(regionList, r)
	}
	sort.Strings(regionList)
	regionDisplay := "Unknown"
	if len(regionList) > 0 {
		regionDisplay = strings.Join(regionList, ", ")
	}
	///print summary here
	tier, avgLevel := pokeapi.LevelFor(levelRanges)

	fmt.Printf("%s is a %s (Level %d) area in the %s region with [%s] habitats.\n",
		prettify(areaName),
		tier,
		avgLevel,      // note: %d for int, not %s
		regionDisplay, // see issue 3
		strings.Join(habitatList, ", "),
	)

	for _, encounter := range location.PokemonEncounters {
		//Collecting Set of encounter types (No dupes)
		methods := make(map[string]struct{})
		for _, vd := range encounter.VersionDetails {
			for _, ed := range vd.EncounterDetails {
				methods[ed.Method.Name] = struct{}{}
			}
		}

		// convert to sorted slice for stable display
		methodList := make([]string, 0, len(methods))
		for m := range methods {
			m = prettify(m)
			methodList = append(methodList, m)
		}
		sort.Strings(methodList)
		fmt.Printf(" - %s [%s]\n", prettify(encounter.Pokemon.Name), strings.Join(methodList, ", "))
	}
	return nil
}

// ///MENTAL NOTE: Add prettify to pokemon names after tests pass
// ///fmt.Printf(" - %s [%s]\n", prettify(encounter.Pokemon.Name), strings.Join(methodList, ", "))
// ///ALSO PRETTIFY LOCATION NAMES
// ///
func prettify(s string) string {
	parts := strings.Split(s, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}
