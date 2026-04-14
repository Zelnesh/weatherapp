package weatherservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type WeatherResponse struct {
	Current WeatherCurrent `json:"current"`
}

type WeatherCurrent struct {
	Temperature float64 `json:"temp_c"`
}

type CacheWeather struct {
	Data      *WeatherResponse
	ExpiresAt time.Time
	Error     bool
}

var (
	weatherCache = map[string]*CacheWeather{}
	mu           sync.RWMutex
	sf           singleflight.Group
)

func cacheKey(latitude, longitude float64) string {
	return fmt.Sprintf("%.2f:%.2f", latitude, longitude)
}

// Retry with exponential backoff
func fetchWithRetry(req *http.Request, client *http.Client) (*http.Response, error) {

	var resp *http.Response
	var err error

	for i := 0; i < 3; i++ {

		resp, err = client.Do(req)

		if err == nil && resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}

		time.Sleep(time.Duration(1<<i) * time.Second)
	}

	return resp, err
}

func GetCurrentWeather(latitude, longitude float64) (*WeatherResponse, error) {

	// Fix invalid coords (Render health checks)
	if latitude == 0 && longitude == 0 {
		latitude = 52.406822
		longitude = -1.519693
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("missing WEATHER_API_KEY")
	}

	key := cacheKey(latitude, longitude)

	// ---------- CACHE READ ----------
	mu.RLock()
	cached, ok := weatherCache[key]
	mu.RUnlock()

	if ok && time.Now().Before(cached.ExpiresAt) {

		if cached.Error && cached.Data != nil {
			return cached.Data, nil
		}

		if cached.Error {
			return nil, fmt.Errorf("cached weather error (cooldown active)")
		}

		return cached.Data, nil
	}

	// ---------- SINGLEFLIGHT ----------
	result, err, _ := sf.Do(key, func() (interface{}, error) {

		// Double check cache inside singleflight
		mu.RLock()
		cached, ok := weatherCache[key]
		mu.RUnlock()

		if ok && time.Now().Before(cached.ExpiresAt) {
			return cached.Data, nil
		}

		url := fmt.Sprintf(
			"https://api.weatherapi.com/v1/current.json?key=%s&q=%.4f,%.4f",
			apiKey,
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

		resp, err := fetchWithRetry(req, client)
		if err != nil {

			// fallback to stale cache
			mu.RLock()
			cached, ok := weatherCache[key]
			mu.RUnlock()

			if ok && cached.Data != nil {
				return cached.Data, nil
			}

			return nil, err
		}

		defer resp.Body.Close()

		// ---------- ERROR HANDLING ----------
		if resp.StatusCode != http.StatusOK {

			var stale *WeatherResponse

			mu.RLock()
			if c, ok := weatherCache[key]; ok {
				stale = c.Data
			}
			mu.RUnlock()

			mu.Lock()
			weatherCache[key] = &CacheWeather{
				Data:      stale,
			 ExpiresAt: time.Now().Add(5 * time.Minute),
				Error:     true,
			}
			mu.Unlock()

			if stale != nil {
				return stale, nil
			}

			return nil, fmt.Errorf("weather API bad response: %s", resp.Status)
		}

		// ---------- SUCCESS ----------
		var weather WeatherResponse

		if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
			return nil, err
		}

		mu.Lock()
		weatherCache[key] = &CacheWeather{
			Data:      &weather,
			 ExpiresAt: time.Now().Add(10 * time.Minute),
				Error:     false,
		}
		mu.Unlock()

		return &weather, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*WeatherResponse), nil
}

// Warmup on startup
func init() {
	go func() {
		time.Sleep(15 * time.Second)
		_, _ = GetCurrentWeather(52.406822, -1.519693)
	}()
}
