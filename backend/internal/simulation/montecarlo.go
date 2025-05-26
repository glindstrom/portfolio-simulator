package simulation

import (
	"math"
	"math/rand"
)

type SimulationResult struct {
	Paths [][]float64 // Each path is a time series of portfolio values
}

type Simulator struct {
	Returns          []float64
	InitialValue     float64
	Years            int
	Simulations      int
	AnnualWithdrawal float64 // % of initial portfolio per year (e.g. 0.04 for 4%)
	InflationRate    float64 // e.g. 0.02 for 2%
}

func (s *Simulator) Simulate() SimulationResult {
	months := s.Years * 12
	result := SimulationResult{
		Paths: make([][]float64, s.Simulations),
	}

	for i := 0; i < s.Simulations; i++ {
		path := make([]float64, months+1)
		value := s.InitialValue
		path[0] = value

		for m := 1; m <= months; m++ {
			// Withdraw at the beginning of each year (month 1, 13, 25, ...)
			if m%12 == 1 {
				year := m / 12
				// Withdrawal grows with inflation relative to initial portfolio value
				withdrawal := s.AnnualWithdrawal * s.InitialValue * math.Pow(1+s.InflationRate, float64(year))
				value -= withdrawal
				if value < 0 {
					value = 0
				}
			}

			// Apply random monthly return from historical returns
			ret := s.Returns[rand.Intn(len(s.Returns))]
			value *= (1 + ret)
			path[m] = value
		}

		result.Paths[i] = path
	}

	return result
}
