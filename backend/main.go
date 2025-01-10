package main

import (
	"fmt"
	"net/http"

	"github.com/jw81/moody-weather/backend/handlers"
)

func main() {
    http.HandleFunc("/weather", handlers.WeatherHandler)
    fmt.Println("Starting server on :8080")
    http.ListenAndServe(":8080", nil)
}