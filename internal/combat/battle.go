package combat

import (
	"fmt"
	"strings"

	"github.com/Mfishburn71/pokedex/internal/pokeapi"
)

func NewBattleMon(p CaughtPokemon) BattleMon {
	var ballAffectionBonus = map[string]float64{
		"pokeball":   1.00,
		"greatball":  1.05,
		"ultraball":  1.10,
		"masterball": 1.25,
	}

	affectionBonus, ok := ballAffectionBonus[p.BallType]
	if !ok {
		affectionBonus = 1.0
	}

	hp := getStat(p.Pokemon.Stats, "hp")

	return BattleMon{
		Name:           prettify(displayName(p)),
		Level:          1,
		MaxHP:          hp,
		CurrentHP:      hp,
		Attack:         getStat(p.Pokemon.Stats, "attack"),
		Defense:        getStat(p.Pokemon.Stats, "defense"),
		Speed:          getStat(p.Pokemon.Stats, "speed"),
		AffectionBonus: affectionBonus,
	}
}

func attack(attacker, defender *BattleMon, battleCfg BattleConfig) string {
	baseDamage := attacker.Attack - defender.Defense
	if baseDamage < 1 {
		baseDamage = 1
	}

	damage := int(float64(baseDamage) * battleCfg.DamageMultiplier * attacker.AffectionBonus)
	if damage < 1 {
		damage = 1
	}

	defender.CurrentHP -= damage
	if defender.CurrentHP < 0 {
		defender.CurrentHP = 0
	}

	return fmt.Sprintf("%s attacked %s for %d damage!", attacker.Name, defender.Name, damage)
}

func SimulateBattle(a, b BattleMon, battleCfg BattleConfig) []string {
	logs := []string{}

	attacker := &a
	defender := &b
	if b.Speed > a.Speed {
		attacker = &b
		defender = &a
	}

	for a.CurrentHP > 0 && b.CurrentHP > 0 {
		logs = append(logs, attack(attacker, defender, battleCfg))
		if defender.CurrentHP == 0 {
			break
		}

		attacker, defender = defender, attacker
	}

	winner := a
	if b.CurrentHP > 0 {
		winner = b
	}

	logs = append(logs, fmt.Sprintf("The winner is %s", winner.Name))
	return logs
}

func getStat(stats []pokeapi.PokemonStat, name string) int {
	for _, s := range stats {
		if s.Stat.Name == name {
			return s.BaseStat
		}
	}
	return 1
}

func displayName(caught CaughtPokemon) string {
	if caught.Nickname != "" {
		return fmt.Sprintf("%s (%s)", caught.Nickname, prettify(caught.Pokemon.Name))
	}
	return prettify(caught.Pokemon.Name)
}
func prettify(s string) string {
	parts := strings.Split(s, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}
