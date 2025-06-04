package entity

import (
	"testing"
)

func TestNewTemperature(t *testing.T) {
	tests := []struct {
		name           string
		celsius        float64
		wantCelsius    float64
		wantFahrenheit float64
		wantKelvin     float64
	}{
		{
			name:           "zero celsius",
			celsius:        0,
			wantCelsius:    0,
			wantFahrenheit: 32,
			wantKelvin:     273,
		},
		{
			name:           "positive celsius",
			celsius:        25,
			wantCelsius:    25,
			wantFahrenheit: 77,
			wantKelvin:     298,
		},
		{
			name:           "negative celsius",
			celsius:        -10,
			wantCelsius:    -10,
			wantFahrenheit: 14,
			wantKelvin:     263,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp := NewTemperature(tt.celsius)

			if temp.Celsius != tt.wantCelsius {
				t.Errorf("NewTemperature().Celsius = %v, want %v", temp.Celsius, tt.wantCelsius)
			}

			if temp.Fahrenheit != tt.wantFahrenheit {
				t.Errorf("NewTemperature().Fahrenheit = %v, want %v", temp.Fahrenheit, tt.wantFahrenheit)
			}

			if temp.Kelvin != tt.wantKelvin {
				t.Errorf("NewTemperature().Kelvin = %v, want %v", temp.Kelvin, tt.wantKelvin)
			}
		})
	}
}
