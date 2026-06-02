package pokeapi

type LocationAreasResp struct {
	Count    int                   `json:"count"`
	Next     *string               `json:"next"`
	Previous *string               `json:"previous"`
	Results  []LocationAreaSummary `json:"results"`
}

type LocationAreaSummary struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationArea struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Method struct {
					Name string `json:"name"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
				MaxLevel int `json:"max_level"`
			} `json:"encounter_details"`
			Version struct {
				Name string `json:"name"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
	EncounterMethodRates []struct {
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
}

type PokemonInfo struct {
	Name           string        `json:"name"`
	BaseExperience int           `json:"base_experience"`
	ID             int           `json:"id"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Stats          []PokemonStat `json:"stats"`
	Types          []PokemonType `json:"types"`
}

type PokemonStat struct {
	BaseStat int      `json:"base_stat"`
	Stat     StatInfo `json:"stat"`
}

type StatInfo struct {
	Name string `json:"name"`
}

type PokemonType struct {
	Type TypeInfo `json:"type"`
}

type TypeInfo struct {
	Name string `json:"name"`
}
