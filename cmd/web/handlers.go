package main

import (

	"html/template"
	"net/http"
	"strings"
	"unicode/utf8"
	"weatherapp.zelnesh.net/internal/geolocation"
	"weatherapp.zelnesh.net/internal/weatherservice"
	"weatherapp.zelnesh.net/internal/geocoding"
)



type FormDataError struct {
	CountryField string
	CityField string
}


func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {

		w.WriteHeader(404)
		w.Write([]byte("Web page not found..."))
		return
	}

	ip:= geo.GetHostIP(r)
	geolocation, err := geo.GetGeoFromIP(ip)
	if err != nil {
		app.logError.Printf("Failed to retreive geolocation: %v", err)
		http.Error(w, "Internal Server Error...", 500)
	}

	app.logInfo.Printf("IP home page: %v", ip)
	app.logInfo.Printf("Geo home page: %+v", geolocation)

	currentWeather, err := weatherservice.GetCurrentWeather(geolocation.Latitude, geolocation.Longitude)
	if err != nil {
		app.logError.Printf("Failed to retreive current weather: %v", err)
		http.Error(w, "Internal Server Error...", 500)
		currentWeather= &weatherservice.WeatherResponse{
			Current: weatherservice.WeatherCurrent{
				Temperature: -999,
			},
		}

	}

	data := struct {
		Title string
		Weather *weatherservice.WeatherResponse
		City string
		Continent string
	}{
		Title: "Current Atmosphere",
		Weather: currentWeather,
		City: geolocation.City,
		Continent: geolocation.Continent,
	}

	tmpl, err := template.ParseFiles("./ui/html/base.tmpl", "./ui/html/home.tmpl")
	if err != nil {
		app.logError.Printf("Failed to parse templates: %v", err)
		http.Error(w, "Internal Server Error...", http.StatusInternalServerError)
	}

	err = tmpl.ExecuteTemplate(w, "base.tmpl", data)
	if err != nil {
		app.logError.Printf("Failed to render template: %v", err)
		http.Error(w, "Internal Server Error...", http.StatusInternalServerError)
	}


}

func (app *application) consultSpace(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

		case http.MethodGet:
			app.consultSpaceGET(w,r)

		case http.MethodPost:
			app.consultSpacePOST(w,r)

		default:
			http.Error(w, "Method Not Allowed...", http.StatusMethodNotAllowed)
	}

}

func (app *application) consultSpaceGET(w http.ResponseWriter, r *http.Request) {

	data := struct{
		Title string
	}{
		Title: "Consult The Space",
	}

	tmpl, err := template.ParseFiles("./ui/html/base.tmpl", "./ui/html/consultspace.tmpl")
	if err != nil {
		app.logError.Printf("Failed to parse templates: %v", err)
		http.Error(w, "Internal Server Error...", http.StatusInternalServerError)
	}

	err = tmpl.ExecuteTemplate(w, "base.tmpl", data)
	if err != nil {
		app.logError.Printf("Failed to render template: %v", err)
		http.Error(w, "Internal Server Error...", http.StatusInternalServerError)
	}
}

func (app *application) consultSpacePOST(w http.ResponseWriter, r *http.Request) {

	var formError FormDataError
	formError.CountryField = "Please enter country code like GR, US, FR"
	formError.CityField = "Filed Ciry can contain 58 characters max..."

	r.ParseForm()
	city := r.FormValue("city")
	country := strings.TrimSpace(r.FormValue("country"))
	country = strings.ToUpper(country)

	if utf8.RuneCountInString(country) != 2 {
		tmpl, err := template.ParseFiles("./ui/html/weather_result_error_country_field.tmpl")
		if err != nil {
			app.logError.Printf("Failed to parse template: %v", err)
			http.Error(w, "Internal Server Error...", http.StatusInternalServerError)
		}
		err = tmpl.ExecuteTemplate(w, "weather_result_error_country_field.tmpl", formError)
		if err != nil {
			app.logError.Printf("Failed to render template: %v", err)
			http.Error(w, "Internal Server Error...", http.StatusInternalServerError)
		}

		return
	}


	geocodingData, err := geocoding.GetGeocoding(city,country)
	app.logInfo.Printf("Geocoding form: %+v", geocodingData.Results[0])

	if len(geocodingData.Results) == 0{
		wrongCity := "City not found, please enter a valid city..."

		tmpl, err := template.ParseFiles("./ui/html/weather_result_error_city_field.tmpl")
		if err != nil {
			app.logError.Printf("Failed to parse template: %v", err)
			http.Error(w, "Internal Server Error...", http.StatusInternalServerError)
		}
		err = tmpl.ExecuteTemplate(w, "weather_result_error_city_field.tmpl", wrongCity)
		if err != nil {
			app.logError.Printf("Failed to render template: %v", err)
			http.Error(w, "Internal Server Error...", http.StatusInternalServerError)
		}
		return
	}


	currentWeather, err := weatherservice.GetCurrentWeather(geocodingData.Results[0].Latitude, geocodingData.Results[0].Longitude)
	app.logInfo.Printf("CurrentWeather form: %+v", currentWeather)


	data := struct{
		City string
		Country string
		Temperature *weatherservice.WeatherResponse
	}{
		City: geocodingData.Results[0].City,
		Country: geocodingData.Results[0].Country,
		Temperature: currentWeather,
	}

	tmpl, err := template.ParseFiles("./ui/html/weather_result.tmpl")
	if err != nil {
		app.logError.Printf("Failed to parse template: %v", err)
		http.Error(w, "Internal Server Error...", http.StatusInternalServerError)
	}

	err = tmpl.ExecuteTemplate(w, "weather_result.tmpl", data)
	if err != nil {
		app.logError.Printf("Failed to render template: %v", err)
		http.Error(w, "Internal Server Error...", http.StatusInternalServerError)
	}
}
