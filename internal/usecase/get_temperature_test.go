package usecase

import (
	"errors"
	"goexpert-clima/internal/entity"
	"goexpert-clima/internal/gateway"
	"testing"
)

// Mock implementations for testing
type MockZipCodeGateway struct {
	GetLocationFunc func(zipCode string) (*gateway.Location, error)
}

func (m *MockZipCodeGateway) GetLocation(zipCode string) (*gateway.Location, error) {
	return m.GetLocationFunc(zipCode)
}

type MockWeatherGateway struct {
	GetTemperatureFunc func(city string) (float64, error)
}

func (m *MockWeatherGateway) GetTemperature(city string) (float64, error) {
	return m.GetTemperatureFunc(city)
}

func TestGetTemperatureUseCase_Execute(t *testing.T) {
	tests := []struct {
		name           string
		zipCode        string
		mockZipCode    func(zipCode string) (*gateway.Location, error)
		mockWeather    func(city string) (float64, error)
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:    "valid zipcode",
			zipCode: "12345678",
			mockZipCode: func(zipCode string) (*gateway.Location, error) {
				return &gateway.Location{City: "São Paulo", State: "SP"}, nil
			},
			mockWeather: func(city string) (float64, error) {
				return 25.5, nil
			},
			wantErr: false,
		},
		{
			name:    "invalid zipcode format",
			zipCode: "1234567", // Only 7 digits
			mockZipCode: func(zipCode string) (*gateway.Location, error) {
				return nil, nil // Should not be called
			},
			mockWeather: func(city string) (float64, error) {
				return 0, nil // Should not be called
			},
			wantErr:        true,
			expectedErrMsg: "invalid zipcode",
		},
		{
			name:    "zipcode not found",
			zipCode: "12345678",
			mockZipCode: func(zipCode string) (*gateway.Location, error) {
				return nil, errors.New("zipcode not found")
			},
			mockWeather: func(city string) (float64, error) {
				return 0, nil // Should not be called
			},
			wantErr:        true,
			expectedErrMsg: "can not find zipcode",
		},
		{
			name:    "weather api error",
			zipCode: "12345678",
			mockZipCode: func(zipCode string) (*gateway.Location, error) {
				return &gateway.Location{City: "São Paulo", State: "SP"}, nil
			},
			mockWeather: func(city string) (float64, error) {
				return 0, errors.New("weather api error")
			},
			wantErr:        true,
			expectedErrMsg: "weather api error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockZipCodeGateway := &MockZipCodeGateway{
				GetLocationFunc: tt.mockZipCode,
			}
			mockWeatherGateway := &MockWeatherGateway{
				GetTemperatureFunc: tt.mockWeather,
			}

			// Create use case with mocks
			uc := NewGetTemperatureUseCase(mockZipCodeGateway, mockWeatherGateway)

			// Execute use case
			temp, err := uc.Execute(tt.zipCode)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check error message if expected
			if tt.wantErr && err != nil && tt.expectedErrMsg != "" {
				if err.Error() != tt.expectedErrMsg && err != ErrInvalidZipCode && err != ErrZipCodeNotFound {
					t.Errorf("Execute() error = %v, expectedErrMsg %v", err, tt.expectedErrMsg)
				}
			}

			// Check result for success case
			if !tt.wantErr {
				if temp == nil {
					t.Errorf("Execute() returned nil temperature for success case")
					return
				}

				// Check temperature values
				expectedTemp := entity.NewTemperature(25.5)
				if temp.Celsius != expectedTemp.Celsius {
					t.Errorf("Execute() temp.Celsius = %v, want %v", temp.Celsius, expectedTemp.Celsius)
				}
				if temp.Fahrenheit != expectedTemp.Fahrenheit {
					t.Errorf("Execute() temp.Fahrenheit = %v, want %v", temp.Fahrenheit, expectedTemp.Fahrenheit)
				}
				if temp.Kelvin != expectedTemp.Kelvin {
					t.Errorf("Execute() temp.Kelvin = %v, want %v", temp.Kelvin, expectedTemp.Kelvin)
				}
			}
		})
	}
}

func TestIsValidZipCode(t *testing.T) {
	tests := []struct {
		name    string
		zipCode string
		want    bool
	}{
		{
			name:    "valid zipcode",
			zipCode: "12345678",
			want:    true,
		},
		{
			name:    "invalid zipcode - too short",
			zipCode: "1234567",
			want:    false,
		},
		{
			name:    "invalid zipcode - too long",
			zipCode: "123456789",
			want:    false,
		},
		{
			name:    "invalid zipcode - contains letters",
			zipCode: "1234567a",
			want:    false,
		},
		{
			name:    "invalid zipcode - empty",
			zipCode: "",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidZipCode(tt.zipCode); got != tt.want {
				t.Errorf("isValidZipCode() = %v, want %v", got, tt.want)
			}
		})
	}
}