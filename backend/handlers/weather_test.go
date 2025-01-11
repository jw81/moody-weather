package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock function to simulate OpenWeatherMap API calls
var mockGetWeatherData = func(zipCode string) (string, error) {
    if zipCode == "12345" {
        return "Sunny (clear sky), 72.0°F (feels like 70.0°F), Wind: 5.0 mph in Springfield", nil
    }
    if zipCode == "00000" {
        return "", errInvalidZipCode
    }
    return "", errAPIUnavailable
}

func TestWeatherHandlerWithAPI(t *testing.T) {
    tests := []struct {
        name           string
        requestBody    map[string]string
        expectedStatus int
        expectedBody   string
    }{
        {
            name: "Valid request with successful API response",
            requestBody: map[string]string{
                "zipCode": "12345",
                "tone":    "nice",
            },
            expectedStatus: http.StatusOK,
            expectedBody:   "Weather for 12345: Sunny (clear sky), 72.0°F (feels like 70.0°F), Wind: 5.0 mph in Springfield",
        },
        {
            name:           "Invalid zip code",
            requestBody:    map[string]string{"zipCode": "00000", "tone": "nice"},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   "Invalid zip code",
        },
        {
            name:           "API unavailable",
            requestBody:    map[string]string{"zipCode": "54321", "tone": "nice"},
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   "Failed to fetch weather data",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Prepare the request body
            requestBody, _ := json.Marshal(tt.requestBody)

            // Create an HTTP request and response recorder
            req := httptest.NewRequest(http.MethodPost, "/weather", bytes.NewBuffer(requestBody))
            req.Header.Set("Content-Type", "application/json")
            rec := httptest.NewRecorder()

            // Replace the real API call with the mock
            getWeatherData = mockGetWeatherData

            // Call the handler
            WeatherHandler(rec, req)

            // Validate the response
            res := rec.Result()
            defer res.Body.Close()

            if res.StatusCode != tt.expectedStatus {
                t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
            }

            body := rec.Body.String()
            if tt.expectedBody != "" && !containsString(body, tt.expectedBody) {
                t.Errorf("expected body to contain %q, got %q", tt.expectedBody, body)
            }
        })
    }
}

func containsString(body, expected string) bool {
    return bytes.Contains([]byte(body), []byte(expected))
}