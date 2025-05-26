package tiingo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestToMonthlyReturns(t *testing.T) {
	raw := []PriceData{
		{Date: parseDate("2024-01-31"), Close: 100},
		{Date: parseDate("2024-02-29"), Close: 110},
		{Date: parseDate("2024-03-31"), Close: 121},
		{Date: parseDate("2024-04-30"), Close: 115.5},
	}

	expected := []float64{
		0.10,            // (110 - 100) / 100
		0.10,            // (121 - 110) / 110
		-0.045454545455, // (115.5 - 121) / 121
	}

	returns := ToMonthlyReturns(raw)

	require.Len(t, returns, len(expected), "unexpected number of monthly returns")

	for i := range expected {
		require.InEpsilonf(t, expected[i], returns[i], 1e-6,
			"return[%d]: expected %.6f, got %.6f", i, expected[i], returns[i])
	}
}

// Helper to parse date strings
func parseDate(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	require.NoError(nil, err)
	return t
}
