package usecase

import (
	"errors"
	"goexpert-clima/internal/entity"
	"goexpert-clima/internal/gateway"
	"log"
	"regexp"
)

var (
	ErrInvalidZipCode   = errors.New("invalid zipcode")
	ErrZipCodeNotFound  = errors.New("can not find zipcode")
)

type GetTemperatureUseCase struct {
	ZipCodeGateway gateway.ZipCodeGateway
	WeatherGateway gateway.WeatherGateway
}

func NewGetTemperatureUseCase(
	zipCodeGateway gateway.ZipCodeGateway,
	weatherGateway gateway.WeatherGateway,
) *GetTemperatureUseCase {
	return &GetTemperatureUseCase{
		ZipCodeGateway: zipCodeGateway,
		WeatherGateway: weatherGateway,
	}
}

func (uc *GetTemperatureUseCase) Execute(zipCode string) (*entity.Temperature, error) {
	log.Printf("Executing GetTemperatureUseCase with zipCode: %s", zipCode)

	// Validate zip code format (8 digits)
	if !isValidZipCode(zipCode) {
		log.Printf("Invalid zipcode format: %s", zipCode)
		return nil, ErrInvalidZipCode
	}

	// Get location from zip code
	log.Printf("Calling ZipCodeGateway.GetLocation with zipCode: %s", zipCode)
	location, err := uc.ZipCodeGateway.GetLocation(zipCode)
	if err != nil {
		log.Printf("Error in ZipCodeGateway.GetLocation: %v", err)
		return nil, ErrZipCodeNotFound
	}
	log.Printf("Location found: City=%s, State=%s", location.City, location.State)

	// Get temperature for the location
	log.Printf("Calling WeatherGateway.GetTemperature for city: %s", location.City)
	celsius, err := uc.WeatherGateway.GetTemperature(location.City)
	if err != nil {
		log.Printf("Error in WeatherGateway.GetTemperature: %v", err)
		return nil, err
	}
	log.Printf("Temperature retrieved: %v°C", celsius)

	// Create and return temperature entity with conversions
	temperature := entity.NewTemperature(celsius)
	log.Printf("Temperature entity created: %v°C, %v°F, %v°K", temperature.Celsius, temperature.Fahrenheit, temperature.Kelvin)
	return temperature, nil
}

// isValidZipCode checks if the zip code has exactly 8 digits
func isValidZipCode(zipCode string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, zipCode)
	return match
}
