package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mikegmatthews/pokedexcli/internal/pokecache"
)

type APIResource struct {
	Count    int         `json:"count"`
	Next     *string     `json:"next"`
	Previous *string     `json:"previous"`
	Results  []APIResult `json:"results"`
}

type APIResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

var cache = pokecache.NewCache(5 * time.Second)

func GetLocationAreas(url string) (*APIResource, error) {
	var areas APIResource
	var data []byte

	data, ok := cache.Get(url)
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("GET error in GetLocationAreas: %v", err)
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Error reading response: %v", err)
		}

		cache.Add(url, data)
	}

	err := json.Unmarshal(data, &areas)
	if err != nil {
		return nil, fmt.Errorf("Error decoding JSON: %v", err)
	}

	return &areas, nil
}
