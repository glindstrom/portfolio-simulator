package service

import (
	"github.com/stretchr/testify/require"
	"math"
	"portfolio-simulator/backend/models"
	"testing"
	"time"
)

func TestCalculateLogReturns(t *testing.T) {
	tests := []struct {
		name    string
		prices  []models.Price
		want    []float64
		wantErr bool
	}{
		{
			name: "Valid prices",
			prices: []models.Price{
				{Date: time.Now(), Close: 100},
				{Date: time.Now(), Close: 110},
				{Date: time.Now(), Close: 121},
			},
			want:    []float64{math.Log(110.0 / 100.0), math.Log(121.0 / 110.0)},
			wantErr: false,
		},
		{
			name:    "Too few prices",
			prices:  []models.Price{{Date: time.Now(), Close: 100}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Zero price",
			prices: []models.Price{
				{Date: time.Now(), Close: 100},
				{Date: time.Now(), Close: 0},
				{Date: time.Now(), Close: 110},
			},
			want:    nil,
			wantErr: true, // Changed to true since no valid returns are generated
		},
		{
			name: "Zero price with valid return",
			prices: []models.Price{
				{Date: time.Now(), Close: 100},
				{Date: time.Now(), Close: 110},
				{Date: time.Now(), Close: 0},
				{Date: time.Now(), Close: 120},
			},
			want:    []float64{math.Log(110.0 / 100.0)}, // Only one valid return
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateLogReturns(tt.prices)
			if tt.wantErr {
				require.Error(t, err, "CalculateLogReturns should return an error")
				require.Nil(t, got, "CalculateLogReturns should return nil on error")
			} else {
				require.NoError(t, err, "CalculateLogReturns should not return an error")
				require.Equal(t, len(tt.want), len(got), "CalculateLogReturns length mismatch")
				for i, v := range got {
					require.InDelta(t, tt.want[i], v, 1e-6, "CalculateLogReturns[%d] value mismatch", i)
				}
			}
		})
	}
}

func TestCalculateStatsFromReturns(t *testing.T) {
	returns := []float64{0.01, 0.02, -0.01, 0.015}
	stats := CalculateStatsFromReturns(returns)

	// Calculate expected values
	dailyMean := (0.01 + 0.02 - 0.01 + 0.015) / 4
	expectedMean := dailyMean * TradingDaysPerYear

	var variance float64
	for _, r := range returns {
		variance += math.Pow(r-dailyMean, 2)
	}
	expectedVolatility := math.Sqrt(variance/4) * math.Sqrt(TradingDaysPerYear)

	require.InDelta(t, expectedMean, stats.MeanReturn, 1e-6, "MeanReturn mismatch")
	require.InDelta(t, expectedVolatility, stats.Volatility, 1e-6, "Volatility mismatch")
}

func TestAnnualizationWithTradingDays(t *testing.T) {
	returns := []float64{0.001, 0.002, -0.001, 0.0015}
	stats := CalculateStatsFromReturns(returns)

	dailyMean := (0.001 + 0.002 - 0.001 + 0.0015) / 4
	expectedMean := dailyMean * TradingDaysPerYear
	require.InDelta(t, expectedMean, stats.MeanReturn, 1e-6, "Annualized MeanReturn mismatch")

	var variance float64
	for _, r := range returns {
		variance += math.Pow(r-dailyMean, 2)
	}
	expectedVolatility := math.Sqrt(variance/4) * math.Sqrt(TradingDaysPerYear)
	require.InDelta(t, expectedVolatility, stats.Volatility, 1e-6, "Annualized Volatility mismatch")
}
