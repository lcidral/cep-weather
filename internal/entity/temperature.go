package entity

// Temperature represents the temperature in different units
type Temperature struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

// NewTemperature creates a new Temperature instance with conversions from Celsius
func NewTemperature(celsius float64) *Temperature {
	return &Temperature{
		Celsius:    celsius,
		Fahrenheit: celsius*1.8 + 32,
		Kelvin:     celsius + 273,
	}
}
