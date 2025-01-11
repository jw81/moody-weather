package services

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWeatherData(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Query().Get("zip") == "12345" {
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(`{
                "weather": [{"main": "Clouds", "description": "broken clouds"}],
                "main": {"temp": 65.0, "feels_like": 63.0},
                "wind": {"speed": 7.0},
                "name": "Springfield"
            }`))
        } else if r.URL.Query().Get("zip") == "00000" {
            w.WriteHeader(http.StatusNotFound)
        } else {
            w.WriteHeader(http.StatusInternalServerError)
        }
    }))
    defer server.Close()

    oldAPIURL := openWeatherMapAPIURL
    openWeatherMapAPIURL = server.URL
    defer func() { openWeatherMapAPIURL = oldAPIURL }()

    tests := []struct {
        zipCode        string
        expectedResult string
        expectError    bool
    }{
        {"12345", "Clouds (broken clouds), 65.0°F (feels like 63.0°F), Wind: 7.0 mph in Springfield", false},
        {"00000", "", true},
        {"54321", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.zipCode, func(t *testing.T) {
            t.Logf("Running test for zip code: %s", tt.zipCode)
            result, err := GetWeatherData(tt.zipCode)
            t.Logf("Result: %q, Error: %v", result, err)
            if tt.expectError && err == nil {
                t.Errorf("expected an error but got none")
            }
            if !tt.expectError && result != tt.expectedResult {
                t.Errorf("expected %q, got %q", tt.expectedResult, result)
            }
        })
    }
}

