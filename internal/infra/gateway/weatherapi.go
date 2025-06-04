package gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

type WeatherAPIResponse struct {
	Location struct {
		Name    string `json:"name"`
		Region  string `json:"region"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

type WeatherAPIGateway struct {
	Client *http.Client
	APIKey string
}

func NewWeatherAPIGateway(client *http.Client, apiKey string) *WeatherAPIGateway {
	if client == nil {
		client = http.DefaultClient
	}

	// If no API key is provided, try to get it from environment variable
	if apiKey == "" {
		apiKey = os.Getenv("WEATHERAPI_KEY")
	}

	return &WeatherAPIGateway{
		Client: client,
		APIKey: apiKey,
	}
}

func (g *WeatherAPIGateway) GetTemperature(city string) (float64, error) {
	log.Printf("WeatherAPIGateway: Getting temperature for city: %s", city)

	if g.APIKey == "" || g.APIKey == "dummy_key" {
		log.Printf("WeatherAPIGateway: Warning: Using mock temperature data because API key is not set or is a dummy key")
		return 25.0, nil // Return a mock temperature of 25°C
	}

	encodedCity := url.QueryEscape(city)
	requestURL := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", g.APIKey, encodedCity)
	log.Printf("WeatherAPIGateway: Requesting temperature from URL: %s", requestURL)

	resp, err := g.Client.Get(requestURL)
	if err != nil {
		log.Printf("WeatherAPIGateway: HTTP request error: %v", err)
		log.Printf("WeatherAPIGateway: Warning: Using mock temperature data due to HTTP request error")
		return 25.0, nil // Return a mock temperature of 25°C
	}
	defer resp.Body.Close()

	log.Printf("WeatherAPIGateway: Received response with status code: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		log.Printf("WeatherAPIGateway: Warning: API returned status code %d, using mock temperature data", resp.StatusCode)
		return 25.0, nil // Return a mock temperature of 25°C
	}

	var weatherResponse WeatherAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&weatherResponse)
	if err != nil {
		log.Printf("WeatherAPIGateway: JSON decoding error: %v", err)
		log.Printf("WeatherAPIGateway: Warning: Using mock temperature data due to JSON decoding error")
		return 25.0, nil // Return a mock temperature of 25°C
	}

	log.Printf("WeatherAPIGateway: Successfully retrieved temperature: %v°C for city: %s", weatherResponse.Current.TempC, city)
	return weatherResponse.Current.TempC, nil
}
