package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
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
	for _, rate := range location.EncounterMethodRates {
		for _, vd := range rate.VersionDetails {
			region[regionFor(vd.Version.Name)] = struct{}{}
		}
	}
	for _, encounter := range location.PokemonEncounters {
		for _, vd := range encounter.VersionDetails {
			for _, ed := range vd.EncounterDetails {
				habitats[habitatFor(ed.Method.Name)] = struct{}{}
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
	tier, avgLevel := levelFor(levelRanges)

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
		fmt.Printf(" - %s [%s]\n", encounter.Pokemon.Name, strings.Join(methodList, ", "))
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

// /////SPLIT THIS PART INTO HABITAT.GO WHEN NEEDED
var methodToHabitat = map[string]string{
	"walk":                "Land",
	"roaming-grass":       "Land",
	"surf":                "Water",
	"old-rod":             "Water",
	"good-rod":            "Water",
	"super-rod":           "Water",
	"rock-smash":          "Rocky",
	"rough-terrain":       "Rocky",
	"headbutt":            "Tree",
	"headbutt-low":        "Tree",
	"headbutt-normal":     "Tree",
	"headbutt-high":       "Tree",
	"honey-tree":          "Tree",
	"dark-grass":          "Tall Grass",
	"grass-spots":         "Tall Grass",
	"cave":                "Cave",
	"cave-spots":          "Cave",
	"seaweed":             "Sea",
	"roaming-water":       "Sea",
	"feebas-tile-fishing": "Sea",
	"overworld-water":     "Sea",
	"super-rod-spots":     "Deep Sea",
	"surf-spots":          "Deep Sea",
	"only-one":            "Unique",
	"gift":                "Gift",
	"gift-egg":            "Gift",
	"npc-trade":           "Gift",
	// ...etc
}

var methodToRegion = map[string]string{
	"red":               "Kanto",
	"blue":              "Kanto",
	"yellow":            "Kanto",
	"gold":              "Johto",
	"silver":            "Johto",
	"crystal":           "Johto",
	"ruby":              "Hoenn",
	"sapphire":          "Hoenn",
	"emerald":           "Hoenn",
	"firered":           "Kanto",
	"leafgreen":         "Kanto",
	"diamond":           "Sinnoh",
	"pearl":             "Sinnoh",
	"platinum":          "Sinnoh",
	"heartgold":         "Johto",
	"soulsilver":        "Johto",
	"black":             "Unova",
	"white":             "Unova",
	"colosseum":         "Orre",
	"xd":                "Orre",
	"black-2":           "Unova",
	"white-2":           "Unova",
	"x":                 "Kalos",
	"y":                 "Kalos",
	"omega-ruby":        "Hoenn",
	"alpha-sapphire":    "Hoenn",
	"sun":               "Alola",
	"moon":              "Alola",
	"ultra-sun":         "Alola",
	"ultra-moon":        "Alola",
	"lets-go-pikachu":   "Kanto",
	"lets-go-eevee":     "Kanto",
	"sword":             "Galar",
	"shield":            "Galar",
	"the-isle-of-armor": "Galar",
	"the-crown-tundra":  "Galar",
	"brilliant-diamond": "Sinnoh",
	"shining-pearl":     "Sinnoh",
	"legends-arceus":    "Hisui",
	"scarlet":           "Paldea",
	"violet":            "Paldea",
	"the-teal-mask":     "Paldea",
	"the-indigo-disk":   "Paldea",
	"red-japan":         "Kanto",
	"green-japan":       "Kanto",
	"blue-japan":        "Kanto",
	"legends-za":        "Kalos",
	"mega-dimension":    "Kalos",
	// ...etc
}

func habitatFor(method string) string {
	if h, ok := methodToHabitat[method]; ok {
		return h
	}
	return "Other" // safe fallback for unmapped methods
}

func regionFor(method string) string {
	if h, ok := methodToRegion[method]; ok {
		return h
	}
	return "Other" // safe fallback for unmapped methods
}

func levelFor(levelRange []int) (string, int) {
	if len(levelRange) == 0 {
		return "Unknown", 0
	}
	sorted := make([]int, len(levelRange))
	copy(sorted, levelRange)
	sort.Ints(sorted)
	avg := levelRange[len(levelRange)/2] //returns the median
	if avg < 21 {
		return "Beginner", avg
	}
	if avg < 41 {
		return "Intermediate", avg
	}
	if avg < 61 {
		return "Advanced", avg
	}
	return "Deadly", avg // captures anything above 60
}
