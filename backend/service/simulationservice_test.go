package service

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"portfolio-simulator/backend/models"
	"portfolio-simulator/backend/repository"
	"testing"
	"time"
)

type mockDataService struct {
	stats map[string]models.Stats
}

func (m *mockDataService) CalculateStats(symbol string) (models.Stats, error) {
	stat, ok := m.stats[symbol]
	if !ok {
		return models.Stats{}, fmt.Errorf("no stats for %s", symbol)
	}
	return stat, nil
}
func TestRunSimulation(t *testing.T) {
	mockStore := repository.NewMockStore()
	dataService := NewDataService(mockStore)
	simService := NewSimulationService(dataService)
	portfolio := models.Portfolio{
		InitialValue: 100000,
		Weights:      map[string]float64{"BTC": 0.05, "EUNL.DE": 0.95},
		SellRate:     0.05,
		TaxRate:      0.30,
		Years:        2,
	}
	mockStore.SetPrices("BTC", []models.Price{
		{Date: time.Now(), Close: 100},
		{Date: time.Now(), Close: 110},
		{Date: time.Now(), Close: 121},
	})
	mockStore.SetPrices("EUNL.DE", []models.Price{
		{Date: time.Now(), Close: 100},
		{Date: time.Now(), Close: 105},
		{Date: time.Now(), Close: 110},
	})
	result, err := simService.RunSimulation(portfolio, 10)
	require.NoError(t, err, "RunSimulation should not return an error")
	require.Len(t, result.Values, 10, "Should have 10 simulations")
	for _, values := range result.Values {
		require.Len(t, values, 3, "Should have 3 yearly values (initial + 2 years)")
		require.GreaterOrEqual(t, values[len(values)-1], 0.0, "Portfolio value should not be negative")
	}
}
func TestRunSimulationNoSell(t *testing.T) {
	mockStore := repository.NewMockStore()
	dataService := NewDataService(mockStore)
	simService := NewSimulationService(dataService)
	portfolio := models.Portfolio{
		InitialValue: 100000,
		Weights:      map[string]float64{"BTC": 0.05, "EUNL.DE": 0.95},
		SellRate:     0.0,
		TaxRate:      0.30,
		Years:        2,
	}
	mockStore.SetPrices("BTC", []models.Price{
		{Date: time.Now(), Close: 100},
		{Date: time.Now(), Close: 110},
		{Date: time.Now(), Close: 121},
	})
	mockStore.SetPrices("EUNL.DE", []models.Price{
		{Date: time.Now(), Close: 100},
		{Date: time.Now(), Close: 105},
		{Date: time.Now(), Close: 110},
	})
	result, err := simService.RunSimulation(portfolio, 10)
	require.NoError(t, err, "RunSimulation should not return an error")
	require.Len(t, result.Values, 10, "Should have 10 simulations")
	for _, values := range result.Values {
		require.Greater(t, values[len(values)-1], 0.0, "Portfolio value should be positive without sell")
	}
}
func TestRunSimulationHighVolatility(t *testing.T) {
	mockStore := repository.NewMockStore()
	dataService := NewDataService(mockStore)
	simService := NewSimulationService(dataService)
	portfolio := models.Portfolio{
		InitialValue: 100000,
		Weights:      map[string]float64{"BTC": 0.05, "EUNL.DE": 0.95},
		SellRate:     0.05,
		TaxRate:      0.30,
		Years:        2,
	}
	mockStore.SetPrices("BTC", []models.Price{
		{Date: time.Now(), Close: 100},
		{Date: time.Now(), Close: 200},
		{Date: time.Now(), Close: 50},
	})
	mockStore.SetPrices("EUNL.DE", []models.Price{
		{Date: time.Now(), Close: 100},
		{Date: time.Now(), Close: 120},
		{Date: time.Now(), Close: 90},
	})
	result, err := simService.RunSimulation(portfolio, 10)
	require.NoError(t, err, "RunSimulation should not return an error")
	require.Len(t, result.Values, 10, "Should have 10 simulations")
	for _, values := range result.Values {
		require.GreaterOrEqual(t, values[len(values)-1], 0.0, "Portfolio value should not be negative with high volatility")
	}
}
