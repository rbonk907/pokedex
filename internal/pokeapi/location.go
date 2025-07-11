package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetLocation(location string) (LocationAreaDetail, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + location

	// check the cache first
	if bytes, found := c.cache.Get(url); found {
		var locationResp LocationAreaDetail
		if err := json.Unmarshal(bytes, &locationResp); err != nil {
			return LocationAreaDetail{}, fmt.Errorf("error unmarshaling json: %w", err)
		}

		return locationResp, nil
	}

	// if not in the cache perform a new http request and add to the cache
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationAreaDetail{}, fmt.Errorf("error creating http request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationAreaDetail{}, fmt.Errorf("error making http request: %w", err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaDetail{}, fmt.Errorf("failed to read data into []byte from response: %w", err)
	}

	var locationResp LocationAreaDetail
	if err := json.Unmarshal(bytes, &locationResp); err != nil {
		return LocationAreaDetail{}, fmt.Errorf("error unmarshaling json: %w", err)
	}

	c.cache.Add(url, bytes)
	return locationResp, nil
}
