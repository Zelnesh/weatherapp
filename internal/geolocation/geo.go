package geo

import (

	"net/http"
	"strings"
	"net"
	"encoding/json"
	"fmt"
)


type GeoResponse struct {
	IP string `json:"ip"`
	City string `json:"city"`
	Continent string `json:"continent"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func GetHostIP(r *http.Request) string {

	var ip string

	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff,",")
		ip = strings.TrimSpace(parts[0])
	} else if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		ip = xrip
	} else {

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {

			return r.RemoteAddr
		} else {

			ip = host
		}
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP != nil {
		if parsedIP.IsLoopback() && parsedIP.To4() == nil {
			ip = "127.0.0.1"
		} else if parsedIP.To4() != nil {
			ip = parsedIP.To4().String()
		} else {
			ip = parsedIP.String()
		}
	}

	return ip

}

func GetGeoFromIP(ip string)(*GeoResponse, error) {

	url := fmt.Sprintf("https://ipwho.is/%s", ip)

	resp, err := http.Get(url)
	if err != nil {
		return nil,err
	}
	defer resp.Body.Close()

	var geolocation GeoResponse

	err = json.NewDecoder(resp.Body).Decode(&geolocation)
	if err != nil {
		return nil, err
	}

	return &geolocation, nil
}
