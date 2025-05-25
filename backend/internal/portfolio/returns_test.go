package portfolio

import (
	"fmt"
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

func TestComputePortfolioReturns_Success(t *testing.T) {
	portfolio := model.Portfolio{
		Assets: []model.Asset{
			{Ticker: "AAA", Weight: 0.5},
			{Ticker: "BBB", Weight: 0.5},
		},
	}

	mockFetch := func(ticker string) ([]float64, error) {
		switch ticker {
		case "AAA":
			return []float64{0.01, 0.02, 0.03}, nil
		case "BBB":
			return []float64{0.02, 0.01, -0.01}, nil
		default:
			return nil, fmt.Errorf("ticker not found")
		}
	}

	expected := []float64{
		0.5*0.01 + 0.5*0.02,    // 0.015
		0.5*0.02 + 0.5*0.01,    // 0.015
		0.5*0.03 + 0.5*(-0.01), // 0.01
	}

	result, err := ComputePortfolioReturns(portfolio, mockFetch)
	require.NoError(t, err)
	require.Len(t, expected, len(result))

	const epsilon = 1e-9
	require.InEpsilonSlice(t, expected, result, epsilon)
}

func TestComputePortfolioReturns_FetchError(t *testing.T) {
	portfolio := model.Portfolio{
		Assets: []model.Asset{
			{Ticker: "AAA", Weight: 1.0},
		},
	}

	mockFetch := func(ticker string) ([]float64, error) {
		return nil, fmt.Errorf("fetch failure")
	}

	_, err := ComputePortfolioReturns(portfolio, mockFetch)
	require.Error(t, err)
	require.Contains(t, err.Error(), "fetch failure")
}

func TestComputePortfolioReturns_EmptyPortfolio(t *testing.T) {
	portfolio := model.Portfolio{}

	mockFetch := func(ticker string) ([]float64, error) {
		return []float64{0.01, 0.02}, nil
	}

	result, err := ComputePortfolioReturns(portfolio, mockFetch)
	require.NoError(t, err)
	require.Nil(t, result)
}

func TestWeightedMonthlyReturns_VaryingLengths(t *testing.T) {
	p := model.Portfolio{
		Assets: []model.Asset{
			{Ticker: "AAPL", Weight: 0.6},
			{Ticker: "GOOGL", Weight: 0.4},
		},
	}

	returnsByAsset := map[string][]float64{
		"AAPL":  {0.01, 0.02, 0.03, 0.04},
		"GOOGL": {0.05, 0.06}, // Bara 2 månader
	}

	result, err := WeightedMonthlyReturns(p, returnsByAsset)
	require.NoError(t, err)
	require.Len(t, result, 2)

	expected := []float64{
		0.6*0.01 + 0.4*0.05, // 0.026
		0.6*0.02 + 0.4*0.06, // 0.036
	}

	const epsilon = 1e-9
	require.InEpsilonSlice(t, expected, result, epsilon)
}
