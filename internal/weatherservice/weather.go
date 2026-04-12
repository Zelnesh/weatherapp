package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WeatherResponse struct {
	Current WeatherCurrent `json:"current"`
}

type WeatherCurrent struct {
	Temperature float64 `json:"temperature_2m"`
}

type CacheWeather struct {
	Data      *WeatherResponse
	ExpiresAt time.Time
}

// key = "lat:lon"
var weatherCache = map[string]*CacheWeather{}

func cacheKey(latitude, longitude float64) string {
	return fmt.Sprintf("%.2f:%.2f", latitude, longitude)
}

func GetCurrentWeather(latitude, longitude float64) (*WeatherResponse, error) {

	key := cacheKey(latitude, longitude)

	// 1. Check cache
	if cached, ok := weatherCache[key]; ok {
		if time.Now().Before(cached.ExpiresAt) {
			return cached.Data, nil
		}
	}

	// 2. Build request
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&current=temperature_2m",
		latitude, longitude,
	)

	client := &http.Client{
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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API bad response: %s", resp.Status)
	}

	// 3. Decode response
	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, err
	}

	// 4. Store in cache (per location)
	weatherCache[key] = &CacheWeather{
		Data:      &weather,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	return &weather, nil
}
