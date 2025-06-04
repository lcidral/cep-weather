package gateway

// ZipCodeGateway defines the interface for CEP (zip code) lookup services
type ZipCodeGateway interface {
	GetLocation(zipCode string) (*Location, error)
}

// WeatherGateway defines the interface for weather services
type WeatherGateway interface {
	GetTemperature(city string) (float64, error)
}

// Location represents a location returned by the CEP service
type Location struct {
	City string
	State string
}
