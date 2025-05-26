// backend/internal/api/types.go

package api

type SimulationRequest struct {
	Ticker      string  `json:"ticker"`
	InitialVal  float64 `json:"initial_value"`
	Periods     int     `json:"periods"`
	Simulations int     `json:"simulations"`
	Method      string  `json:"method"` // "normal" or "bootstrap"
}

type SimulationResponse struct {
	Paths      [][]float64 `json:"paths"`
	FinalStats struct {
		Mean   float64 `json:"mean"`
		Median float64 `json:"median"`
		Min    float64 `json:"min"`
		Max    float64 `json:"max"`
	} `json:"final_stats"`
}
