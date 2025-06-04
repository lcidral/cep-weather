package handlers

import (
	"encoding/json"
	"goexpert-clima/internal/usecase"
	"log"
	"net/http"
)

type TemperatureHandler struct {
	GetTemperatureUseCase usecase.GetTemperatureUseCaseInterface
}

func NewTemperatureHandler(getTemperatureUseCase usecase.GetTemperatureUseCaseInterface) *TemperatureHandler {
	return &TemperatureHandler{
		GetTemperatureUseCase: getTemperatureUseCase,
	}
}

func (h *TemperatureHandler) GetTemperature(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request for temperature: %s %s", r.Method, r.URL.String())

	// Only accept GET requests
	if r.Method != http.MethodGet {
		log.Printf("Method not allowed: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Get zip code from query parameter
	zipCode := r.URL.Query().Get("zipcode")
	log.Printf("Processing request for zipcode: %s", zipCode)
	if zipCode == "" {
		log.Print("Error: zipcode parameter is missing")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("zipcode is required"))
		return
	}

	// Execute use case
	log.Printf("Calling GetTemperatureUseCase.Execute with zipcode: %s", zipCode)
	temperature, err := h.GetTemperatureUseCase.Execute(zipCode)
	if err != nil {
		log.Printf("Error in GetTemperatureUseCase.Execute: %v", err)
		switch err {
		case usecase.ErrInvalidZipCode:
			log.Printf("Invalid zipcode format: %s", zipCode)
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("invalid zipcode"))
		case usecase.ErrZipCodeNotFound:
			log.Printf("Zipcode not found: %s", zipCode)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("can not find zipcode"))
		default:
			log.Printf("Internal server error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}
		return
	}

	// Return success response
	log.Printf("Successfully retrieved temperature for zipcode %s: %v°C", zipCode, temperature.Celsius)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(temperature)
}
