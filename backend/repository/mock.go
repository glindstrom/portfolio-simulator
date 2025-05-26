package repository

import (
	"fmt"
	"portfolio-simulator/backend/models"
)

type MockStore struct {
	prices map[string][]models.Price
}

func NewMockStore() *MockStore {
	return &MockStore{prices: make(map[string][]models.Price)}
}
func (m *MockStore) GetPrices(symbol string) ([]models.Price, error) {
	prices, ok := m.prices[symbol]
	if !ok {
		return nil, fmt.Errorf("no prices for %s", symbol)
	}
	return prices, nil
}
func (m *MockStore) SavePrices(symbol string, prices []models.Price) error {
	m.prices[symbol] = prices
	return nil
}
func (m *MockStore) SetPrices(symbol string, prices []models.Price) {
	m.prices[symbol] = prices
}
