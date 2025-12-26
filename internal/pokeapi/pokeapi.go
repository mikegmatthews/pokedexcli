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

type LocationArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

var cache *pokecache.Cache

func init() {
	fmt.Println("Initializing PokeApi...")
	cache = pokecache.NewCache(5 * time.Second)
}

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

func GetPokemonEncounters(areaName string) ([]string, error) {
	var pokemon []string

	urlFmt := "https://pokeapi.co/api/v2/location-area/%s/"
	url := fmt.Sprintf(urlFmt, areaName)

	data, ok := cache.Get(url)
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("GET error in GetPokemonEncounters: %v", err)
		}
		defer resp.Body.Close()

		laData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Error reading response: %v", err)
		}

		var locArea LocationArea
		err = json.Unmarshal(laData, &locArea)
		if err != nil {
			return nil, fmt.Errorf("Error getting LocationArea: %s", laData)
		}

		for _, enc := range locArea.PokemonEncounters {
			pokemon = append(pokemon, enc.Pokemon.Name)
		}

		data, err = json.Marshal(pokemon)
		if err != nil {
			return nil, fmt.Errorf("Error encoding JSON data: %v", err)
		}
		cache.Add(url, data)
	} else {
		err := json.Unmarshal(data, &pokemon)
		if err != nil {
			return nil, fmt.Errorf("Error decoding JSON: %v", err)
		}
	}

	return pokemon, nil
}
