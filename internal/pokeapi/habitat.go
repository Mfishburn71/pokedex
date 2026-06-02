package pokeapi

import "sort"

var MethodToHabitat = map[string]string{
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

var MethodToRegion = map[string]string{
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

func HabitatFor(method string) string {
	if h, ok := MethodToHabitat[method]; ok {
		return h
	}
	return "Other" // safe fallback for unmapped methods
}

func RegionFor(method string) string {
	if h, ok := MethodToRegion[method]; ok {
		return h
	}
	return "Other" // safe fallback for unmapped methods
}

func LevelFor(levelRange []int) (string, int) {
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
