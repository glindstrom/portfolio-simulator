package service

import (
	"fmt"
	"log"
	"math"
	"portfolio-simulator/backend/models"
	"portfolio-simulator/backend/repository"
)

const TradingDaysPerYear = 252 // Number of trading days in a year for annualization
type DataService struct {
	Store repository.DataStore
}

func NewDataService(store repository.DataStore) *DataService {
	return &DataService{Store: store}
}

// CalculateLogReturns computes daily log returns from prices.
func CalculateLogReturns(prices []models.Price) ([]float64, error) {
	if len(prices) < 2 {
		return nil, fmt.Errorf("insufficient data: %d prices", len(prices))
	}
	var returns []float64
	for i := 1; i < len(prices); i++ {
		prev := prices[i-1].Close
		curr := prices[i].Close
		if prev > 0 && curr > 0 {
			ret := math.Log(curr / prev)
			if !math.IsNaN(ret) && !math.IsInf(ret, 0) {
				returns = append(returns, ret)
			}
		}
	}
	if len(returns) == 0 {
		return nil, fmt.Errorf("no valid returns")
	}
	return returns, nil
}

// CalculateStatsFromReturns computes annualized mean return and volatility.
func CalculateStatsFromReturns(returns []float64) models.Stats {
	var sum float64
	for _, r := range returns {
		sum += r
	}
	mean := (sum / float64(len(returns))) * TradingDaysPerYear
	var variance float64
	for _, r := range returns {
		variance += math.Pow(r-sum/float64(len(returns)), 2)
	}
	volatility := math.Sqrt(variance/float64(len(returns))) * math.Sqrt(TradingDaysPerYear)
	return models.Stats{MeanReturn: mean, Volatility: volatility}
}
func (s *DataService) CalculateStats(symbol string) (models.Stats, error) {
	prices, err := s.Store.GetPrices(symbol)
	if err != nil {
		return models.Stats{}, err
	}
	returns, err := CalculateLogReturns(prices)
	if err != nil {
		log.Printf("%s: %v", symbol, err)
		return models.Stats{}, err
	}
	log.Printf("%s: Found %d valid returns", symbol, len(returns))
	stats := CalculateStatsFromReturns(returns)
	log.Printf("%s: Raw Mean=%.6f, Volatility=%.6f", symbol, stats.MeanReturn, stats.Volatility)
	return stats, nil
}
