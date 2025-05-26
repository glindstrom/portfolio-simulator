package simulation

import (
	"errors"
	"math"
	"math/rand"
	"sort"
)

type SimulationResult struct {
	Paths      [][]float64 // Each path is a slice of portfolio values over time
	FinalStats SummaryStats
}

type SummaryStats struct {
	Mean   float64
	Median float64
	Min    float64
	Max    float64
}

type Simulator struct {
	InitialValue float64
	Returns      []float64 // Historical returns or simulated returns
}

// Simulate runs a Monte Carlo simulation using normally distributed returns.
// N = number of simulations, periods = time steps (e.g. months or years)
func (s *Simulator) Simulate(N int, periods int) (*SimulationResult, error) {
	if len(s.Returns) == 0 {
		return nil, errors.New("returns slice is empty")
	}
	mean, std := meanStd(s.Returns)
	if std == 0 {
		return nil, errors.New("standard deviation of returns is zero")
	}

	paths := make([][]float64, N)
	for i := 0; i < N; i++ {
		path := make([]float64, periods+1)
		path[0] = s.InitialValue
		for t := 1; t <= periods; t++ {
			// simulate return with normal distribution
			r := rand.NormFloat64()*std + mean
			path[t] = path[t-1] * (1 + r)
		}
		paths[i] = path
	}

	finalVals := extractFinalValues(paths)
	summary := calculateSummary(finalVals)

	return &SimulationResult{
		Paths:      paths,
		FinalStats: summary,
	}, nil
}

// BootstrapSim runs a Monte Carlo simulation using bootstrapping (sampling with replacement) from historical returns.
func (s *Simulator) BootstrapSim(N int, periods int) (*SimulationResult, error) {
	if len(s.Returns) == 0 {
		return nil, errors.New("returns slice is empty")
	}

	nReturns := len(s.Returns)
	paths := make([][]float64, N)
	for i := 0; i < N; i++ {
		path := make([]float64, periods+1)
		path[0] = s.InitialValue
		for t := 1; t <= periods; t++ {
			r := s.Returns[rand.Intn(nReturns)] // bootstrap sample
			path[t] = path[t-1] * (1 + r)
		}
		paths[i] = path
	}

	finalVals := extractFinalValues(paths)
	summary := calculateSummary(finalVals)

	return &SimulationResult{
		Paths:      paths,
		FinalStats: summary,
	}, nil
}

// meanStd calculates mean and standard deviation of a slice of floats.
func meanStd(arr []float64) (mean, std float64) {
	n := float64(len(arr))
	if n == 0 {
		return 0, 0
	}
	sum := 0.0
	for _, v := range arr {
		sum += v
	}
	mean = sum / n

	var sqDiffSum float64
	for _, v := range arr {
		diff := v - mean
		sqDiffSum += diff * diff
	}
	std = math.Sqrt(sqDiffSum / n)
	return
}

// extractFinalValues extracts the last value from each simulation path.
func extractFinalValues(paths [][]float64) []float64 {
	finalVals := make([]float64, len(paths))
	for i, path := range paths {
		finalVals[i] = path[len(path)-1]
	}
	return finalVals
}

// calculateSummary computes mean, median, min and max from a slice of float64.
func calculateSummary(values []float64) SummaryStats {
	if len(values) == 0 {
		return SummaryStats{}
	}
	sort.Float64s(values)

	mean := 0.0
	for _, v := range values {
		mean += v
	}
	mean /= float64(len(values))

	median := values[len(values)/2]
	min := values[0]
	max := values[len(values)-1]

	return SummaryStats{
		Mean:   mean,
		Median: median,
		Min:    min,
		Max:    max,
	}
}
