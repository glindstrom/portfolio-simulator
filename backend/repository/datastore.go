package repository

import "portfolio-simulator/backend/models"

type DataStore interface {
	SavePrices(symbol string, prices []models.Price) error
	GetPrices(symbol string) ([]models.Price, error)
}
