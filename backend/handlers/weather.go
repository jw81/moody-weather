package handlers

import (
	"encoding/json"
	"errors"
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

func WeatherHandler(
    getWeatherData func(string) (string, error),
    getOpenAIResponse func(string, string) (string, error),
    responseWriter http.ResponseWriter,
    request *http.Request,
) {
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

    openAIResponse, err := getOpenAIResponse(weatherData, weatherRequest.Tone)
    if err != nil {
        http.Error(responseWriter, "Failed to generate response", http.StatusInternalServerError)
        return
    }

    response := WeatherResponse{
        Message: openAIResponse,
    }

    responseWriter.Header().Set("Content-Type", "application/json")
    json.NewEncoder(responseWriter).Encode(response)
}
