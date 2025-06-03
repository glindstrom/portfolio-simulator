package simulation

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSimulateNormal_NoWithdrawals tests SimulateNormal without withdrawals.
func TestSimulateNormal_NoWithdrawals(t *testing.T) {
	params := Params{
		InitialValue: 1000,
		Returns:      []float64{0.01, -0.02, 0.03, 0.015, -0.005},
		Simulations:  100,
		Periods:      10,
	}
	result, err := SimulateNormal(params)
	require.NoError(t, err)
	require.NotNil(t, result, "Result should not be nil on success")
	require.Len(t, result.Paths, params.Simulations)

	for _, path := range result.Paths {
		require.Len(t, path, params.Periods+1)
		require.Equal(t, params.InitialValue, path[0])
		for i := 1; i <= params.Periods; i++ {
			require.GreaterOrEqual(t, path[i], 0.0, "Path value should not be negative")
		}
	}

	require.LessOrEqual(t, result.FinalStats.Min, result.FinalStats.Mean)
	require.LessOrEqual(t, result.FinalStats.Min, result.FinalStats.Median)
	require.GreaterOrEqual(t, result.FinalStats.Max, result.FinalStats.Mean)
	require.GreaterOrEqual(t, result.FinalStats.Max, result.FinalStats.Median)
	if len(result.Paths) > 0 {
		require.Equal(t, 1.0, result.SuccessRate, "Success rate should be 100% with no withdrawals")
	}
}

// TestSimulateBootstrap_NoWithdrawals tests SimulateBootstrap without withdrawals.
func TestSimulateBootstrap_NoWithdrawals(t *testing.T) {
	params := Params{
		InitialValue: 1000,
		Returns:      []float64{0.01, -0.02, 0.03, 0.015, -0.005},
		Simulations:  100,
		Periods:      10,
	}
	result, err := SimulateBootstrap(params)
	require.NoError(t, err)
	require.NotNil(t, result, "Result should not be nil on success")
	require.Len(t, result.Paths, params.Simulations)

	for _, path := range result.Paths {
		require.Len(t, path, params.Periods+1)
		require.Equal(t, params.InitialValue, path[0])
		for i := 1; i <= params.Periods; i++ {
			require.GreaterOrEqual(t, path[i], 0.0, "Path value should not be negative")
		}
	}

	require.LessOrEqual(t, result.FinalStats.Min, result.FinalStats.Mean)
	require.LessOrEqual(t, result.FinalStats.Min, result.FinalStats.Median)
	require.GreaterOrEqual(t, result.FinalStats.Max, result.FinalStats.Mean)
	require.GreaterOrEqual(t, result.FinalStats.Max, result.FinalStats.Median)
	if len(result.Paths) > 0 {
		require.Equal(t, 1.0, result.SuccessRate, "Success rate should be 100% with no withdrawals")
	}
}

// TestSimulateNormal_WithWithdrawals tests SimulateNormal with withdrawals and inflation.
func TestSimulateNormal_WithWithdrawals(t *testing.T) {
	params := Params{
		InitialValue:     10000,
		Returns:          []float64{0.01, 0.005, 0.02, -0.01, 0.015, 0.008, -0.003},
		WithdrawalRate:   0.04,
		InflationPerYear: 0.02,
		Simulations:      100,
		Periods:          12 * 5,
	}
	result, err := SimulateNormal(params)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Paths, params.Simulations)
	require.True(t, result.SuccessRate >= 0 && result.SuccessRate <= 1, "Success rate must be between 0 and 1")
}

// TestSimulateBootstrap_WithWithdrawals tests SimulateBootstrap with withdrawals and inflation.
func TestSimulateBootstrap_WithWithdrawals(t *testing.T) {
	params := Params{
		InitialValue:     10000,
		Returns:          []float64{0.01, 0.005, 0.02, -0.01, 0.015, 0.008, -0.003},
		WithdrawalRate:   0.04,
		InflationPerYear: 0.02,
		Simulations:      100,
		Periods:          12 * 5,
	}
	result, err := SimulateBootstrap(params)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Paths, params.Simulations)
	require.True(t, result.SuccessRate >= 0 && result.SuccessRate <= 1, "Success rate must be between 0 and 1")
}

func TestSimulateNormal_EmptyReturns(t *testing.T) {
	params := Params{
		InitialValue: 1000,
		Returns:      []float64{},
		Simulations:  100,
		Periods:      10,
	}
	_, err := SimulateNormal(params)
	require.Error(t, err)
	require.Contains(t, err.Error(), "returns slice is empty")
}

func TestSimulateBootstrap_EmptyReturns(t *testing.T) {
	params := Params{
		InitialValue: 1000,
		Returns:      []float64{},
		Simulations:  100,
		Periods:      10,
	}
	_, err := SimulateBootstrap(params)
	require.Error(t, err)
	require.Contains(t, err.Error(), "returns slice is empty")
}

func TestSimulateNormal_ZeroStdDev(t *testing.T) {
	params := Params{
		InitialValue: 1000,
		Returns:      []float64{0.01, 0.01, 0.01, 0.01},
		Simulations:  100,
		Periods:      10,
	}
	_, err := SimulateNormal(params)
	require.Error(t, err)
	require.Contains(t, err.Error(), "standard deviation of returns is zero")
}

func TestRunSimulationPaths_SuccessRateCorrect(t *testing.T) {
	params := Params{
		InitialValue:     1000,
		WithdrawalRate:   0.01,
		InflationPerYear: 0.0,
		Simulations:      4,
		Periods:          2,
		Returns:          []float64{0.0},
	}
	generateZeroReturn := func() float64 { return 0.0 }
	result, err := runSimulationPaths(params, generateZeroReturn)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 1.0, result.SuccessRate)
}

func TestRunSimulationPaths_PartialFailures(t *testing.T) {
	params := Params{
		InitialValue:     100,
		WithdrawalRate:   0.5,
		InflationPerYear: 0.0,
		Simulations:      4,
		Periods:          2,
		Returns:          []float64{0.0},
	}
	deterministicReturns := []float64{
		0.1, 0.1,
		-0.9, -0.9,
		0.1, 0.1,
		-0.9, -0.9,
	}
	idx := 0
	generateDeterministicReturn := func() float64 {
		val := deterministicReturns[idx]
		idx++
		return val
	}

	result, err := runSimulationPaths(params, generateDeterministicReturn)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Paths, params.Simulations)
	require.Equal(t, 0.5, result.SuccessRate, "Expected 2/4 paths to succeed")
}

func TestRunSimulationPaths_AllPathsFail(t *testing.T) {
	params := Params{
		InitialValue:     100,
		WithdrawalRate:   7.2,
		InflationPerYear: 0.0,
		Simulations:      4,
		Periods:          2,
		Returns:          []float64{0.0},
	}
	generateZeroReturn := func() float64 { return 0.0 }
	result, err := runSimulationPaths(params, generateZeroReturn)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 0.0, result.SuccessRate, "All paths should fail")
}

func TestMeanStd(t *testing.T) {
	testCases := []struct {
		name      string
		data      []float64
		expMean   float64
		expStdDev float64 // Expected Sample Standard Deviation
	}{
		{"empty slice", []float64{}, 0.0, 0.0},
		{"single value", []float64{5.0}, 5.0, 0.0},                                   // Sample std dev for n=1 is 0
		{"multiple values", []float64{1.0, 2.0, 3.0, 4.0, 5.0}, 3.0, math.Sqrt(2.5)}, // Corrected
		{"negative values", []float64{-1.0, -2.0, -3.0}, -2.0, math.Sqrt(1.0)},       // Corrected
		{"mixed values", []float64{-1.0, 0.0, 1.0}, 0.0, math.Sqrt(1.0)},             // Corrected
		{"all same values", []float64{2.0, 2.0, 2.0}, 2.0, 0.0},                      // Std dev is 0
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mean, stdDev := meanStd(tc.data)
			require.InDelta(t, tc.expMean, mean, 0.0001)
			require.InDelta(t, tc.expStdDev, stdDev, 0.0001)
		})
	}
}

func TestCalculateSummary(t *testing.T) {
	testCases := []struct {
		name      string
		values    []float64
		expMin    float64
		expMax    float64
		expMean   float64
		expMedian float64 // Expected Median (average of two middle for even)
	}{
		{"empty slice", []float64{}, 0, 0, 0, 0},
		{"single value", []float64{42.0}, 42.0, 42.0, 42.0, 42.0},
		{"sorted odd", []float64{10.0, 20.0, 30.0, 40.0, 50.0}, 10.0, 50.0, 30.0, 30.0},
		{"unsorted odd", []float64{30.0, 10.0, 50.0, 20.0, 40.0}, 10.0, 50.0, 30.0, 30.0},
		{"sorted even", []float64{10.0, 20.0, 30.0, 40.0}, 10.0, 40.0, 25.0, 25.0},   // Corrected median
		{"unsorted even", []float64{40.0, 10.0, 30.0, 20.0}, 10.0, 40.0, 25.0, 25.0}, // Corrected median
		{"negative values", []float64{-30.0, -10.0, -50.0, -20.0, -40.0}, -50.0, -10.0, -30.0, -30.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			summary := calculateSummary(tc.values)
			require.Equal(t, tc.expMin, summary.Min)
			require.Equal(t, tc.expMax, summary.Max)
			require.InDelta(t, tc.expMean, summary.Mean, 0.0001)
			require.Equal(t, tc.expMedian, summary.Median)
		})
	}
}
