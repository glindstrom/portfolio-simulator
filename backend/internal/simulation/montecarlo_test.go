package simulation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimulate(t *testing.T) {
	sim := &Simulator{
		Returns:      []float64{0.01, -0.005, 0.02}, // Simplified historical returns
		InitialValue: 10000,
		Years:        5,
		Simulations:  100,
	}

	result := sim.Simulate()

	require.Len(t, result.Paths, sim.Simulations)

	expectedMonths := sim.Years * 12
	for _, path := range result.Paths {
		require.Len(t, path, expectedMonths+1) // +1 for initial value
		for _, val := range path {
			require.Greater(t, val, 0.0)
		}
	}
}
