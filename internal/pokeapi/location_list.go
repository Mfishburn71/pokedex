package pokeapi

import (
	"encoding/json"
	"fmt"
	//"io"
	//"net/http"
)

func (c *Client) ListLocations(url string) (LocationAreasResp, error) {
	//res, err := http.Get(url)
	data, err := c.getBytes(url)
	if err != nil {
		return LocationAreasResp{}, fmt.Errorf("error making request: %w", err)
	}

	/*
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return LocationAreasResp{}, err
		}
	*/
	var resp LocationAreasResp
	if err := json.Unmarshal(data, &resp); err != nil {
		return LocationAreasResp{}, err
	}

	return resp, nil
}
