package pokeapi

import ("fmt"
        "encoding/json"
				"net/http"
				"errors"
				"time"
				"io"
				"pokedex/internal/pokecache"
        )

type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ResourceList struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string    `json:"previous"`
	Results  []NamedAPIResource `json:"results"`
}

var cache = pokecache.NewCache(5 * time.Second)

func GetAPIPage(url string) (ResourceList, error) {
	if stored, exists := cache.Get(url); exists {
		var result ResourceList
		if err := json.Unmarshal(stored, &result); err != nil {
			errString := fmt.Sprintf("Failed to decode json: %v", err)
			return ResourceList{}, errors.New(errString)
		}
		return result, nil
	}

	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errString := fmt.Sprintf("Failed to create request: %v", err)
		return ResourceList{}, errors.New(errString)
	}

	res, err := client.Do(req)
	if err != nil {
		errString := fmt.Sprintf("Failed to get request: %v", err)
		return ResourceList{}, errors.New(errString)
	}

	defer res.Body.Close()

	if statusCode := res.StatusCode; statusCode < 200 || statusCode > 299 {
		errString := fmt.Sprintf("Error Status Code: %v", statusCode)
		return ResourceList{}, errors.New(errString)
	}
	
	var result ResourceList
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		errString := fmt.Sprintf("Failed to convert request to []byte: %v", err)
		return ResourceList{}, errors.New(errString)
	}
	
	if err := json.Unmarshal(bytes, &result); err != nil {
		errString := fmt.Sprintf("Failed to decode json: %v", err)
		return ResourceList{}, errors.New(errString)
	}

	cache.Add(url, bytes)
	return result, nil
}



type AreaInfo struct {
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


func GetLocationArea(areaName string) (AreaInfo, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + areaName

	if stored, exists := cache.Get(url); exists {
		var result AreaInfo
		if err := json.Unmarshal(stored, &result); err != nil {
			errString := fmt.Sprintf("Failed to decode json: %v", err)
			return AreaInfo{}, errors.New(errString)
		}
		return result, nil
	}

	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errString := fmt.Sprintf("Failed to create request: %v", err)
		return AreaInfo{}, errors.New(errString)
	}

	res, err := client.Do(req)
	if err != nil {
		errString := fmt.Sprintf("Failed to get request: %v", err)
		return AreaInfo{}, errors.New(errString)
	}

	defer res.Body.Close()

	if statusCode := res.StatusCode; statusCode < 200 || statusCode > 299 {
		errString := fmt.Sprintf("Error Status Code: %v", statusCode)
		return AreaInfo{}, errors.New(errString)
	}
	
	var result AreaInfo
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		errString := fmt.Sprintf("Failed to convert request to []byte: %v", err)
		return AreaInfo{}, errors.New(errString)
	}
	
	if err := json.Unmarshal(bytes, &result); err != nil {
		errString := fmt.Sprintf("Failed to decode json: %v", err)
		return AreaInfo{}, errors.New(errString)
	}

	cache.Add(url, bytes)
	return result, nil
}



