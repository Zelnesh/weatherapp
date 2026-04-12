package weather

import (

	"fmt"
	"net/http"
	"encoding/json"
	"time"
)

type WeatherResponse struct {
	Current WeatherCurrent `json:"current"`
}

type WeatherCurrent struct {
	Temperature float64 `json:"temperature_2m"`
}


func GetCurrentWeather(latitude, longitude float64) (*WeatherResponse, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&current=temperature_2m", latitude, longitude)

	client := &http.Client {
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "weather-app")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("Weather service API bad response code: %s", resp.Status)
		return nil, err
	}

	defer resp.Body.Close()

	var weather WeatherResponse

	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return nil, err
	}

	return &weather, nil
}
