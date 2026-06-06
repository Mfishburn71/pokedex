//go:build js && wasm

package pokeapi

import (
	"errors"
	"syscall/js"
)

func (c *Client) getBytes(url string) ([]byte, error) {
	data, ok := c.cache.Get(url)
	if ok {
		return data, nil
	}

	result := js.Global().Call("syncFetch", url)
	if result.Get("error").Truthy() {
		return nil, errors.New(result.Get("error").String())
	}

	buf := js.Global().Get("Uint8Array").New(result.Get("data"))
	bytes := make([]byte, buf.Get("length").Int())
	js.CopyBytesToGo(bytes, buf)
	c.cache.Add(url, bytes)
	return bytes, nil
}
