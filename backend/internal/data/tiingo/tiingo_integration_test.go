//go:build integration

package tiingo_test

import (
	"log"
	"os"
	"testing"

	"portfolio-simulator/backend/internal/data/tiingo" // Import the package to be tested

	"github.com/stretchr/testify/require"
)

var (
	// Use the correct, exported type from the tiingo package.
	// The 'tiingo.' prefix is needed because this is in package 'tiingo_test'.
	testService     *tiingo.Service
	cachedReturns   map[string][]float64
	isAPIKeyPresent bool
)

// TestMain sets up the service for integration tests and handles API key checks.
func TestMain(m *testing.M) {
	apiKey := os.Getenv("TIINGO_API_KEY")
	if apiKey == "" {
		log.Println("WARNING: TIINGO_API_KEY environment variable not set. Integration tests will be skipped.")
		isAPIKeyPresent = false
		// Do not call os.Exit here if you want individual tests to check isAPIKeyPresent and skip.
		// If you want to exit early if no key, then os.Exit(0) could be used, but tests won't report as skipped.
	} else {
		isAPIKeyPresent = true
		// Use the renamed constructor.
		testService = tiingo.NewService()
	}

	cachedReturns = make(map[string][]float64)
	exitCode := m.Run() // Run all tests.
	os.Exit(exitCode)
}

// skipIfNotIntegrationReady checks if the API key is present and skips the test if not.
func skipIfNotIntegrationReady(t *testing.T) {
	if !isAPIKeyPresent {
		t.Skip("Skipping integration test: TIINGO_API_KEY not set.")
	}
	if testService == nil || testService.APIKey == "" { // Double check service was initialized if APIKey was initially present
		t.Fatal("Test service not initialized correctly, API key might be missing or NewService failed silently.")
	}
}

// getReturns is a helper to fetch (and cache) returns for a ticker during tests.
func getReturns(t *testing.T, ticker string) []float64 {
	skipIfNotIntegrationReady(t) // Ensure API key is available before making a call.

	if returns, ok := cachedReturns[ticker]; ok {
		return returns
	}

	t.Logf("Fetching live Tiingo data for ticker: %s", ticker)
	returns, err := testService.GetMonthlyReturns(ticker)
	require.NoError(t, err, "Fetching returns for %s should not produce an error", ticker)

	// It's good to have some data, but exact length can vary.
	// For stable tickers like SPY, a significant history is expected.
	// For newer or less common tickers, this might need adjustment.
	require.NotEmpty(t, returns, "Expected some returns for ticker %s, but got none. Check data source or ticker validity.", ticker)
	if len(returns) < 10 && (ticker == "SPY" || ticker == "AAPL") { // Example for well-established tickers
		t.Logf("Warning: Received only %d return periods for %s, expected more for an established ticker.", len(returns), ticker)
	}

	cachedReturns[ticker] = returns
	return returns
}

// TestGetMonthlyReturns_SPY tests fetching returns for SPY.
func TestGetMonthlyReturns_SPY(t *testing.T) {
	returns := getReturns(t, "SPY")
	require.NotEmpty(t, returns, "Should get returns for SPY")
}

// TestGetMonthlyReturns_AAPL tests fetching returns for AAPL.
func TestGetMonthlyReturns_AAPL(t *testing.T) {
	returns := getReturns(t, "AAPL")
	require.NotEmpty(t, returns, "Should get returns for AAPL")
}

// TestGetMonthlyReturns_BTCUSD tests fetching returns for BTCUSD.
// The success of this test depends on whether the Tiingo endpoint configured in
// TiingoService supports "btcusd" and provides sufficient historical data.
func TestGetMonthlyReturns_BTCUSD(t *testing.T) {
	returns := getReturns(t, "btcusd")
	require.NotEmpty(t, returns, "Should get returns for btcusd if supported")
}

// TestGetMonthlyReturns_InvalidTicker tests behavior for a non-existent ticker.
func TestGetMonthlyReturns_InvalidTicker(t *testing.T) {
	skipIfNotIntegrationReady(t)

	invalidTicker := "THISISNOTAVALIDTICKERXYZ"
	t.Logf("Fetching live Tiingo data for invalid ticker: %s", invalidTicker)
	_, err := testService.GetMonthlyReturns(invalidTicker)

	// Tiingo should return an error for an invalid ticker.
	// The exact error message might vary, so checking for a non-nil error is a good start.
	// Your service's GetMonthlyPrices method should translate Tiingo's error response
	// (e.g., 404 Not Found) into a Go error.
	require.Error(t, err, "Expected an error for an invalid ticker")
	if err != nil {
		// You might want to check for specific error content if Tiingo has consistent error messages.
		// e.g., require.Contains(t, err.Error(), "not found") or similar.
		t.Logf("Received expected error for invalid ticker %s: %v", invalidTicker, err)
	}
}
