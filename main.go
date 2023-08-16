package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string)(apiConfigData, error){
	bytes, err := os.ReadFile(filename)

	if err != nil {
		return apiConfigData{}, err
	}

	var c apiConfigData
	err = json.Unmarshal(bytes, &c)

	if err != nil {
		return apiConfigData{}, err
	}

	return c, err
}

func query(city string)(weatherData, error){
	apiConfigData, err := loadApiConfig(".apiConfig")
	if err!=nil {
		return weatherData{}, err
	}

	var reqUrl string = "https://api.openweathermap.org/data/2.5/weather?q="+city+"&appid="+apiConfigData.OpenWeatherMapApiKey
	resp, err := http.Get(reqUrl)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()
	
	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}

	d.Main.Kelvin -= 273.15

	return d, err
}

func main(){
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		data, err := query(city)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})
	http.ListenAndServe(":8080", nil)
}