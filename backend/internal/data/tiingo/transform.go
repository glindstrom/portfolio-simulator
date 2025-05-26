package tiingo

import (
	"sort"
)

func ToMonthlyReturns(data []PriceData) []float64 {
	if len(data) == 0 {
		return nil
	}

	// Sort by date ascending
	sort.Slice(data, func(i, j int) bool {
		return data[i].Date.Before(data[j].Date)
	})

	monthlyPrices := make(map[string]float64) // key: "YYYY-MM"

	for _, d := range data {
		key := d.Date.Format("2006-01")
		// always overwrite with the latest price in month
		monthlyPrices[key] = d.Close
	}

	// Extract sorted keys
	keys := make([]string, 0, len(monthlyPrices))
	for k := range monthlyPrices {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Compute monthly returns
	var returns []float64
	var prev float64

	for _, k := range keys {
		price := monthlyPrices[k]
		if prev != 0 {
			ret := (price - prev) / prev
			returns = append(returns, ret)
		}
		prev = price
	}

	return returns
}
