package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jw81/moody-weather/backend/validation"
)

type WeatherRequest struct {
    ZipCode string `json:"zipCode"`
    Tone    string `json:"tone"`
}

type WeatherResponse struct {
    Message string `json:"message"`
}

var (
    errInvalidZipCode = errors.New("invalid zip code")
    errAPIUnavailable = errors.New("unable to reach OpenWeatherMap API")
)

// OpenWeatherMap API response structure
type APIWeatherResponse struct {
    Weather []struct {
        Main        string `json:"main"`
        Description string `json:"description"`
    } `json:"weather"`
    Main struct {
        Temp      float64 `json:"temp"`
        FeelsLike float64 `json:"feels_like"`
    } `json:"main"`
    Wind struct {
        Speed float64 `json:"speed"`
    } `json:"wind"`
    Name string `json:"name"`
}

// WeatherHandler handles incoming weather requests
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

    weatherData, err := getWeatherData(weatherRequest.ZipCode)
    if err != nil {
        if errors.Is(err, errInvalidZipCode) {
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

var getWeatherData = func(zipCode string) (string, error) {
    apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
    if apiKey == "" {
        return "", errors.New("API key is not set")
    }

    url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?zip=%s&appid=%s&units=imperial", zipCode, apiKey)

    resp, err := http.Get(url)
    if err != nil {
        return "", errAPIUnavailable
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        if resp.StatusCode == http.StatusNotFound {
            return "", errInvalidZipCode
        }
        return "", fmt.Errorf("unexpected API response status: %d", resp.StatusCode)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    var weatherData APIWeatherResponse
    err = json.Unmarshal(body, &weatherData)
    if err != nil {
        return "", err
    }

    if len(weatherData.Weather) == 0 {
        return "", errors.New("missing weather data in response")
    }

    condition := weatherData.Weather[0].Main
    description := weatherData.Weather[0].Description
    temperature := weatherData.Main.Temp
    feelsLike := weatherData.Main.FeelsLike
    windSpeed := weatherData.Wind.Speed
    city := weatherData.Name

    return fmt.Sprintf("%s (%s), %.1f°F (feels like %.1f°F), Wind: %.1f mph in %s", condition, description, temperature, feelsLike, windSpeed, city), nil
}
