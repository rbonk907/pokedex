package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetPokemon(name string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + name

	if bytes, found := c.cache.Get(url); found {
		var pokemonResp Pokemon
		if err := json.Unmarshal(bytes, &pokemonResp); err != nil {
			return Pokemon{}, fmt.Errorf("error unmarshaling json: %w", err)
		}

		return pokemonResp, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Pokemon{}, fmt.Errorf("error creating http request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Pokemon{}, fmt.Errorf("error making http request: %w", err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, fmt.Errorf("failed to read data into []byte from response: %w", err)
	}

	var pokemonResp Pokemon
	if err := json.Unmarshal(bytes, &pokemonResp); err != nil {
		return Pokemon{}, fmt.Errorf("error unmarshaling json: %w", err)
	}

	c.cache.Add(url, bytes)
	return pokemonResp, nil
}
