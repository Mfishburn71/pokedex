package main

import (
	"errors"
	"fmt"

	"github.com/Mfishburn71/pokedex/internal/combat"
)

func commandBattle(cfg *config, args ...string) error {
	if len(args) != 2 {
		return errors.New("you need to give me two pokemon names")
	}

	name1 := args[0]
	name2 := args[1]

	if name1 == name2 {
		return errors.New("please provide two different pokemon")
	}

	caught1, ok := cfg.Pokedex[name1]
	if !ok {
		return fmt.Errorf("you have not caught %s", name1)
	}
	battleMon1 := combat.NewBattleMon(caught1)

	caught2, ok := cfg.Pokedex[name2]
	if !ok {
		return fmt.Errorf("you have not caught %s", name2)
	}
	battleMon2 := combat.NewBattleMon(caught2)
	battleCfg := combat.BattleConfig{
		DamageMultiplier: 1.0,
	}
	var logs []string

	logs = append(logs, fmt.Sprintf("%s HP: %d", battleMon1.Name, battleMon1.MaxHP))
	logs = append(logs, fmt.Sprintf("%s HP: %d", battleMon2.Name, battleMon2.MaxHP))

	if battleMon1.AffectionBonus != 1 {
		logs = append(logs, fmt.Sprintf("%s is getting a %v bonus from how much he loves his trainer!", battleMon1.Name, battleMon1.AffectionBonus))
	}
	if battleMon2.AffectionBonus != 1 {
		logs = append(logs, fmt.Sprintf("%s is getting a %v bonus from how much he loves his trainer!", battleMon2.Name, battleMon2.AffectionBonus))
	}

	logs = append(logs, combat.SimulateBattle(battleMon1, battleMon2, battleCfg)...)

	for _, line := range logs {
		fmt.Println(line)
	}

	return nil
}
