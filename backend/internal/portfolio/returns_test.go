package portfolio

import (
	"testing"

	"github.com/stretchr/testify/require"
	"portfolio-simulator/backend/internal/portfolio/model"
)

func TestWeightedMonthlyReturns(t *testing.T) {
	portfolio := model.Portfolio{
		Assets: []model.Asset{
			{Ticker: "AAA", Weight: 0.6},
			{Ticker: "BBB", Weight: 0.4},
		},
	}

	returnsByAsset := map[string][]float64{
		"AAA": {0.01, 0.02, -0.01},
		"BBB": {0.03, -0.02, 0.00},
	}

	expected := []float64{
		0.6*0.01 + 0.4*0.03,    // 0.018
		0.6*0.02 + 0.4*(-0.02), // 0.004
		0.6*(-0.01) + 0.4*0.00, // -0.006
	}

	result, err := WeightedMonthlyReturns(portfolio, returnsByAsset)
	require.NoError(t, err)
	require.Len(t, expected, len(result))

	const epsilon = 1e-9
	require.InEpsilonSlice(t, expected, result, epsilon)
}

func TestWeightedMonthlyReturns_MissingAsset(t *testing.T) {
	portfolio := model.Portfolio{
		Assets: []model.Asset{
			{Ticker: "AAA", Weight: 1.0},
			{Ticker: "CCC", Weight: 0.0}, // missing in returnsByAsset
		},
	}

	returnsByAsset := map[string][]float64{
		"AAA": {0.01, 0.02, -0.01},
	}

	_, err := WeightedMonthlyReturns(portfolio, returnsByAsset)
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing returns for asset CCC")
}

func TestWeightedMonthlyReturns_EmptyPortfolio(t *testing.T) {
	portfolio := model.Portfolio{}

	returnsByAsset := map[string][]float64{
		"AAA": {0.01, 0.02, -0.01},
	}

	result, err := WeightedMonthlyReturns(portfolio, returnsByAsset)
	require.NoError(t, err)
	require.Nil(t, result)
}
