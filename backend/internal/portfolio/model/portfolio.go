package model

// Asset represents a single portfolio asset with its ticker symbol and allocation weight.
type Asset struct {
	Ticker string  // e.g. "BTC", "EUNL"
	Weight float64 // Allocation as a fraction, e.g. 0.05 for 5%
}

// Portfolio holds multiple assets with their allocations.
type Portfolio struct {
	Assets []Asset
}
