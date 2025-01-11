package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
    ErrInvalidZipCode   = errors.New("invalid zip code")
    ErrAPIUnavailable   = errors.New("unable to reach OpenWeatherMap API")
    openWeatherMapAPIURL = "https://api.openweathermap.org/data/2.5/weather"
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

// fetches weather data from OpenWeatherMap
func GetWeatherData(zipCode string) (string, error) {
    apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
    if apiKey == "" {
        return "", errors.New("API key is not set")
    }

    url := fmt.Sprintf("%s?zip=%s&appid=%s&units=imperial", openWeatherMapAPIURL, zipCode, apiKey)

    resp, err := http.Get(url)
    if err != nil {
        return "", ErrAPIUnavailable
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        if resp.StatusCode == http.StatusNotFound {
            return "", ErrInvalidZipCode
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
