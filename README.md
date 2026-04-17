# 🐙 Lovecraftian Weather App

A dark, atmospheric weather application inspired by Lovecraftian cosmic horror.  
Check the weather while gazing into the abyss... and sometimes, the abyss gazes back.

## 🌩️ Features

- Current weather conditions
- City & Country search
- IP-based geolocation (City & Continent)
- Lovecraftian-themed UI
- HTMX-powered dynamic updates
- Server-side validation
- Security headers middleware
- Lightweight and fast (written in Go)

## 🧠 Tech Stack

Backend

Go (Golang)
Go HTML Templates
HTTP Server (net/http)
______________________
Frontend

HTML Templates
HTMX
CSS (Custom Lovecraftian styling)
______________________
Performance & Concurrency

In-memory caching
sync.RWMutex (thread-safe cache access)
singleflight (request deduplication)
Rate-limit protection logic
______________________
External APIs

WeatherAPI
Geocoding API
IPWho.is (Geolocation)


