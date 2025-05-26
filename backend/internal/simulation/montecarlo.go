package simulation

import (
	"errors"
	"math"
	"math/rand"
	"sort"
)

type Params struct {
	InitialValue     float64
	Returns          []float64
	WithdrawalRate   float64
	InflationPerYear float64
	Periods          int
	Simulations      int
}

type Result struct {
	Paths       [][]float64
	FinalStats  SummaryStats
	SuccessRate float64
}

type SummaryStats struct {
	Mean   float64
	Median float64
	Min    float64
	Max    float64
}

func SimulateNormal(params Params) (*Result, error) {
	if len(params.Returns) == 0 {
		return nil, errors.New("returns slice is empty")
	}
	mean, std := meanStd(params.Returns)
	if std == 0 {
		return nil, errors.New("standard deviation of returns is zero")
	}

	generateReturn := func() float64 {
		return rand.NormFloat64()*std + mean
	}
	return runSimulationPaths(params, generateReturn)
}

func SimulateBootstrap(params Params) (*Result, error) {
	if len(params.Returns) == 0 {
		return nil, errors.New("returns slice is empty")
	}

	generateReturn := func() float64 {
		return params.Returns[rand.Intn(len(params.Returns))]
	}
	return runSimulationPaths(params, generateReturn)
}

func runSimulationPaths(params Params, generateReturn func() float64) (*Result, error) {
	N := params.Simulations
	periods := params.Periods

	paths := make([][]float64, N)
	finalVals := make([]float64, N)
	successCount := 0

	withdrawal := params.InitialValue * params.WithdrawalRate
	inflationFactor := 1.0 + params.InflationPerYear
	adjustedWithdrawals := make([]float64, periods+1)
	for t := range adjustedWithdrawals {
		adjustedWithdrawals[t] = withdrawal * math.Pow(inflationFactor, float64(t)/12)
	}

	for i := 0; i < N; i++ {
		path := make([]float64, periods+1)
		path[0] = params.InitialValue
		success := true

		for t := 1; t <= periods; t++ {
			r := generateReturn()
			path[t] = path[t-1] * (1 + r)

			if params.WithdrawalRate > 0 {
				path[t] -= adjustedWithdrawals[t]
				if path[t] <= 0 {
					path[t] = 0
					success = false
					break
				}
			}
		}

		paths[i] = path
		finalVals[i] = path[len(path)-1]
		if success {
			successCount++
		}
	}

	summary := calculateSummary(finalVals)
	successRate := float64(successCount) / float64(N)

	return &Result{
		Paths:       paths,
		FinalStats:  summary,
		SuccessRate: successRate,
	}, nil
}

// meanStd calculates the mean and standard deviation of a slice of float64.
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

// calculateSummary computes summary statistics for a slice of float64 values.
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
	minValue := values[0]
	maxValue := values[len(values)-1]

	return SummaryStats{
		Mean:   mean,
		Median: median,
		Min:    minValue,
		Max:    maxValue,
	}
}
