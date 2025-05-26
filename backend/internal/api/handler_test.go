package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockTiingo implements TiingoFetcher interface for mocking
type mockTiingo struct {
	returns []float64
	err     error
}

func (m *mockTiingo) GetMonthlyReturns(ticker string) ([]float64, error) {
	return m.returns, m.err
}

func TestRunSimulation_Normal(t *testing.T) {
	mock := &mockTiingo{
		returns: []float64{0.01, 0.015, -0.005, 0.02, 0.0},
	}

	handler := &Handler{Tiingo: mock}

	reqBody := SimulationRequest{
		Ticker:      "BTC",
		InitialVal:  1000,
		Periods:     3,
		Simulations: 5,
		Method:      "normal",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/simulate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.RunSimulation(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}

	var resp SimulationResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Paths) != 5 {
		t.Errorf("expected 5 simulation paths, got %d", len(resp.Paths))
	}
	if resp.FinalStats.Mean <= 0 {
		t.Errorf("expected positive mean, got %f", resp.FinalStats.Mean)
	}
}

func TestRunSimulation_FailureInTiingo(t *testing.T) {
	mock := &mockTiingo{
		err: errors.New("API down"),
	}
	handler := &Handler{Tiingo: mock}

	reqBody := SimulationRequest{
		Ticker:      "BTC",
		InitialVal:  1000,
		Periods:     3,
		Simulations: 5,
		Method:      "normal",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/simulate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.RunSimulation(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 Internal Server Error, got %d", rr.Code)
	}
}
