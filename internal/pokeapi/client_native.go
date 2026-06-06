//go:build !js

package pokeapi

import (
	"fmt"
	"io"
	"net/http"
)

func (c *Client) getBytes(url string) ([]byte, error) {
	data, ok := c.cache.Get(url)
	if ok {
		return data, nil
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("bad status code: %d", res.StatusCode)
	}
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	c.cache.Add(url, data)
	return data, nil
}
