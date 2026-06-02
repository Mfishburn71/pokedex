package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Mfishburn71/pokedex/internal/combat"
)

type SaveData struct {
	TrainerName string                          `json:"trainer_name"`
	Pokedex     map[string]combat.CaughtPokemon `json:"pokedex"`
	Party       []string                        `json:"party"`
}

func saveToFile(filename string, data SaveData) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, bytes, 0644)
}
func loadFromFile(filename string) (SaveData, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return SaveData{}, err
	}

	var data SaveData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return SaveData{}, err
	}

	return data, nil
}

func commandSave(cfg *config, args ...string) error {
	data := SaveData{
		TrainerName: cfg.TrainerName,
		Pokedex:     cfg.Pokedex,
		Party:       cfg.Party,
	}

	if err := saveToFile("save.json", data); err != nil {
		return err
	}
	fmt.Println("Pokedex saved to save.json")
	return nil
}

func commandLoad(cfg *config, args ...string) error {
	data, err := loadFromFile("save.json")
	if err != nil {
		return err
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
	fmt.Println("Pokedex loaded from save.json")
	return nil
}

func commandDelete(_ *config, _ ...string) error {
	err := os.Remove("save.json")
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	fmt.Println("Pokedex cleaned and ready for the next trainer")
	return nil
}
