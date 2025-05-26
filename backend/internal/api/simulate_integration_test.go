//go:build integration

package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"portfolio-simulator/backend/internal/api"
	"portfolio-simulator/backend/internal/data/tiingo"
	"testing"
)

func TestRunSimulation_WithPortfolio_NormalAndBootstrap(t *testing.T) {
	if os.Getenv("TIINGO_API_KEY") == "" {
		t.Skip("TIINGO_API_KEY not set; skipping integration test")
	}

	handler := &api.Handler{
		Tiingo: tiingo.NewTiingoService(),
	}

	methods := []string{"normal", "bootstrap"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			// Build the simulation request with withdrawal and inflation parameters
			reqBody := api.SimulationRequest{
				Portfolio: []api.AssetRequest{
					{Ticker: "AAPL", Weight: 0.5},
					{Ticker: "MSFT", Weight: 0.3},
					{Ticker: "SPY", Weight: 0.2},
				},
				InitialVal:  10000,
				Periods:     12,
				Simulations: 100,
				Method:      method,
				Withdrawal:  0.04, // 4% annual withdrawal
				Inflation:   0.02, // 2% annual inflation
			}

			body, err := json.Marshal(reqBody)
			if err != nil {
				t.Fatalf("failed to marshal request: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/simulate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.RunSimulation(rr, req)

			if rr.Code != http.StatusOK {
				t.Logf("Body: %s", rr.Body.String())
				t.Fatalf("expected 200 OK, got %d", rr.Code)
			}

			var resp api.SimulationResponse
			if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			// Log basic statistics
			t.Logf("[%s] Mean: %.2f, Median: %.2f, Min: %.2f, Max: %.2f, SuccessRate: %.2f",
				method,
				resp.FinalStats.Mean,
				resp.FinalStats.Median,
				resp.FinalStats.Min,
				resp.FinalStats.Max,
				resp.SuccessRate,
			)

			// Validate that we still have the expected number of paths
			if len(resp.Paths) != 100 {
				t.Errorf("[%s] expected 100 simulation paths, got %d", method, len(resp.Paths))
			}

			// SuccessRate should be between 0 and 1
			if resp.SuccessRate < 0 || resp.SuccessRate > 1 {
				t.Errorf("[%s] unexpected success rate: %.2f", method, resp.SuccessRate)
			}
		})
	}
}
