package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"math" // Import math package for Pow
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// mockFetcher implements the PriceFetcher interface for mocking.
type mockFetcher struct {
	returns []float64
	err     error
}

func (m *mockFetcher) GetMonthlyReturns(ticker string) ([]float64, error) {
	return m.returns, m.err
}

func TestRunSimulation_Normal(t *testing.T) {
	mock := &mockFetcher{
		returns: []float64{0.01, 0.015, -0.005, 0.02, 0.0},
	}

	handler := &Handler{Fetcher: mock}

	reqBody := SimulationRequest{
		Portfolio: []AssetRequest{
			{Ticker: "MOCK_NORMAL", Weight: 1.0},
		},
		InitialVal:  1000,
		Periods:     3, // 3 months
		Simulations: 5,
		Method:      "normal",
		// WithdrawalRate and Inflation default to 0.0
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err, "Failed to marshal request body")

	req := httptest.NewRequest(http.MethodPost, "/api/simulate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.RunSimulation(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, "Expected status OK. Body: %s", rr.Body.String())

	var resp SimulationResponse
	err = json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err, "Failed to decode response")

	require.Len(t, resp.Paths, reqBody.Simulations, "Number of simulation paths mismatch")
	if len(resp.Paths) > 0 {
		require.Len(t, resp.Paths[0], reqBody.Periods+1, "Path length mismatch for simulation paths")
	}
	require.NotEqual(t, 0.0, resp.FinalStats.Mean, "Expected non-zero mean for successful simulation")

	// Test SimulatedCAGR
	if reqBody.InitialVal > 0 && reqBody.Periods > 0 {
		years := float64(reqBody.Periods) / 12.0
		require.True(t, years > 0, "Years must be positive for CAGR calculation")

		var expectedCAGR float64
		if resp.FinalStats.Mean >= 0 { // Ensure mean is not negative before Pow
			expectedCAGR = math.Pow(resp.FinalStats.Mean/reqBody.InitialVal, 1.0/years) - 1.0
		} else if reqBody.InitialVal > 0 { // If mean is negative, CAGR effectively -100% or worse
			expectedCAGR = -1.0
		}
		require.InDelta(t, expectedCAGR, resp.SimulatedCAGR, 0.0001, "SimulatedCAGR calculation mismatch")
	} else {
		// If initial value or periods are zero/negative (though validation should prevent this),
		// CAGR might be zero or not calculated.
		require.Equal(t, 0.0, resp.SimulatedCAGR, "Expected CAGR to be 0.0 if InitialVal or Periods are not positive")
	}
}

func TestRunSimulation_BootstrapMethod(t *testing.T) {
	mock := &mockFetcher{
		returns: []float64{0.01, 0.015, -0.005, 0.02, 0.03, -0.01},
	}
	handler := &Handler{Fetcher: mock}

	reqBody := SimulationRequest{
		Portfolio:   []AssetRequest{{Ticker: "MOCK_BOOT", Weight: 1.0}},
		InitialVal:  5000,
		Periods:     4, // 4 months
		Simulations: 3,
		Method:      "bootstrap",
		// WithdrawalRate and Inflation default to 0.0
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err, "Failed to marshal request body for bootstrap test")

	req := httptest.NewRequest(http.MethodPost, "/api/simulate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.RunSimulation(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, "Expected status OK for bootstrap. Body: %s", rr.Body.String())

	var resp SimulationResponse
	err = json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err, "Failed to decode bootstrap response")

	require.Len(t, resp.Paths, reqBody.Simulations, "Number of bootstrap simulation paths mismatch")
	if len(resp.Paths) > 0 {
		require.Len(t, resp.Paths[0], reqBody.Periods+1, "Path length mismatch for bootstrap simulation paths")
	}
	require.NotEqual(t, 0.0, resp.FinalStats.Mean, "Expected non-zero mean for successful bootstrap simulation")

	// Test SimulatedCAGR
	if reqBody.InitialVal > 0 && reqBody.Periods > 0 {
		years := float64(reqBody.Periods) / 12.0
		require.True(t, years > 0, "Years must be positive for CAGR calculation")

		var expectedCAGR float64
		if resp.FinalStats.Mean >= 0 {
			expectedCAGR = math.Pow(resp.FinalStats.Mean/reqBody.InitialVal, 1.0/years) - 1.0
		} else if reqBody.InitialVal > 0 {
			expectedCAGR = -1.0
		}
		require.InDelta(t, expectedCAGR, resp.SimulatedCAGR, 0.0001, "SimulatedCAGR calculation mismatch")
	} else {
		require.Equal(t, 0.0, resp.SimulatedCAGR, "Expected CAGR to be 0.0 if InitialVal or Periods are not positive")
	}
}

// TestRunSimulation_FailureInFetcher, TestRunSimulation_InvalidRequestBody, TestRunSimulation_ValidationErrors, TestRunSimulation_NoReturnsFromFetcher
// remain the same as they don't reach the point of successful response decoding for CAGR.
// ... (rest of the test file as previously provided) ...
func TestRunSimulation_FailureInFetcher(t *testing.T) {
	mock := &mockFetcher{
		err: errors.New("API provider is down"),
	}
	handler := &Handler{Fetcher: mock}

	reqBody := SimulationRequest{
		Portfolio: []AssetRequest{
			{Ticker: "MOCK_FAIL", Weight: 1.0},
		},
		InitialVal:  1000,
		Periods:     3,
		Simulations: 5,
		Method:      "normal",
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err, "Failed to marshal request body")

	req := httptest.NewRequest(http.MethodPost, "/api/simulate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.RunSimulation(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status Internal Server Error. Body: %s", rr.Body.String())
	require.Contains(t, rr.Body.String(), "Failed to fetch returns for ticker MOCK_FAIL", "Response body should contain fetch error")
}

func TestRunSimulation_InvalidRequestBody(t *testing.T) {
	handler := &Handler{Fetcher: &mockFetcher{}} // Fetcher setup doesn't matter much here

	malformedJSON := `{"initialValue": 1000, "periods": 12, "portfolio": [`                        // Missing closing bracket and brace
	req := httptest.NewRequest(http.MethodPost, "/api/simulate", strings.NewReader(malformedJSON)) // Use strings.Reader
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.RunSimulation(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code, "Expected BadRequest for malformed JSON. Body: %s", rr.Body.String())
	require.Contains(t, rr.Body.String(), "Invalid request body", "Response body should indicate invalid request")
}

func TestRunSimulation_ValidationErrors(t *testing.T) {
	validPortfolio := []AssetRequest{{Ticker: "VALID", Weight: 1.0}}
	minimalValidReturns := []float64{0.01}
	fetcherForValidationTests := &mockFetcher{returns: minimalValidReturns}

	baseRequest := func() SimulationRequest {
		return SimulationRequest{
			Portfolio:   validPortfolio,
			InitialVal:  1000,
			Periods:     12,
			Simulations: 10,
			Method:      "normal",
			Withdrawal:  0.04,
			Inflation:   0.02,
		}
	}

	testCases := []struct {
		name          string
		modifier      func(r *SimulationRequest)
		expectedError string
	}{
		{"zero initial value", func(r *SimulationRequest) { r.InitialVal = 0 }, "initial value must be greater than 0"},
		{"negative initial value", func(r *SimulationRequest) { r.InitialVal = -100 }, "initial value must be greater than 0"},
		{"zero periods", func(r *SimulationRequest) { r.Periods = 0 }, "periods must be between 1 and 1200"},
		{"periods too high", func(r *SimulationRequest) { r.Periods = 1201 }, "periods must be between 1 and 1200"},
		{"zero simulations", func(r *SimulationRequest) { r.Simulations = 0 }, "simulations must be between 1 and 10000"},
		{"simulations too high", func(r *SimulationRequest) { r.Simulations = 10001 }, "simulations must be between 1 and 10000"},
		{"negative withdrawal", func(r *SimulationRequest) { r.Withdrawal = -0.1 }, "withdrawal rate must be between 0 and 1"},
		{"withdrawal too high", func(r *SimulationRequest) { r.Withdrawal = 1.1 }, "withdrawal rate must be between 0 and 1"},
		{"negative inflation", func(r *SimulationRequest) { r.Inflation = -0.1 }, "inflation must be between 0 and 1"},
		{"inflation too high", func(r *SimulationRequest) { r.Inflation = 1.1 }, "inflation must be between 0 and 1"},
		{"empty portfolio", func(r *SimulationRequest) { r.Portfolio = []AssetRequest{} }, "portfolio must be provided and cannot be empty"},
		{"asset with empty ticker", func(r *SimulationRequest) { r.Portfolio = []AssetRequest{{Ticker: "", Weight: 1.0}} }, "each asset in portfolio must have a ticker"},
		{"asset weight zero", func(r *SimulationRequest) { r.Portfolio = []AssetRequest{{Ticker: "T", Weight: 0.0}} }, "asset weights must be greater than 0"},
		{"asset weight negative", func(r *SimulationRequest) { r.Portfolio = []AssetRequest{{Ticker: "T", Weight: -0.5}} }, "asset weights must be greater than 0"},
		{"asset weight too high", func(r *SimulationRequest) { r.Portfolio = []AssetRequest{{Ticker: "T", Weight: 1.1}} }, "individual asset weight cannot exceed 1.0"},
		{"portfolio weights sum not 1", func(r *SimulationRequest) {
			r.Portfolio = []AssetRequest{{Ticker: "T1", Weight: 0.5}, {Ticker: "T2", Weight: 0.6}}
		}, "sum of portfolio weights must be approximately 1.0"},
		{"invalid method", func(r *SimulationRequest) { r.Method = "unknown" }, "method must be 'normal' or 'bootstrap'"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := &Handler{Fetcher: fetcherForValidationTests}
			reqBody := baseRequest()
			tc.modifier(&reqBody)

			body, err := json.Marshal(reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/simulate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.RunSimulation(rr, req)

			require.Equal(t, http.StatusBadRequest, rr.Code, "Expected BadRequest for validation error '%s'. Body: %s", tc.name, rr.Body.String())
			require.Contains(t, rr.Body.String(), tc.expectedError, "Response body for '%s' does not contain expected error message", tc.name)
		})
	}
}

func TestRunSimulation_NoReturnsFromFetcher(t *testing.T) {
	mock := &mockFetcher{
		returns: []float64{},
		err:     nil,
	}
	handler := &Handler{Fetcher: mock}

	reqBody := SimulationRequest{
		Portfolio:   []AssetRequest{{Ticker: "MOCK_NO_RETURNS", Weight: 1.0}},
		InitialVal:  1000,
		Periods:     3,
		Simulations: 5,
		Method:      "normal",
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/simulate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.RunSimulation(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code, "Expected InternalServerError when fetcher returns no data. Body: %s", rr.Body.String())
	require.Contains(t, rr.Body.String(), "returns slice is empty", "Response body should indicate empty returns issue for simulation")
}
