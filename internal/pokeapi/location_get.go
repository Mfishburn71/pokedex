package pokeapi

import (
	"encoding/json"
	"fmt"
	//"io"
	//"net/http"
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
	// 3. On miss: make the HTTP request
	// 4. Read the body
	var resp LocationArea
	if err := json.Unmarshal(data, &resp); err != nil {
		return LocationArea{}, err
	}
	// 5. Add body to cache
	// 6. Unmarshal into LocationArea
	// 7. Return it
	return resp, nil
}
