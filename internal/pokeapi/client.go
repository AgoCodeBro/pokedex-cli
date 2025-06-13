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
