package portfolio

import (
	"fmt"
	"portfolio-simulator/backend/internal/portfolio/model"
)

// WeightedMonthlyReturns computes the weighted sum of monthly returns for the portfolio.
// Assumes each asset has the same length of returns slice.
func WeightedMonthlyReturns(p model.Portfolio, returnsByAsset map[string][]float64) ([]float64, error) {
	if len(p.Assets) == 0 {
		return nil, nil
	}

	// Number of months, assuming all have equal length
	numMonths := -1
	for _, asset := range p.Assets {
		ret, ok := returnsByAsset[asset.Ticker]
		if !ok {
			return nil, fmt.Errorf("missing returns for asset %s", asset.Ticker)
		}
		if numMonths == -1 {
			numMonths = len(ret)
		} else if numMonths != len(ret) {
			return nil, fmt.Errorf("returns length mismatch for asset %s", asset.Ticker)
		}
	}

	weightedReturns := make([]float64, numMonths)
	for i := 0; i < numMonths; i++ {
		sum := 0.0
		for _, asset := range p.Assets {
			sum += asset.Weight * returnsByAsset[asset.Ticker][i]
		}
		weightedReturns[i] = sum
	}

	return weightedReturns, nil
}
