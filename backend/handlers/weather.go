package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jw81/moody-weather/backend/services"
	"github.com/jw81/moody-weather/backend/validation"
)

type WeatherRequest struct {
    ZipCode string `json:"zipCode"`
    Tone    string `json:"tone"`
}

type WeatherResponse struct {
    Message string `json:"message"`
}

// WeatherHandler handles incoming weather requests
func WeatherHandler(getWeatherData func(string) (string, error), responseWriter http.ResponseWriter, request *http.Request) {
    if request.Method != http.MethodPost {
        http.Error(responseWriter, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var weatherRequest WeatherRequest
    err := json.NewDecoder(request.Body).Decode(&weatherRequest)
    if err != nil {
        http.Error(responseWriter, "Invalid request body", http.StatusBadRequest)
        return
    }

    if !validation.IsValidZipCode(weatherRequest.ZipCode) {
        http.Error(responseWriter, "Invalid zip code", http.StatusBadRequest)
        return
    }
    if !validation.IsValidTone(weatherRequest.Tone) {
        http.Error(responseWriter, "Invalid tone. Valid options are: nice, normal, snarky.", http.StatusBadRequest)
        return
    }

    weatherData, err := getWeatherData(weatherRequest.ZipCode)
    if err != nil {
        if errors.Is(err, services.ErrInvalidZipCode) {
            http.Error(responseWriter, "Invalid zip code", http.StatusBadRequest)
            return
        }
        http.Error(responseWriter, "Failed to fetch weather data", http.StatusInternalServerError)
        return
    }

    response := WeatherResponse{
        Message: fmt.Sprintf("Weather for %s: %s", weatherRequest.ZipCode, weatherData),
    }

    responseWriter.Header().Set("Content-Type", "application/json")
    json.NewEncoder(responseWriter).Encode(response)
}
