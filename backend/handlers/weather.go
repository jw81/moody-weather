package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jw81/moody-weather/backend/validation"
)

type WeatherRequest struct {
	ZipCode string `json:"zipCode"`
	Tone    string `json:"tone"`
}

type WeatherResponse struct {
	Message string `json:"message"`
}

func WeatherHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	response := WeatherResponse{
		Message: fmt.Sprintf("Weather data for zip code %s with tone %s.", weatherRequest.ZipCode, weatherRequest.Tone),
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(response)
}