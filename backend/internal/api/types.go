// backend/internal/api/types.go

package api

import (
	"errors"
	"math"
	"strings"
)

type AssetRequest struct {
	Ticker string  `json:"ticker"` // Asset identifier (e.g. AAPL)
	Weight float64 `json:"weight"` // Portfolio weight (e.g. 0.25 for 25%)
}

type SimulationRequest struct {
	Ticker      string         `json:"ticker,omitempty"`    // Used if simulating a single asset
	Portfolio   []AssetRequest `json:"portfolio,omitempty"` // Optional: a list of assets with weights
	InitialVal  float64        `json:"initialValue"`        // Starting portfolio value
	Withdrawal  float64        `json:"withdrawalRate"`      // Annual withdrawal rate (e.g. 0.04 = 4%)
	Inflation   float64        `json:"inflation"`           // Annual inflation rate (e.g. 0.02 = 2%)
	Simulations int            `json:"simulations"`         // Number of simulation paths
	Periods     int            `json:"periods"`             // Number of periods (e.g. months)
	Method      string         `json:"method"`              // "normal" or "bootstrap"
}

func (r *SimulationRequest) Validate() error {
	if r.InitialVal <= 0 {
		return errors.New("initial value must be greater than 0")
	}
	if r.Periods <= 0 || r.Periods > 1200 {
		return errors.New("periods must be between 1 and 1200")
	}
	if r.Simulations <= 0 || r.Simulations > 10000 {
		return errors.New("simulations must be between 1 and 10000")
	}
	if r.Withdrawal < 0 || r.Withdrawal > 1 {
		return errors.New("withdrawal rate must be between 0 and 1")
	}
	if r.Inflation < 0 || r.Inflation > 1 {
		return errors.New("inflation must be between 0 and 1")
	}

	// Allow either a portfolio or a single ticker
	if len(r.Portfolio) == 0 && r.Ticker == "" {
		return errors.New("either portfolio or ticker must be provided")
	}

	if len(r.Portfolio) > 0 {
		totalWeight := 0.0
		for _, a := range r.Portfolio {
			if a.Ticker == "" {
				return errors.New("each asset in portfolio must have a ticker")
			}
			if a.Weight <= 0 {
				return errors.New("asset weights must be greater than 0")
			}
			totalWeight += a.Weight
		}
		if math.Abs(totalWeight-1.0) > 0.01 {
			return errors.New("sum of portfolio weights must be approximately 1.0")
		}
	}

	if method := strings.ToLower(r.Method); method != "normal" && method != "bootstrap" {
		return errors.New("method must be 'normal' or 'bootstrap'")
	}

	return nil
}

type SummaryStatsResponse struct {
	Mean   float64 `json:"mean"`   // Average final portfolio value
	Median float64 `json:"median"` // Median final portfolio value
	Min    float64 `json:"min"`    // Minimum final portfolio value
	Max    float64 `json:"max"`    // Maximum final portfolio value
}

type SimulationResponse struct {
	Paths       [][]float64          `json:"paths"`        // Simulated portfolio paths
	FinalStats  SummaryStatsResponse `json:"final_stats"`  // Summary statistics of final values
	SuccessRate float64              `json:"success_rate"` // Proportion of paths that never depleted
}
