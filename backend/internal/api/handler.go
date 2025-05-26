package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"portfolio-simulator/backend/internal/portfolio"
	"portfolio-simulator/backend/internal/portfolio/model"
	"portfolio-simulator/backend/internal/simulation"
)

type TiingoFetcher interface {
	GetMonthlyReturns(ticker string) ([]float64, error)
}

type Handler struct {
	Tiingo TiingoFetcher
}

func (h *Handler) RunSimulation(w http.ResponseWriter, r *http.Request) {
	var req SimulationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var p model.Portfolio
	if len(req.Portfolio) > 0 {
		for _, ar := range req.Portfolio {
			p.Assets = append(p.Assets, model.Asset{
				Ticker: ar.Ticker,
				Weight: ar.Weight,
			})
		}
	} else if req.Ticker != "" {
		p.Assets = append(p.Assets, model.Asset{
			Ticker: req.Ticker,
			Weight: 1.0,
		})
	} else {
		http.Error(w, "Missing portfolio or ticker", http.StatusBadRequest)
		return
	}

	returns, err := portfolio.ComputePortfolioReturns(p, h.Tiingo.GetMonthlyReturns)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch returns: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("Fetched %d months of returns", len(returns))

	params := simulation.Params{
		InitialValue:     req.InitialVal,
		Returns:          returns,
		WithdrawalRate:   req.Withdrawal,
		InflationPerYear: req.Inflation,
		Periods:          req.Periods,
		Simulations:      req.Simulations,
	}

	var simResult *simulation.Result

	switch req.Method {
	case "normal":
		simResult, err = simulation.SimulateNormal(params)
	case "bootstrap":
		simResult, err = simulation.SimulateBootstrap(params)
	default:
		http.Error(w, "Unknown simulation method", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Simulation error: %v", err), http.StatusInternalServerError)
		return
	}

	resp := SimulationResponse{
		Paths: simResult.Paths,
		FinalStats: SummaryStatsResponse{
			Mean:   simResult.FinalStats.Mean,
			Median: simResult.FinalStats.Median,
			Min:    simResult.FinalStats.Min,
			Max:    simResult.FinalStats.Max,
		},
		SuccessRate: simResult.SuccessRate,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
