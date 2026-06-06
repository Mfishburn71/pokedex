//go:build js && wasm

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"syscall/js"

	"github.com/Mfishburn71/pokedex/internal/combat"
	"github.com/Mfishburn71/pokedex/internal/pokeapi"
)

var cfg = config{}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg = config{
		PokeAPIClient: pokeapi.NewClient(ctx),
		Pokedex:       make(map[string]combat.CaughtPokemon),
		CurrentBall:   "pokeball",
		Party:         []string{},
	}

	// Register functions
	js.Global().Set("showHelp", js.FuncOf(jsHelp))
	js.Global().Set("listMap", js.FuncOf(jsMap))
	js.Global().Set("listMapB", js.FuncOf(jsMapB))
	js.Global().Set("searchArea", js.FuncOf(jsSearch))
	js.Global().Set("searchRegion", js.FuncOf(jsRegion))
	js.Global().Set("exploreArea", js.FuncOf(jsExplore))
	js.Global().Set("catchPokemon", js.FuncOf(jsCatch))
	js.Global().Set("openPokedex", js.FuncOf(jsPokedex))
	js.Global().Set("inspectPokemon", js.FuncOf(jsInspect))
	js.Global().Set("equipokeball", js.FuncOf(jsBall))
	js.Global().Set("createBattle", js.FuncOf(jsBattle))
	js.Global().Set("saveData", js.FuncOf(jsSave))
	js.Global().Set("loadData", js.FuncOf(jsLoad))
	js.Global().Set("listParty", js.FuncOf(jsParty))
	js.Global().Set("addParty", js.FuncOf(jsPartyAdd))
	js.Global().Set("removeParty", js.FuncOf(jsPartyRemove))
	js.Global().Set("setNickname", js.FuncOf(jsNickname))
	js.Global().Set("deleteData", js.FuncOf(jsDelete))
	js.Global().Set("setTrainerName", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) > 0 {
			cfg.TrainerName = args[0].String()
		}
		return nil
	}))

	js.Global().Set("getTrainerInfo", js.FuncOf(func(this js.Value, args []js.Value) any {
		type info struct {
			Name string `json:"name"`
			Ball string `json:"ball"`
		}
		data, _ := json.Marshal(info{Name: cfg.TrainerName, Ball: cfg.CurrentBall})
		return string(data)
	}))

	save := js.Global().Get("localStorage").Call("getItem", "pokedex_save")
	if !save.IsNull() && !save.IsUndefined() {
		// existing save found, initialize from it
		initFromSave(save.String())
	} else {
		// signal JS to show the name modal
		js.Global().Call("showNameModal")
	}

	// Keep the program alive
	select {}
}

type catchResult struct {
	Caught  bool                  `json:"caught"`
	Name    string                `json:"name"`
	Message string                `json:"message,omitempty"`
	Pokemon *combat.CaughtPokemon `json:"pokemon,omitempty"`
	Sprite  string                `json:"sprite,omitempty"`
}

func jsCatch(this js.Value, args []js.Value) any {
	name := ""
	if len(args) > 0 {
		name = args[0].String()
	}

	// Check duplicate before fetching
	if _, ok := cfg.Pokedex[name]; ok {
		result := catchResult{Name: name, Message: prettify(name) + " is already in your Pokedex!"}
		data, _ := json.Marshal(result)
		return string(data)
	}

	// Fetch pokemon info upfront so we have the sprite regardless of outcome
	pokemon, err := cfg.PokeAPIClient.GetPokemon(name)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}

	err = commandCatch(&cfg, name)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}

	sprite := pokemon.Sprites.FrontDefault

	caught, ok := cfg.Pokedex[name]
	if !ok {
		result := catchResult{Name: name, Message: name + " escaped!", Sprite: sprite}
		data, _ := json.Marshal(result)
		return string(data)
	}

	result := catchResult{Caught: true, Name: name, Message: name + " was caught!", Pokemon: &caught, Sprite: sprite}
	data, _ := json.Marshal(result)
	return string(data)
}

func jsPokedex(this js.Value, args []js.Value) any {
	data, err := json.Marshal(cfg.Pokedex)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	return string(data)
}

func jsMap(this js.Value, args []js.Value) any {
	url := "https://pokeapi.co/api/v2" + "/location-area"
	if cfg.Next != nil {
		url = *cfg.Next
	}
	resp, err := cfg.PokeAPIClient.ListLocations(url)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	cfg.Next = resp.Next
	cfg.Previous = resp.Previous
	data, _ := json.Marshal(resp.Results)
	return string(data)
}

func jsMapB(this js.Value, args []js.Value) any {
	if cfg.Previous == nil {
		return `{"error": "you're on the first page"}`
	}
	resp, err := cfg.PokeAPIClient.ListLocations(*cfg.Previous)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	cfg.Next = resp.Next
	cfg.Previous = resp.Previous
	data, _ := json.Marshal(resp.Results)
	return string(data)
}

func jsHelp(this js.Value, args []js.Value) any {
	type helpEntry struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	entries := []helpEntry{
		{"Pokedex", "List the Pokemon in your Pokedex. Manage your party and nicknames here."},
		{"Inspect", "Detailed information on a caught Pokemon. You must have caught the pokemon before loading its data."},
		{"Catch", "Throw a Pokeball at a Pokemon. Requires a pokemon name. For best results, equip a high-grade pokeball."},
		{"Pokeballs", "Choose which pokeball type to throw"},
		{"Battle", "Simulate a battle between two of your Pokemon"},

		{"Map", "Display the master map list, 20 locations at a time. Press again to advance the list."},
		{"Map Back", "Go back to the last map list, if you advanced too far."},
		{"Explore", "Show Pokemon habitats and difficulty for a given area. Name must be spelled exactly; get the names from the Map."},
		{"Search", "Search for specific map locations from the master list, like 'Eterna'. One word please. "},
		{"Region", "Lists only locations inside a specific region like `Kanto`. You can also use a search term like `cave` to list only the caves in a given region"},

		{"Help", "This is the screen you're reading now! Lists all the Pokedex features."},
		{"Party", "Shows your current party. You can add Pokemon from the Pokedex screen."},
		{"Nickname", "Give a nickname to a Pokemon. This works for all Pokedex commands, including battles. Found on the Pokedex screen."},
		{"Trainer Data", "Save, load, or delete trainer data. Make sure you save before you close the browser!"},
	}

	data, _ := json.Marshal(entries)
	return string(data)
}

func jsSearch(this js.Value, args []js.Value) any {
	query := ""
	if len(args) > 0 {
		query = strings.ToLower(args[0].String())
	}

	results := []string{}
	url := pokeapi.BaseURL + "/location-area"

	for {
		resp, err := cfg.PokeAPIClient.ListLocations(url)
		if err != nil {
			return `{"error": "` + err.Error() + `"}`
		}
		for _, area := range resp.Results {
			if strings.Contains(strings.ToLower(area.Name), query) {
				results = append(results, prettify(area.Name))
			}
		}
		if resp.Next == nil {
			break
		}
		url = *resp.Next
	}

	type searchResult struct {
		Count   int      `json:"count"`
		Results []string `json:"results"`
	}

	data, _ := json.Marshal(searchResult{
		Count:   len(results),
		Results: results,
	})
	return string(data)
}

func jsRegion(this js.Value, args []js.Value) any {
	regionQuery := ""
	query := ""

	if len(args) > 0 {
		regionQuery = strings.ToLower(args[0].String())
	}
	if len(args) > 1 {
		query = strings.ToLower(args[1].String())
	}

	if !matchesRegion(regionQuery) {
		return `{"error": "unknown region"}`
	}

	results := []string{}
	url := pokeapi.BaseURL + "/location-area"

	for {
		resp, err := cfg.PokeAPIClient.ListLocations(url)
		if err != nil {
			return `{"error": "` + err.Error() + `"}`
		}

		pageResults, err := collectRegionMatchesForPage(&cfg, resp.Results, regionQuery, query)
		if err != nil {
			return `{"error": "` + err.Error() + `"}`
		}
		results = append(results, pageResults...)

		if resp.Next == nil {
			break
		}
		url = *resp.Next
	}

	type regionResult struct {
		Count   int      `json:"count"`
		Region  string   `json:"region"`
		Results []string `json:"results"`
	}

	data, _ := json.Marshal(regionResult{
		Count:   len(results),
		Region:  regionQuery,
		Results: results,
	})
	return string(data)
}

type exploreResult struct {
	Name     string             `json:"name"`
	Display  string             `json:"display"`
	Tier     string             `json:"tier"`
	Level    int                `json:"level"`
	Regions  []string           `json:"regions"`
	Habitats []string           `json:"habitats"`
	Pokemon  []exploreEncounter `json:"pokemon"`
}

type exploreEncounter struct {
	Name    string   `json:"name"`
	Display string   `json:"display"`
	Methods []string `json:"methods"`
}

func jsExplore(this js.Value, args []js.Value) any {
	areaName := ""
	if len(args) > 0 {
		areaName = args[0].String()
	}

	location, err := cfg.PokeAPIClient.GetLocationArea(areaName)
	if err != nil {
		return `{"error": "I'm sorry, that's not in my records. Please check your spelling"}`
	}

	habitats := make(map[string]struct{})
	regions := make(map[string]struct{})
	levelRanges := []int{}

	for _, encounter := range location.PokemonEncounters {
		for _, vd := range encounter.VersionDetails {
			regions[pokeapi.RegionFor(vd.Version.Name)] = struct{}{}
			for _, ed := range vd.EncounterDetails {
				habitats[pokeapi.HabitatFor(ed.Method.Name)] = struct{}{}
				levelRanges = append(levelRanges, ed.MaxLevel)
				levelRanges = append(levelRanges, ed.MinLevel)
			}
		}
	}

	habitatList := make([]string, 0, len(habitats))
	for h := range habitats {
		habitatList = append(habitatList, h)
	}
	sort.Strings(habitatList)

	regionList := make([]string, 0, len(regions))
	for r := range regions {
		regionList = append(regionList, r)
	}
	sort.Strings(regionList)

	tier, avgLevel := pokeapi.LevelFor(levelRanges)

	encounters := []exploreEncounter{}
	for _, encounter := range location.PokemonEncounters {
		methods := make(map[string]struct{})
		for _, vd := range encounter.VersionDetails {
			for _, ed := range vd.EncounterDetails {
				methods[ed.Method.Name] = struct{}{}
			}
		}
		methodList := make([]string, 0, len(methods))
		for m := range methods {
			methodList = append(methodList, prettify(m))
		}
		sort.Strings(methodList)

		encounters = append(encounters, exploreEncounter{
			Name:    encounter.Pokemon.Name,
			Display: prettify(encounter.Pokemon.Name),
			Methods: methodList,
		})
	}

	data, _ := json.Marshal(exploreResult{
		Name:     areaName,
		Display:  prettify(areaName),
		Tier:     tier,
		Level:    avgLevel,
		Regions:  regionList,
		Habitats: habitatList,
		Pokemon:  encounters,
	})
	return string(data)
}

type battleResult struct {
	Fighter1 string   `json:"fighter1"`
	Fighter2 string   `json:"fighter2"`
	Logs     []string `json:"logs"`
}

func jsBattle(this js.Value, args []js.Value) any {
	name1 := ""
	name2 := ""
	if len(args) > 0 {
		name1 = args[0].String()
	}
	if len(args) > 1 {
		name2 = args[1].String()
	}

	if name1 == name2 {
		return `{"error": "please provide two different pokemon"}`
	}

	caught1, ok := cfg.Pokedex[name1]
	if !ok {
		return `{"error": "you have not caught ` + name1 + `"}`
	}

	caught2, ok := cfg.Pokedex[name2]
	if !ok {
		return `{"error": "you have not caught ` + name2 + `"}`
	}

	battleMon1 := combat.NewBattleMon(caught1)
	battleMon2 := combat.NewBattleMon(caught2)
	battleCfg := combat.BattleConfig{DamageMultiplier: 1.0}

	logs := []string{}
	logs = append(logs, fmt.Sprintf("%s HP: %d", battleMon1.Name, battleMon1.MaxHP))
	logs = append(logs, fmt.Sprintf("%s HP: %d", battleMon2.Name, battleMon2.MaxHP))

	if battleMon1.AffectionBonus != 1 {
		logs = append(logs, fmt.Sprintf("%s is getting a %v bonus from how much he loves his trainer!", battleMon1.Name, battleMon1.AffectionBonus))
	}
	if battleMon2.AffectionBonus != 1 {
		logs = append(logs, fmt.Sprintf("%s is getting a %v bonus from how much he loves his trainer!", battleMon2.Name, battleMon2.AffectionBonus))
	}

	logs = append(logs, combat.SimulateBattle(battleMon1, battleMon2, battleCfg)...)

	data, _ := json.Marshal(battleResult{
		Fighter1: battleMon1.Name,
		Fighter2: battleMon2.Name,
		Logs:     logs,
	})
	return string(data)
}

type partyMember struct {
	Slot    int    `json:"slot"`
	Name    string `json:"name"`
	Display string `json:"display"`
}

type partyResult struct {
	Party   []partyMember `json:"party"`
	Message string        `json:"message,omitempty"`
}

func buildParty(cfg *config) []partyMember {
	members := []partyMember{}
	for i, mon := range cfg.Party {
		display := prettify(mon)
		if caught, ok := cfg.Pokedex[mon]; ok {
			display = prettify(displayName(caught))
		}
		members = append(members, partyMember{
			Slot:    i + 1,
			Name:    mon,
			Display: display,
		})
	}
	return members
}

func jsParty(this js.Value, args []js.Value) any {
	data, _ := json.Marshal(partyResult{
		Party: buildParty(&cfg),
	})
	return string(data)
}

func jsPartyAdd(this js.Value, args []js.Value) any {
	name := ""
	if len(args) > 0 {
		name = args[0].String()
	}
	err := commandPartyAdd(&cfg, name)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	data, _ := json.Marshal(partyResult{
		Message: prettify(name) + " added to your party!",
		Party:   buildParty(&cfg),
	})
	return string(data)
}

func jsPartyRemove(this js.Value, args []js.Value) any {
	name := ""
	if len(args) > 0 {
		name = args[0].String()
	}
	err := commandPartyRemove(&cfg, name)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	data, _ := json.Marshal(partyResult{
		Message: prettify(name) + " removed from your party.",
		Party:   buildParty(&cfg),
	})
	return string(data)
}

func jsNickname(this js.Value, args []js.Value) any {
	name := ""
	nickname := ""
	if len(args) > 0 {
		name = args[0].String()
	}
	if len(args) > 1 {
		nickname = args[1].String()
	}
	err := commandNickname(&cfg, name, nickname)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	data, _ := json.Marshal(partyResult{
		Message: prettify(name) + " is now known as " + prettify(nickname) + "!",
		Party:   buildParty(&cfg),
	})
	return string(data)
}

func jsSave(this js.Value, args []js.Value) any {
	data := SaveData{
		TrainerName: cfg.TrainerName,
		Pokedex:     cfg.Pokedex,
		Party:       cfg.Party,
	}
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	// return the JSON string to JS, let JS store it
	return string(bytes)
}

func jsLoad(this js.Value, args []js.Value) any {
	// JS passes the localStorage string in as args[0]
	if len(args) == 0 {
		return `{"error": "no save data provided"}`
	}
	var data SaveData
	if err := json.Unmarshal([]byte(args[0].String()), &data); err != nil {
		return `{"error": "corrupted save data"}`
	}
	cfg.TrainerName = data.TrainerName
	cfg.Pokedex = data.Pokedex
	cfg.Party = data.Party
	if cfg.Pokedex == nil {
		cfg.Pokedex = make(map[string]combat.CaughtPokemon)
	}
	if cfg.Party == nil {
		cfg.Party = []string{}
	}
	return `{"message": "Pokedex loaded successfully"}`
}

func jsDelete(this js.Value, args []js.Value) any {
	cfg.Pokedex = make(map[string]combat.CaughtPokemon)
	cfg.Party = []string{}
	cfg.TrainerName = ""
	// signal JS to clear localStorage
	return `{"message": "Pokedex cleaned and ready for the next trainer"}`
}

type inspectResult struct {
	Name           string      `json:"name"`
	Display        string      `json:"display"`
	Nickname       string      `json:"nickname,omitempty"`
	Height         int         `json:"height"`
	Weight         int         `json:"weight"`
	BaseExperience int         `json:"base_experience"`
	BallType       string      `json:"ball_type"`
	Sprite         string      `json:"sprite"`
	Types          []string    `json:"types"`
	Stats          []statEntry `json:"stats"`
}

type statEntry struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func jsInspect(this js.Value, args []js.Value) any {
	name := ""
	if len(args) > 0 {
		name = args[0].String()
	}

	pokemon, ok := cfg.Pokedex[name]
	if !ok {
		return `{"error": "you have not caught that pokemon"}`
	}

	types := []string{}
	for _, t := range pokemon.Pokemon.Types {
		types = append(types, prettify(t.Type.Name))
	}

	stats := []statEntry{}
	for _, s := range pokemon.Pokemon.Stats {
		stats = append(stats, statEntry{
			Name:  prettify(s.Stat.Name),
			Value: s.BaseStat,
		})
	}

	data, _ := json.Marshal(inspectResult{
		Name:           pokemon.Pokemon.Name,
		Display:        prettify(displayName(pokemon)),
		Nickname:       pokemon.Nickname,
		Height:         pokemon.Pokemon.Height,
		Weight:         pokemon.Pokemon.Weight,
		BaseExperience: pokemon.Pokemon.BaseExperience,
		BallType:       pokemon.BallType,
		Sprite:         pokemon.Pokemon.Sprites.FrontDefault,
		Types:          types,
		Stats:          stats,
	})
	return string(data)
}

func initFromSave(saveJSON string) {
	var data SaveData
	if err := json.Unmarshal([]byte(saveJSON), &data); err != nil {
		js.Global().Call("showNameModal")
		return
	}
	cfg.TrainerName = data.TrainerName
	cfg.Pokedex = data.Pokedex
	cfg.Party = data.Party
	if cfg.Pokedex == nil {
		cfg.Pokedex = make(map[string]combat.CaughtPokemon)
	}
	if cfg.Party == nil {
		cfg.Party = []string{}
	}
	js.Global().Call("renderTrainerHeader", cfg.TrainerName)
}

func jsBall(this js.Value, args []js.Value) any {

	ballName := ""
	if len(args) > 0 {
		ballName = args[0].String()
	}

	var goArgs []string
	if ballName != "" {
		goArgs = []string{ballName}
	}

	err := commandBall(&cfg, goArgs...)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	type ballResult struct {
		Ball    string `json:"ball"`
		Message string `json:"message"`
	}
	data, _ := json.Marshal(ballResult{
		Ball:    cfg.CurrentBall,
		Message: "Equipped " + prettify(cfg.CurrentBall) + "!",
	})
	return string(data)
}
