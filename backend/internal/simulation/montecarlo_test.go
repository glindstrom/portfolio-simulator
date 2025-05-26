package simulation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimulate(t *testing.T) {
	// Arrange
	sim := Simulator{
		InitialValue: 1000,
		Returns:      []float64{0.01, -0.02, 0.03, 0.015, -0.005},
	}

	N := 1000
	periods := 10

	// Act
	result, err := sim.Simulate(N, periods)

	// Assert
	require.NoError(t, err)
	require.Len(t, result.Paths, N)

	for _, path := range result.Paths {
		require.Len(t, path, periods+1)
		require.Equal(t, 1000.0, path[0]) // initial value should be unchanged
		for i := 1; i <= periods; i++ {
			require.Greater(t, path[i], 0.0) // portfolio value should be positive
		}
	}

	// Check summary stats sanity
	require.LessOrEqual(t, result.FinalStats.Min, result.FinalStats.Mean)
	require.LessOrEqual(t, result.FinalStats.Min, result.FinalStats.Median)
	require.GreaterOrEqual(t, result.FinalStats.Max, result.FinalStats.Mean)
	require.GreaterOrEqual(t, result.FinalStats.Max, result.FinalStats.Median)
}

func TestBootstrapSim(t *testing.T) {
	// Arrange
	sim := Simulator{
		InitialValue: 1000,
		Returns:      []float64{0.01, -0.02, 0.03, 0.015, -0.005},
	}

	N := 1000
	periods := 10

	// Act
	result, err := sim.BootstrapSim(N, periods)

	// Assert
	require.NoError(t, err)
	require.Len(t, result.Paths, N)

	for _, path := range result.Paths {
		require.Len(t, path, periods+1)
		require.Equal(t, 1000.0, path[0])
		for i := 1; i <= periods; i++ {
			require.Greater(t, path[i], 0.0)
		}
	}

	require.LessOrEqual(t, result.FinalStats.Min, result.FinalStats.Mean)
	require.LessOrEqual(t, result.FinalStats.Min, result.FinalStats.Median)
	require.GreaterOrEqual(t, result.FinalStats.Max, result.FinalStats.Mean)
	require.GreaterOrEqual(t, result.FinalStats.Max, result.FinalStats.Median)
}

func TestSimulate_EmptyReturns(t *testing.T) {
	sim := Simulator{
		InitialValue: 1000,
		Returns:      []float64{},
	}
	_, err := sim.Simulate(100, 10)
	require.Error(t, err)
}

func TestBootstrapSim_EmptyReturns(t *testing.T) {
	sim := Simulator{
		InitialValue: 1000,
		Returns:      []float64{},
	}
	_, err := sim.BootstrapSim(100, 10)
	require.Error(t, err)
}
