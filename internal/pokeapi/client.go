package pokeapi

import ("fmt"
        "encoding/json"
				"net/http"
				"errors"
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


func GetAPIPage(url string) (ResourceList, error) {
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
		errString := fmt.Sprintf("Error Status Code: v%", statusCode)
		return ResourceList{}, errors.New(errString)
	}
	
	var result ResourceList
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&result); err != nil {
		errString := fmt.Sprintf("Failed to decode json: %v")
		return ResourceList{}, errors.New(errString)
	}

	return result, nil
}
