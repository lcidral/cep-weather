package main

import (
	"fmt"
	infraGateway "goexpert-clima/internal/infra/gateway"
	"goexpert-clima/internal/infra/web/handlers"
	"goexpert-clima/internal/usecase"
	"log"
	"net/http"
	"os"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get WeatherAPI key from environment
	weatherAPIKey := os.Getenv("WEATHERAPI_KEY")
	if weatherAPIKey == "" {
		log.Println("WARNING: WEATHERAPI_KEY environment variable is not set. Some functionality may be limited.")
		weatherAPIKey = "dummy_key" // Provide a dummy key to allow the application to start
	}

	// Create HTTP client
	httpClient := http.DefaultClient

	// Create gateways
	zipCodeGateway := infraGateway.NewViaCEPGateway(httpClient)
	weatherGateway := infraGateway.NewWeatherAPIGateway(httpClient, weatherAPIKey)

	// Create use case
	getTemperatureUseCase := usecase.NewGetTemperatureUseCase(zipCodeGateway, weatherGateway)

	// Create handler
	temperatureHandler := handlers.NewTemperatureHandler(getTemperatureUseCase)

	// Set up HTTP routes
	http.HandleFunc("/temperature", temperatureHandler.GetTemperature)

	// Start server
	fmt.Printf("Server running on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
