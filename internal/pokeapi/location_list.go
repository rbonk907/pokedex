package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetLocationList(url string) (LocationAreaList, error) {
	// check the cache first
	if bytes, found := c.cache.Get(url); found {
		var locationList LocationAreaList
		if err := json.Unmarshal(bytes, &locationList); err != nil {
			return LocationAreaList{}, fmt.Errorf("error unmarshaling json: %w", err)
		}

		return locationList, nil
	}

	// if not in the cache perform a new http request and add to the cache
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationAreaList{}, fmt.Errorf("error creating http request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationAreaList{}, fmt.Errorf("error making http request: %w", err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaList{}, fmt.Errorf("failed to read data into []byte from response: %w", err)
	}

	var locationList LocationAreaList
	if err := json.Unmarshal(bytes, &locationList); err != nil {
		return LocationAreaList{}, fmt.Errorf("error unmarshaling json: %w", err)
	}

	c.cache.Add(url, bytes)
	return locationList, nil
}
