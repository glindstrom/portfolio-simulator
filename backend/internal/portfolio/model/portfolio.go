package model

// Asset represents a single asset within a portfolio, defined by its
// ticker symbol and its allocation weight.
type Asset struct {
	Ticker string  // Ticker symbol (e.g., "AAPL", "SPY").
	Weight float64 // Allocation weight within the portfolio (e.g., 0.6 for 60%).
}

// Portfolio represents a collection of assets.
type Portfolio struct {
	Assets []Asset
}
