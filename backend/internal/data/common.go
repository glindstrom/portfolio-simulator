package data

import (
	"sort"
	"time"
)

// PriceData holds a single historical price point.
// It's assumed that for a given series, these points represent consistent
// period-end data (e.g., month-end closing prices).
type PriceData struct {
	Date  time.Time // Date of the price record.
	Close float64   // Closing price for the recorded period.
}

// ToMonthlyReturns calculates sequential percentage returns from a slice of PriceData.
// The input 'prices' slice is expected to contain chronologically ordered data points,
// typically representing monthly values.
func ToMonthlyReturns(prices []PriceData) []float64 {
	// Ensure prices are sorted by date for correct sequential return calculation.
	// Data sources should ideally provide sorted data; this is a safeguard.
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Date.Before(prices[j].Date)
	})

	if len(prices) < 2 {
		return nil // Need at least two data points to calculate one return.
	}

	returns := make([]float64, 0, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		previousClose := prices[i-1].Close
		currentClose := prices[i].Close

		if previousClose == 0 {
			// Avoid division by zero. Treat return as 0.0 if previous close was zero
			// (e.g., due to new listing or data error).
			returns = append(returns, 0.0)
			continue
		}

		ret := (currentClose - previousClose) / previousClose
		returns = append(returns, ret)
	}
	return returns
}
