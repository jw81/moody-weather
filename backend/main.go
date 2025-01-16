package main

import (
	"fmt"
	"net/http"

	"github.com/jw81/moody-weather/backend/handlers"
	"github.com/jw81/moody-weather/backend/services"
)

func main() {
    http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
        handlers.WeatherHandler(services.GetWeatherData, services.GetOpenAIResponse, w, r)
    })    
    fmt.Println("Starting server on :8080")
    http.ListenAndServe(":8080", nil)
}