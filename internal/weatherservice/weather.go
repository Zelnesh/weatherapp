package weatherservice

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
	Error     bool
}

var weatherCache = map[string]*CacheWeather{}

func cacheKey(latitude, longitude float64) string {
	return fmt.Sprintf("%.2f:%.2f", latitude, longitude)
}

func GetCurrentWeather(latitude, longitude float64) (*WeatherResponse, error) {

	key := cacheKey(latitude, longitude)

	// Check cache first
	if cached, ok := weatherCache[key]; ok {

		// If cache still valid
		if time.Now().Before(cached.ExpiresAt) {

			// If error cached, return stale data if exists
			if cached.Error && cached.Data != nil {
				return cached.Data, nil
			}

			if cached.Error {
				return nil, fmt.Errorf("cached weather error (temporary cooldown)")
			}

			return cached.Data, nil
		}
	}

	// Try to fetch fresh data
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&current=temperature_2m",
		latitude,
		longitude,
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

		// fallback to stale cache
		if cached, ok := weatherCache[key]; ok && cached.Data != nil {
			return cached.Data, nil
		}

		return nil, err
	}
	defer resp.Body.Close()

	// Handle rate limit or API errors
	if resp.StatusCode != http.StatusOK {

		var staleData *WeatherResponse

		if cached, ok := weatherCache[key]; ok {
			staleData = cached.Data
		}

		weatherCache[key] = &CacheWeather{
			Data:      staleData,
			ExpiresAt: time.Now().Add(10 * time.Minute), // longer cooldown
			Error:     true,
		}

		if staleData != nil {
			return staleData, nil
		}

		return nil, fmt.Errorf("weather API bad response: %s", resp.Status)
	}

	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, err
	}

	// Store successful response
	weatherCache[key] = &CacheWeather{
		Data:      &weather,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Error:     false,
	}

	return &weather, nil
}
