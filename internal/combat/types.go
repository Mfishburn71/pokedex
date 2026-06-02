package combat

import "github.com/Mfishburn71/pokedex/internal/pokeapi"

type BattleMon struct {
	Name           string
	MaxHP          int
	CurrentHP      int
	Attack         int
	Defense        int
	Speed          int
	AffectionBonus float64
	//TODO: Implement Level, Attack and Defense bonus if transferring to full game
	Level        int
	AttackBonus  int
	DefenseBonus int
}

type BattleConfig struct {
	DamageMultiplier float64
	//TODO: Crit Chance, Level Scale, ad Type effectiveness for full game
	CritChance        float64
	LevelScale        bool
	TypeEffectiveness bool
	BallBonus         bool
}

type CaughtPokemon struct {
	Pokemon  pokeapi.PokemonInfo `json:"pokemon"`
	BallType string              `json:"ball_type"`
	Level    int                 `json:"level"`
	Nickname string              `json:"nickname"`
}
