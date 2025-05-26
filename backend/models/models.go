package models

import "time"

type Price struct {
	Date  time.Time
	Close float64
}
type Stats struct {
	MeanReturn float64
	Volatility float64
}
type Portfolio struct {
	InitialValue float64
	Weights      map[string]float64 // e.g., {"BTC-USD": 0.05, "EUNL.DE": 0.95}
	SellRate     float64            // e.g., 0.05 for 5% annual sale
	Years        int                // e.g., 5, 10, 15
}
type SimulationResult struct {
	Values [][]float64 // Portfolio values for each simulation run
}
