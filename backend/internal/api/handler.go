package api

import (
	"encoding/json"
	"net/http"

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
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	returns, err := h.Tiingo.GetMonthlyReturns(req.Ticker)
	if err != nil {
		http.Error(w, "failed to fetch returns: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sim := &simulation.Simulator{
		InitialValue: req.InitialVal,
		Returns:      returns,
	}

	var result *simulation.SimulationResult
	if req.Method == "bootstrap" {
		result, err = sim.BootstrapSim(req.Simulations, req.Periods)
	} else {
		result, err = sim.Simulate(req.Simulations, req.Periods)
	}
	if err != nil {
		http.Error(w, "simulation error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := SimulationResponse{
		Paths: result.Paths,
	}
	resp.FinalStats.Mean = result.FinalStats.Mean
	resp.FinalStats.Median = result.FinalStats.Median
	resp.FinalStats.Min = result.FinalStats.Min
	resp.FinalStats.Max = result.FinalStats.Max

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
