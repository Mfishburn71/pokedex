package pokeapi

import (
	"encoding/json"
	"fmt"
)

const BaseURL = "https://pokeapi.co/api/v2"

func (c *Client) GetLocationArea(areaName string) (LocationArea, error) {
	// 1. Build the full URL with areaName appended
	url := BaseURL + "/location-area/" + areaName
	// 2. Check the cache for that URL
	//    - If hit: unmarshal cached bytes into LocationArea and return
	data, err := c.getBytes(url)
	if err != nil {
		return LocationArea{}, fmt.Errorf("error making request: %w", err)
	}
	var resp LocationArea
	if err := json.Unmarshal(data, &resp); err != nil {
		return LocationArea{}, err
	}
	return resp, nil
}
