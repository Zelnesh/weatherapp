package geocoding

import (

	"net/http"
	"fmt"
	"encoding/json"
)


type GeocodingResponse struct {
	Results []GeocodingResults `json"results"`
}

type GeocodingResults struct {
	City string `json:"name"`
	Country string `json:"country"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func GetGeocoding(city, country string) (*GeocodingResponse, error) {

	url := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&country=%s", city,country)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var geocoding GeocodingResponse

	err = json.NewDecoder(resp.Body).Decode(&geocoding)
	if err != nil {
		return nil, err
	}

	return &geocoding, nil
}
