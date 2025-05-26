//go:build integration
// +build integration

package tiingo

import (
	"os"
	"testing"
	"time"
)

func TestGetDailyPrices(t *testing.T) {
	if os.Getenv("TIINGO_API_KEY") == "" {
		t.Skip("TIINGO_API_KEY not set")
	}

	start := time.Now().AddDate(0, 0, -30)
	end := time.Now()

	data, err := GetDailyPrices("AAPL", start, end)
	if err != nil {
		t.Fatalf("Failed to get daily prices: %v", err)
	}

	if len(data) < 10 {
		t.Fatalf("Expected at least 10 price entries, got %d", len(data))
	}

	for _, d := range data {
		t.Logf("%s: %.2f", d.Date.Format("2006-01-02"), d.Close)
	}
}
