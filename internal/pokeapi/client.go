package pokeapi

import (
	//"encoding/json"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Mfishburn71/pokedex/internal/pokecache"
)

type Client struct {
	cache *pokecache.Cache
}

func NewClient() Client {
	return Client{
		cache: pokecache.NewCache(5 * time.Second),
	}
}

/* Legacy function - Kept for internal backup
func (c *Client) clientRequestHelper(url string) (LocationAreasResp, error) {

	cachedData, ok := c.cache.Get(url)
	if ok {
		var result LocationAreasResp
		if err := json.Unmarshal(cachedData, &result); err != nil {
			return LocationAreasResp{}, err
		}
		return result, nil
	} else {
		res, err := http.Get(url)
		if err != nil {
			return LocationAreasResp{}, fmt.Errorf("error making request: %w", err)
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return LocationAreasResp{}, err
		}
		c.cache.Add(url, data)

		var resp LocationAreasResp
		if err := json.Unmarshal(data, &resp); err != nil {
			return LocationAreasResp{}, err
		}

		return resp, nil
	}

}*/

func (c *Client) getBytes(url string) ([]byte, error) {
	// check cache with url key
	data, ok := c.cache.Get(url)
	if ok {
		return data, nil
	} else {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		if res.StatusCode > 299 {
			return nil, fmt.Errorf("bad status code: %d", res.StatusCode)
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		c.cache.Add(url, data)
		return data, nil
	}
}

func (c *Client) GetPokemon(pokemonName string) (PokemonInfo, error) {
	url := BaseURL + "/pokemon/" + pokemonName
	data, err := c.getBytes(url)
	if err != nil {
		if strings.Contains(err.Error(), "bad status code") {
			return PokemonInfo{}, errors.New("I'm sorry, that's not in my records. Please check your spelling")
		}
		return PokemonInfo{}, err
	}
	var resp PokemonInfo
	if err := json.Unmarshal(data, &resp); err != nil {
		return PokemonInfo{}, err
	}
	return resp, nil
}
