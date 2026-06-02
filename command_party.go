package main

import (
	"errors"
	"fmt"
	"strings"
)

func commandParty(cfg *config, _ ...string) error {
	if len(cfg.Party) == 0 {
		fmt.Println("You don't have a party yet!")
		return nil
	}

	for index, mon := range cfg.Party {
		caught, ok := cfg.Pokedex[mon]
		if ok { ///Check for a nickname
			fmt.Printf("Slot #%d: %s\n", index+1, prettify(displayName(caught)))
		} else {
			fmt.Printf("Slot #%d: %s\n", index+1, prettify(mon))
		}
	}

	return nil
}

func commandPartyAdd(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("you need to give me a pokemon name")
	}
	if len(args) > 1 {
		return errors.New("please only provide one pokemon at a time")
	}
	if len(cfg.Party) >= 6 {
		return errors.New("Party is full! You can only hold 6 Pokemon at a time!")
	}

	name := strings.ToLower(args[0])
	_, ok := cfg.Pokedex[name]
	if !ok {
		_, err := cfg.PokeAPIClient.GetPokemon(name)
		if err != nil {
			return err
		}
		return errors.New("you have not caught that pokemon")
	}
	for _, member := range cfg.Party {
		if member == name {
			return errors.New("That pokemon is already in your party")
		}
	}

	cfg.Party = append(cfg.Party, name)
	fmt.Printf("Added %s to your party\n", prettify(name))
	for index, mon := range cfg.Party {
		caught, ok := cfg.Pokedex[mon]
		if ok { ///Check for a nickname
			fmt.Printf("Slot #%d: %s\n", index+1, prettify(displayName(caught)))
		} else {
			fmt.Printf("Slot #%d: %s\n", index+1, prettify(mon))
		}
	}
	return nil
}

func commandPartyRemove(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("you need to give me a pokemon name")
	}
	if len(args) > 1 {
		return errors.New("please only provide one pokemon at a time")
	}
	if len(cfg.Party) == 0 {
		return errors.New("Party is empty! You should add a companion!")
	}
	name := strings.ToLower(args[0])
	_, ok := cfg.Pokedex[name]
	if !ok {
		_, err := cfg.PokeAPIClient.GetPokemon(name)
		if err != nil {
			return err
		}
		return errors.New("you have not caught that pokemon")
	}
	for i, member := range cfg.Party {
		if member == name {
			cfg.Party = append(cfg.Party[:i], cfg.Party[i+1:]...)
			fmt.Printf("Removed %s from your party\n", prettify(name))
			for index, mon := range cfg.Party {
				caught, ok := cfg.Pokedex[mon]
				if ok { ///Check for a nickname
					fmt.Printf("Slot #%d: %s\n", index+1, prettify(displayName(caught)))
				} else {
					fmt.Printf("Slot #%d: %s\n", index+1, prettify(mon))
				}
			}
			return nil
		}
	}
	return errors.New("That pokemon is not in your party")
}

func commandNickname(cfg *config, args ...string) error {
	if len(args) < 2 {
		return errors.New("Give me the pokemon name first, then your nickname")
	}
	if len(args) > 2 {
		return errors.New("please only provide one pokemon at a time")
	}
	name := strings.ToLower(args[0])
	_, ok := cfg.Pokedex[name]
	if !ok {
		_, err := cfg.PokeAPIClient.GetPokemon(name)
		if err != nil {
			return err
		}
		return errors.New("you have not caught that pokemon")
	}
	caught := cfg.Pokedex[name]
	caught.Nickname = args[1]
	cfg.Pokedex[name] = caught

	fmt.Printf("%s's new name is %s.\nI think they like it!\n", prettify(name), prettify(caught.Nickname))
	return nil
}
