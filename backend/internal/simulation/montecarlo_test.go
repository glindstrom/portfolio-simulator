package simulation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimulateWithWithdrawal(t *testing.T) {
	sim := &Simulator{
		Returns:          []float64{0.01, -0.005, 0.02},
		InitialValue:     10000,
		Years:            5,
		Simulations:      100,
		AnnualWithdrawal: 0.04, // 4% withdrawal per year of initial value
		InflationRate:    0.02, // 2% annual inflation rate
	}

	result := sim.Simulate()

	require.Len(t, result.Paths, sim.Simulations)

	months := sim.Years * 12
	for _, path := range result.Paths {
		require.Len(t, path, months+1)

		for _, val := range path {
			// Portfolio value should never be negative
			require.GreaterOrEqual(t, val, 0.0)
		}
	}
}
