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
	Weights      map[string]float64
	SellRate     float64
	TaxRate      float64
	Years        int
}
type SimulationResult struct {
	Values [][]float64 // Portfolio values for each simulation run
}
