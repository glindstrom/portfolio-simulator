package api

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"

	"portfolio-simulator/backend/internal/portfolio"
	"portfolio-simulator/backend/internal/portfolio/model"
	"portfolio-simulator/backend/internal/simulation"
)

// PriceFetcher defines the interface for fetching monthly returns for a ticker.
type PriceFetcher interface {
	GetMonthlyReturns(ticker string) ([]float64, error)
}

// Handler holds dependencies for API handlers, such as data fetchers.
type Handler struct {
	Fetcher PriceFetcher // Consolidated to a single fetcher.
}

// RunSimulation handles requests to run a portfolio simulation.
func (h *Handler) RunSimulation(w http.ResponseWriter, r *http.Request) {
	if h.Fetcher == nil {
		log.Println("Error: Handler's PriceFetcher is not initialized.")
		http.Error(w, "Internal server error: Fetcher service not available", http.StatusInternalServerError)
		return
	}

	var req SimulationRequest // Defined in types.go; AssetRequest within it no longer has AssetType.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var p model.Portfolio
	for _, ar := range req.Portfolio {
		p.Assets = append(p.Assets, model.Asset{
			Ticker: ar.Ticker,
			Weight: ar.Weight,
		})
	}

	returnsByAsset := make(map[string][]float64)
	for _, asset := range p.Assets {
		log.Printf("Fetching returns for ticker: %s", asset.Ticker)
		assetReturns, fetchErr := h.Fetcher.GetMonthlyReturns(asset.Ticker)
		if fetchErr != nil {
			log.Printf("Error fetching returns for %s: %v", asset.Ticker, fetchErr)
			http.Error(w, fmt.Sprintf("Failed to fetch returns for ticker %s", asset.Ticker), http.StatusInternalServerError)
			return
		}
		// It's possible a fetcher returns no error but also no returns (e.g., new ticker with no history).
		if len(assetReturns) == 0 {
			log.Printf("Warning: No returns fetched for %s. This might affect simulation results if its weight is > 0.", asset.Ticker)
		}
		returnsByAsset[asset.Ticker] = assetReturns
	}

	portfolioReturns, err := portfolio.WeightedMonthlyReturns(p, returnsByAsset)
	if err != nil {
		log.Printf("Error computing weighted portfolio returns: %v", err)
		http.Error(w, "Failed to compute weighted portfolio returns", http.StatusInternalServerError)
		return
	}

	// The simulation functions expect a non-empty returns slice if the portfolio is non-empty and assets are valid.
	if len(portfolioReturns) == 0 && len(p.Assets) > 0 {
		log.Println("Warning: Weighted portfolio returns are empty, though assets were provided. Check asset return data, lengths, and weights.")
		// SimulateNormal/Bootstrap will error if params.Returns is empty.
		// We can pre-emptively return an error or let the simulation functions handle it.
		// For now, let the simulation functions handle it, as they have specific error messages.
	}
	log.Printf("Computed portfolio returns for %d months.", len(portfolioReturns))

	params := simulation.Params{
		InitialValue:     req.InitialVal,
		Returns:          portfolioReturns,
		WithdrawalRate:   req.Withdrawal, // This is withdrawalRate from request
		InflationPerYear: req.Inflation,
		Periods:          req.Periods,
		Simulations:      req.Simulations,
	}

	var simResult *simulation.Result
	// req.Method is already validated to be "normal" or "bootstrap"
	switch req.Method {
	case "normal":
		simResult, err = simulation.SimulateNormal(params)
	case "bootstrap":
		simResult, err = simulation.SimulateBootstrap(params)
	}

	if err != nil {
		log.Printf("Simulation error (method: %s): %v", req.Method, err)
		http.Error(w, fmt.Sprintf("Simulation error: %v", err), http.StatusInternalServerError)
		return
	}

	if simResult == nil {
		log.Printf("Error: Simulation completed without error, but simResult is nil (method: %s)", req.Method)
		http.Error(w, "Internal server error: Simulation returned no result", http.StatusInternalServerError)
		return
	}

	var simulatedCAGR float64
	if req.InitialVal > 0 && params.Periods > 0 {
		years := float64(params.Periods) / 12.0
		if years > 0 {
			if simResult.FinalStats.Mean >= 0 {
				simulatedCAGR = math.Pow(simResult.FinalStats.Mean/req.InitialVal, 1.0/years) - 1.0
			} else {
				if req.InitialVal > 0 {
					simulatedCAGR = -1.0 // Represents 100% loss or more if mean becomes negative
				}
			}
		}
	}

	resp := SimulationResponse{
		Paths:         simResult.Paths,
		FinalStats:    SummaryStatsResponse(simResult.FinalStats),
		SuccessRate:   simResult.SuccessRate,
		SimulatedCAGR: simulatedCAGR,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding or writing simulation response: %v", err)
	}
}
