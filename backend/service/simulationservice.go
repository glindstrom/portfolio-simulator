package service

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"portfolio-simulator/backend/models"
)

type SimulationService struct {
	DataService *DataService
}

func NewSimulationService(dataService *DataService) *SimulationService {
	return &SimulationService{DataService: dataService}
}
func (s *SimulationService) RunSimulation(portfolio models.Portfolio, nSimulations int) (models.SimulationResult, error) {
	symbols := []string{"BTC", "EUNL.DE"}
	stats := make(map[string]models.Stats)
	for _, symbol := range symbols {
		stat, err := s.DataService.CalculateStats(symbol)
		if err != nil {
			return models.SimulationResult{}, fmt.Errorf("stats error for %s: %v", symbol, err)
		}
		log.Printf("%s: MeanReturn=%.4f, Volatility=%.4f", symbol, stat.MeanReturn, stat.Volatility)
		if math.IsNaN(stat.MeanReturn) || math.IsNaN(stat.Volatility) {
			return models.SimulationResult{}, fmt.Errorf("invalid stats for %s", symbol)
		}
		stats[symbol] = stat
	}
	result := models.SimulationResult{Values: make([][]float64, nSimulations)}
	for i := 0; i < nSimulations; i++ {
		value := portfolio.InitialValue
		yearlyValues := []float64{value}
		for year := 0; year < portfolio.Years; year++ {
			portfolioReturn := 0.0
			for symbol, weight := range portfolio.Weights {
				stat := stats[symbol]
				annualReturn := rand.NormFloat64()*stat.Volatility + stat.MeanReturn
				if math.IsNaN(annualReturn) || math.IsInf(annualReturn, 0) {
					annualReturn = 0.0 // Fallback
				}
				portfolioReturn += weight * annualReturn
			}
			if math.IsNaN(portfolioReturn) || math.IsInf(portfolioReturn, 0) {
				portfolioReturn = 0.0
			}
			value *= math.Exp(portfolioReturn)
			if math.IsNaN(value) || math.IsInf(value, 0) {
				value = yearlyValues[len(yearlyValues)-1] // Keep last valid value
			}
			if portfolio.SellRate > 0 {
				sellAmount := value * portfolio.SellRate
				gain := sellAmount * (value/portfolio.InitialValue - 1)
				tax := 0.0
				if gain > 0 {
					if gain <= 30000 {
						tax = gain * 0.3
					} else {
						tax = 30000*0.3 + (gain-30000)*0.34
					}
				}
				value -= sellAmount + tax
				if math.IsNaN(value) || value < 0 {
					value = 0.0 // Prevent negative or NaN
				}
			}
			yearlyValues = append(yearlyValues, value)
		}
		result.Values[i] = yearlyValues
	}
	return result, nil
}
