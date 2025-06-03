package simulation

import (
	"errors"
	"math"
	"math/rand"
	"sort"
)

// Params defines the parameters required for a Monte Carlo simulation.
type Params struct {
	InitialValue     float64   // Starting value of the portfolio.
	Returns          []float64 // Historical returns (e.g., monthly) to base the simulation on.
	WithdrawalRate   float64   // Annual withdrawal rate from the portfolio (e.g., 0.04 for 4%).
	InflationPerYear float64   // Annual inflation rate (e.g., 0.02 for 2%).
	Periods          int       // Total number of periods (e.g., months) for the simulation.
	Simulations      int       // Number of Monte Carlo paths to simulate.
}

// Result holds the outcomes of a Monte Carlo simulation.
type Result struct {
	Paths       [][]float64  // Each inner slice represents a single simulated path of portfolio values.
	FinalStats  SummaryStats // Summary statistics of the final portfolio values across all paths.
	SuccessRate float64      // Proportion of paths that did not deplete before the end of the simulation period.
}

// SummaryStats provides descriptive statistics for a set of values, typically final portfolio values.
type SummaryStats struct {
	Mean   float64 // Average value.
	Median float64 // Median value.
	Min    float64 // Minimum value.
	Max    float64 // Maximum value.
}

// SimulateNormal runs Monte Carlo simulations assuming returns follow a normal distribution.
// The distribution is parameterized by the mean and sample standard deviation of the provided historical returns.
func SimulateNormal(params Params) (*Result, error) {
	if len(params.Returns) == 0 {
		return nil, errors.New("simulation: returns slice is empty, cannot calculate mean/std for normal distribution")
	}
	mean, std := meanStd(params.Returns) // Now uses sample standard deviation

	// If std is 0 (e.g., all historical returns are identical, or only one return data point),
	// the simulation becomes deterministic based on the mean return. This is generally acceptable.
	// A specific check for std == 0 for >1 data points could be an error if a stochastic result is strictly required.
	// The user's test 'TestSimulateNormal_ZeroStdDev' expects an error if std is 0 from multiple identical returns.
	if std == 0 && len(params.Returns) > 1 {
		return nil, errors.New("simulation: standard deviation of returns is zero based on provided historical data, cannot reliably perform stochastic normal simulation")
	}

	generateReturn := func() float64 {
		return rand.NormFloat64()*std + mean
	}
	return runSimulationPaths(params, generateReturn)
}

// SimulateBootstrap runs Monte Carlo simulations by randomly sampling from the provided historical returns.
func SimulateBootstrap(params Params) (*Result, error) {
	if len(params.Returns) == 0 {
		return nil, errors.New("simulation: returns slice is empty, cannot bootstrap")
	}

	generateReturn := func() float64 {
		return params.Returns[rand.Intn(len(params.Returns))]
	}
	return runSimulationPaths(params, generateReturn)
}

// runSimulationPaths executes the core Monte Carlo simulation logic for a given return generation function.
func runSimulationPaths(params Params, generateReturn func() float64) (*Result, error) {
	N := params.Simulations
	periods := params.Periods

	if periods <= 0 {
		return nil, errors.New("simulation: number of periods must be positive")
	}

	paths := make([][]float64, N)
	finalVals := make([]float64, N)
	successCount := 0

	var adjustedMonthlyWithdrawals []float64
	if params.WithdrawalRate > 0 {
		initialAnnualWithdrawal := params.InitialValue * params.WithdrawalRate
		baseMonthlyWithdrawal := initialAnnualWithdrawal / 12.0

		adjustedMonthlyWithdrawals = make([]float64, periods+1) // Index 0 unused, 1 to periods used.

		for t := 1; t <= periods; t++ {
			yearFractionForInflation := float64(t-1) / 12.0
			adjustedMonthlyWithdrawals[t] = baseMonthlyWithdrawal * math.Pow(1.0+params.InflationPerYear, yearFractionForInflation)
		}
	}

	for i := 0; i < N; i++ {
		path := make([]float64, periods+1)
		path[0] = params.InitialValue
		currentSuccess := true

		for t := 1; t <= periods; t++ {
			monthlyReturn := generateReturn()
			currentPortfolioValue := path[t-1]

			currentPortfolioValue = currentPortfolioValue * (1 + monthlyReturn)

			if params.WithdrawalRate > 0 {
				currentPortfolioValue -= adjustedMonthlyWithdrawals[t]
				if currentPortfolioValue <= 0 {
					currentPortfolioValue = 0
					currentSuccess = false
				}
			}
			path[t] = currentPortfolioValue
			if !currentSuccess {
				break
			}
		}

		paths[i] = path
		finalVals[i] = path[len(path)-1]
		if currentSuccess {
			successCount++
		}
	}

	summary := calculateSummary(finalVals)
	successRate := 0.0
	if N > 0 {
		successRate = float64(successCount) / float64(N)
	}

	return &Result{
		Paths:       paths,
		FinalStats:  summary,
		SuccessRate: successRate,
	}, nil
}

// meanStd calculates the mean and sample standard deviation of a slice of float64.
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

	if n < 2 { // Sample standard deviation requires at least 2 data points.
		return mean, 0 // Standard deviation is 0 for a single data point.
	}

	var sqDiffSum float64
	for _, v := range arr {
		diff := v - mean
		sqDiffSum += diff * diff
	}
	// Use (n-1) for sample standard deviation.
	std = math.Sqrt(sqDiffSum / (n - 1))
	return
}

// calculateSummary computes summary statistics for a slice of float64 values.
func calculateSummary(values []float64) SummaryStats {
	if len(values) == 0 {
		return SummaryStats{} // Return zeroed stats for empty input.
	}

	// Create a copy to sort, preserving the original slice if it's from elsewhere.
	sortedValues := make([]float64, len(values))
	copy(sortedValues, values)
	sort.Float64s(sortedValues)

	sum := 0.0
	for _, v := range sortedValues {
		sum += v
	}
	mean := sum / float64(len(sortedValues))

	var median float64
	n := len(sortedValues)
	mid := n / 2
	if n%2 == 0 {
		// Even number of values: average of the two middle elements.
		median = (sortedValues[mid-1] + sortedValues[mid]) / 2.0
	} else {
		// Odd number of values: the middle element.
		median = sortedValues[mid]
	}

	minValue := sortedValues[0]
	maxValue := sortedValues[len(sortedValues)-1]

	return SummaryStats{
		Mean:   mean,
		Median: median,
		Min:    minValue,
		Max:    maxValue,
	}
}
