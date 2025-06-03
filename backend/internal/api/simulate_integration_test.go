//go:build integration

package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	// Assuming these are the correct import paths based on your project structure
	"portfolio-simulator/backend/internal/api"
	"portfolio-simulator/backend/internal/data/tiingo"

	"github.com/stretchr/testify/require" // Import testify/require
)

func TestRunSimulation_WithPortfolio_NormalAndBootstrap(t *testing.T) {
	if os.Getenv("TIINGO_API_KEY") == "" {
		t.Skip("TIINGO_API_KEY not set; skipping integration test")
	}

	handler := &api.Handler{
		Fetcher: tiingo.NewService(),
	}

	methods := []string{"normal", "bootstrap"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			reqBody := api.SimulationRequest{
				Portfolio: []api.AssetRequest{
					{Ticker: "AAPL", Weight: 0.1},
					{Ticker: "MSFT", Weight: 0.1},
					{Ticker: "SPY", Weight: 0.8}, // SPY is an ETF, testing Tiingo's capability
				},
				InitialVal:  10000,
				Periods:     12, // 1 year of monthly periods
				Simulations: 100,
				Method:      method,
				Withdrawal:  0.04, // 4% annual withdrawal rate
				Inflation:   0.02, // 2% annual inflation rate
			}

			body, err := json.Marshal(reqBody)
			require.NoError(t, err, "Failed to marshal request for method %s", method)

			req := httptest.NewRequest(http.MethodPost, "/api/simulate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.RunSimulation(rr, req)

			require.Equal(t, http.StatusOK, rr.Code, "[%s] Expected status 200 OK, got %d. Body: %s", method, rr.Code, rr.Body.String())

			var resp api.SimulationResponse
			err = json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err, "[%s] Failed to decode response", method)

			// Log basic statistics for manual inspection if needed
			t.Logf("[%s] Mean: %.2f, Median: %.2f, Min: %.2f, Max: %.2f, SuccessRate: %.2f%%",
				method,
				resp.FinalStats.Mean,
				resp.FinalStats.Median,
				resp.FinalStats.Min,
				resp.FinalStats.Max,
				resp.SuccessRate*100, // Display success rate as percentage
			)

			require.Len(t, resp.Paths, 100, "[%s] Expected 100 simulation paths", method)
			if len(resp.Paths) > 0 {
				require.Len(t, resp.Paths[0], reqBody.Periods+1, "[%s] Expected path length to be periods+1", method)
			}

			// SuccessRate should be between 0 and 1 (inclusive)
			require.GreaterOrEqual(t, resp.SuccessRate, 0.0, "[%s] Success rate should be >= 0", method)
			require.LessOrEqual(t, resp.SuccessRate, 1.0, "[%s] Success rate should be <= 1", method)

			// Basic sanity checks for financial stats

			if resp.SuccessRate > 0 { // Only check these if there's some success, otherwise Min/Max/Mean might be 0 or negative due to depletion
				require.True(t, resp.FinalStats.Max >= resp.FinalStats.Min, "[%s] Max should be >= Min", method)
			}
		})
	}
}
