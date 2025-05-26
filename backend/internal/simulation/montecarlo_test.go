package simulation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimulateNormal(t *testing.T) {
	params := Params{
		InitialValue: 1000,
		Returns:      []float64{0.01, -0.02, 0.03, 0.015, -0.005},
		Simulations:  1000,
		Periods:      10,
	}
	result, err := SimulateNormal(params)
	require.NoError(t, err)
	require.Len(t, result.Paths, params.Simulations)

	for _, path := range result.Paths {
		require.Len(t, path, params.Periods+1)
		require.Equal(t, params.InitialValue, path[0])
		for i := 1; i <= params.Periods; i++ {
			require.Greater(t, path[i], 0.0)
		}
	}

	require.LessOrEqual(t, result.FinalStats.Min, result.FinalStats.Mean)
	require.LessOrEqual(t, result.FinalStats.Min, result.FinalStats.Median)
	require.GreaterOrEqual(t, result.FinalStats.Max, result.FinalStats.Mean)
	require.GreaterOrEqual(t, result.FinalStats.Max, result.FinalStats.Median)
}

func TestSimulateBootstrap(t *testing.T) {
	params := Params{
		InitialValue: 1000,
		Returns:      []float64{0.01, -0.02, 0.03, 0.015, -0.005},
		Simulations:  1000,
		Periods:      10,
	}
	result, err := SimulateBootstrap(params)
	require.NoError(t, err)
	require.Len(t, result.Paths, params.Simulations)

	for _, path := range result.Paths {
		require.Len(t, path, params.Periods+1)
		require.Equal(t, params.InitialValue, path[0])
		for i := 1; i <= params.Periods; i++ {
			require.Greater(t, path[i], 0.0)
		}
	}

	require.LessOrEqual(t, result.FinalStats.Min, result.FinalStats.Mean)
	require.LessOrEqual(t, result.FinalStats.Min, result.FinalStats.Median)
	require.GreaterOrEqual(t, result.FinalStats.Max, result.FinalStats.Mean)
	require.GreaterOrEqual(t, result.FinalStats.Max, result.FinalStats.Median)
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
}

func TestRunSimulationPaths_SuccessRateCorrect(t *testing.T) {
	params := Params{
		InitialValue:     1000,
		WithdrawalRate:   0.1,
		InflationPerYear: 0.0,
		Simulations:      4,
		Periods:          2,
	}
	generateZeroReturn := func() float64 {
		return 0.0
	}
	result, err := runSimulationPaths(params, generateZeroReturn)
	require.NoError(t, err)
	require.Equal(t, 1.0, result.SuccessRate)
}

func TestRunSimulationPaths_PartialFailures(t *testing.T) {
	params := Params{
		InitialValue:     100,
		WithdrawalRate:   0.5,
		InflationPerYear: 0.0,
		Simulations:      4,
		Periods:          2,
	}
	returns := []float64{
		0.1, 0.1, // Sim 0
		-0.9, -0.9, // Sim 1
		0.1, 0.1, // Sim 2
		-0.9, -0.9, // Sim 3
	}
	i := 0
	generateReturn := func() float64 {
		val := returns[i%len(returns)]
		i++
		return val
	}

	result, err := runSimulationPaths(params, generateReturn)
	require.NoError(t, err)
	require.Len(t, result.Paths, params.Simulations)
	require.Equal(t, 0.5, result.SuccessRate)
}
