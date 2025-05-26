//go:build integration

package tiingo_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	service       *TiingoService
	cachedReturns map[string][]float64
)

func TestMain(m *testing.M) {
	service = NewTiingoService()
	cachedReturns = make(map[string][]float64)
	os.Exit(m.Run())
}

func getReturns(t *testing.T, ticker string) []float64 {
	if returns, ok := cachedReturns[ticker]; ok {
		return returns
	}

	returns, err := service.GetMonthlyReturns(ticker)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(returns), 10)
	cachedReturns[ticker] = returns
	return returns
}

func TestGetMonthlyReturns_SPY(t *testing.T) {
	returns := getReturns(t, "SPY")
	require.NotEmpty(t, returns)
}

func TestGetMonthlyReturns_BTC(t *testing.T) {
	returns := getReturns(t, "btcusd")
	require.NotEmpty(t, returns)
}
