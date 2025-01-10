package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeatherHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid request",
			requestBody: map[string]string{
				"zipCode": "12345",
				"tone":    "nice",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Weather data for zip code 12345 with tone nice.",
		},
		{
			name:           "Missing zipCode",
			requestBody:    map[string]string{"tone": "nice"},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid zip code",
		},
		{
			name:           "Invalid tone",
			requestBody:    map[string]string{"zipCode": "12345", "tone": "angry"},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid tone. Valid options are: nice, normal, snarky.",
		},
		{
			name:           "Empty request body",
			requestBody:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare the request body
			var requestBody []byte
			if tt.requestBody != nil {
				requestBody, _ = json.Marshal(tt.requestBody)
			}

			// Create an HTTP request and response recorder
			req := httptest.NewRequest(http.MethodPost, "/weather", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

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
