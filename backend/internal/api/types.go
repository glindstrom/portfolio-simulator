package api

import (
	"errors"
	"math"
	"strings"
)

// AssetRequest defines a single asset within a portfolio, including its ticker and weight.
type AssetRequest struct {
	Ticker string  `json:"ticker"` // Asset identifier (e.g. AAPL, SPY, BTCUSD)
	Weight float64 `json:"weight"` // Portfolio weight (e.g. 0.25 for 25%)
}

type SimulationRequest struct {
	Portfolio   []AssetRequest `json:"portfolio"`      // Now mandatory
	InitialVal  float64        `json:"initialValue"`   // Starting portfolio value
	Withdrawal  float64        `json:"withdrawalRate"` // Annual withdrawal rate (e.g. 0.04 = 4%)
	Inflation   float64        `json:"inflation"`      // Annual inflation rate (e.g. 0.02 = 2%)
	Simulations int            `json:"simulations"`    // Number of simulation paths
	Periods     int            `json:"periods"`        // Number of periods (e.g. months)
	Method      string         `json:"method"`         // "normal" or "bootstrap"
}

// Validate checks the SimulationRequest for correctness and completeness.
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

	if len(r.Portfolio) == 0 {
		return errors.New("portfolio must be provided and cannot be empty")
	}

	totalWeight := 0.0
	for _, a := range r.Portfolio {
		if a.Ticker == "" {
			return errors.New("each asset in portfolio must have a ticker")
		}
		if a.Weight <= 0 {
			return errors.New("asset weights must be greater than 0")
		}
		if a.Weight > 1.0 {
			return errors.New("individual asset weight cannot exceed 1.0 (100%)")
		}
		// AssetType checks are removed.
		totalWeight += a.Weight
	}
	if math.Abs(totalWeight-1.0) > 0.01 {
		return errors.New("sum of portfolio weights must be approximately 1.0")
	}

	if method := strings.ToLower(r.Method); method != "normal" && method != "bootstrap" {
		return errors.New("method must be 'normal' or 'bootstrap'")
	}

	return nil
}

// SummaryStatsResponse and SimulationResponse remain the same.
type SummaryStatsResponse struct {
	Mean   float64 `json:"mean"`
	Median float64 `json:"median"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
}

type SimulationResponse struct {
	Paths         [][]float64          `json:"paths"`
	FinalStats    SummaryStatsResponse `json:"finalStats"`
	SuccessRate   float64              `json:"successRate"`
	SimulatedCAGR float64              `json:"simulatedCAGR"`
}
