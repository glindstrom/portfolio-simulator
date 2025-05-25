package portfolio

import (
	"fmt"
	"portfolio-simulator/backend/internal/portfolio/model"
)

// WeightedMonthlyReturns computes the weighted sum of monthly returns for the portfolio.
func WeightedMonthlyReturns(p model.Portfolio, returnsByAsset map[string][]float64) ([]float64, error) {
	if len(p.Assets) == 0 {
		return nil, nil
	}

	// Find minimum length among all assets' returns
	minMonths := -1
	for _, asset := range p.Assets {
		ret, ok := returnsByAsset[asset.Ticker]
		if !ok {
			return nil, fmt.Errorf("missing returns for asset %s", asset.Ticker)
		}
		if minMonths == -1 || len(ret) < minMonths {
			minMonths = len(ret)
		}
	}

	// Compute weighted returns up to minMonths
	weightedReturns := make([]float64, minMonths)
	for i := 0; i < minMonths; i++ {
		sum := 0.0
		for _, asset := range p.Assets {
			sum += asset.Weight * returnsByAsset[asset.Ticker][i]
		}
		weightedReturns[i] = sum
	}

	return weightedReturns, nil
}

func ComputePortfolioReturns(
	p model.Portfolio,
	fetchFunc func(string) ([]float64, error),
) ([]float64, error) {
	returnsByAsset := make(map[string][]float64)

	for _, asset := range p.Assets {
		ret, err := fetchFunc(asset.Ticker)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch returns for %s: %w", asset.Ticker, err)
		}
		returnsByAsset[asset.Ticker] = ret
	}

	return WeightedMonthlyReturns(p, returnsByAsset)
}
