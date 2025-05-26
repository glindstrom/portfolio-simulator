package simulation

import (
	"math/rand"
	"time"
)

type SimulationResult struct {
	Paths [][]float64 // Each path is a time series of portfolio values
}

type Simulator struct {
	Returns      []float64 // Historical monthly returns (as decimal, e.g., 0.01 for 1%)
	InitialValue float64
	Years        int
	Simulations  int
}

func (s *Simulator) Simulate() SimulationResult {
	rand.Seed(time.Now().UnixNano())

	months := s.Years * 12
	result := SimulationResult{
		Paths: make([][]float64, s.Simulations),
	}

	for i := 0; i < s.Simulations; i++ {
		path := make([]float64, months+1)
		path[0] = s.InitialValue
		for m := 1; m <= months; m++ {
			ret := s.Returns[rand.Intn(len(s.Returns))]
			path[m] = path[m-1] * (1 + ret)
		}
		result.Paths[i] = path
	}

	return result
}
