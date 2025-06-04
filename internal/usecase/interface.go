package usecase

import "goexpert-clima/internal/entity"

// GetTemperatureUseCaseInterface defines the interface for the GetTemperatureUseCase
type GetTemperatureUseCaseInterface interface {
	Execute(zipCode string) (*entity.Temperature, error)
}