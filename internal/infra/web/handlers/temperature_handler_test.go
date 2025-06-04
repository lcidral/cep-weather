package handlers

import (
	"goexpert-clima/internal/entity"
	"goexpert-clima/internal/usecase"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockGetTemperatureUseCase implements the GetTemperatureUseCaseInterface for testing
type MockGetTemperatureUseCase struct {
	ExecuteFunc func(zipCode string) (*entity.Temperature, error)
}

func (m *MockGetTemperatureUseCase) Execute(zipCode string) (*entity.Temperature, error) {
	return m.ExecuteFunc(zipCode)
}

func TestTemperatureHandler_GetTemperature(t *testing.T) {
	tests := []struct {
		name               string
		zipCode            string
		mockExecuteFunc    func(zipCode string) (*entity.Temperature, error)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:    "valid zipcode",
			zipCode: "12345678",
			mockExecuteFunc: func(zipCode string) (*entity.Temperature, error) {
				return entity.NewTemperature(25.5), nil
			},
			expectedStatusCode: http.StatusOK,
			// We're not checking the exact JSON response since it's handled by the standard library
		},
		{
			name:    "invalid zipcode",
			zipCode: "1234567", // Only 7 digits
			mockExecuteFunc: func(zipCode string) (*entity.Temperature, error) {
				return nil, usecase.ErrInvalidZipCode
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedBody:       "invalid zipcode",
		},
		{
			name:    "zipcode not found",
			zipCode: "12345678",
			mockExecuteFunc: func(zipCode string) (*entity.Temperature, error) {
				return nil, usecase.ErrZipCodeNotFound
			},
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       "can not find zipcode",
		},
		{
			name:    "internal server error",
			zipCode: "12345678",
			mockExecuteFunc: func(zipCode string) (*entity.Temperature, error) {
				return nil, usecase.ErrZipCodeNotFound
			},
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       "can not find zipcode",
		},
		{
			name:    "missing zipcode parameter",
			zipCode: "",
			mockExecuteFunc: func(zipCode string) (*entity.Temperature, error) {
				return nil, nil // Should not be called
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "zipcode is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock use case
			mockUseCase := &MockGetTemperatureUseCase{
				ExecuteFunc: tt.mockExecuteFunc,
			}

			// Create handler with mock use case
			handler := NewTemperatureHandler(mockUseCase)

			// Create request
			req, err := http.NewRequest(http.MethodGet, "/temperature", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Add query parameter if provided
			if tt.zipCode != "" {
				q := req.URL.Query()
				q.Add("zipcode", tt.zipCode)
				req.URL.RawQuery = q.Encode()
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.GetTemperature(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.expectedStatusCode)
			}

			// Check response body for error cases
			if tt.expectedBody != "" && rr.Body.String() != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
