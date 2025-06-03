package tiingo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"portfolio-simulator/backend/internal/data" // Adjust import path as per your module structure
)

// Service fetches historical price data from the Tiingo API.
type Service struct {
	APIKey string
	Client *http.Client
}

// NewService creates a new Tiingo Service instance.
// It reads the TIINGO_API_KEY from environment variables.
func NewService() *Service {
	apiKey := os.Getenv("TIINGO_API_KEY")
	if apiKey == "" {
		log.Println("Warning: TIINGO_API_KEY environment variable not set. Tiingo API calls may fail.")
	}
	return &Service{
		APIKey: apiKey,
		Client: http.DefaultClient,
	}
}

// GetMonthlyReturns fetches and calculates monthly percentage returns for a given ticker.
func (s *Service) GetMonthlyReturns(ticker string) ([]float64, error) {
	prices, err := s.GetMonthlyPrices(ticker)
	if err != nil {
		return nil, fmt.Errorf("tiingo: GetMonthlyPrices for %s failed: %w", ticker, err)
	}
	return data.ToMonthlyReturns(prices), nil
}

// GetMonthlyPrices fetches historical monthly closing prices for a ticker from Tiingo.
// Data is resampled to monthly frequency by the Tiingo API.
func (s *Service) GetMonthlyPrices(ticker string) ([]data.PriceData, error) {
	if s.APIKey == "" {
		return nil, fmt.Errorf("tiingo: API key is not configured")
	}

	startDate := "2000-01-01"
	endDate := time.Now().Format("2006-01-02")

	url := fmt.Sprintf("https://api.tiingo.com/tiingo/daily/%s/prices?startDate=%s&endDate=%s&resampleFreq=monthly&token=%s",
		ticker, startDate, endDate, s.APIKey)

	resp, err := s.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("tiingo: request for %s failed: %w", ticker, err)
	}
	defer func() {
		errClose := resp.Body.Close()
		if errClose != nil {
			log.Printf("tiingo: failed to close response body for %s: %v", ticker, errClose)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		// TODO: Attempt to read and log the error message from Tiingo's response body.
		return nil, fmt.Errorf("tiingo: unexpected status code %d for %s", resp.StatusCode, ticker)
	}

	var rawTiingoPrices []struct {
		Date     string  `json:"date"`
		AdjClose float64 `json:"adjClose"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&rawTiingoPrices); err != nil {
		return nil, fmt.Errorf("tiingo: failed to decode JSON response for %s: %w", ticker, err)
	}

	if len(rawTiingoPrices) == 0 {
		return nil, fmt.Errorf("tiingo: no price data returned for %s (period: %s to %s)", ticker, startDate, endDate)
	}

	var prices []data.PriceData
	for _, rawPrice := range rawTiingoPrices {
		parsedDate, err := time.Parse(time.RFC3339, rawPrice.Date)
		if err != nil {
			log.Printf("tiingo: skipping price record for %s due to unparseable date '%s': %v", ticker, rawPrice.Date, err)
			continue
		}
		prices = append(prices, data.PriceData{
			Date:  parsedDate,
			Close: rawPrice.AdjClose,
		})
	}
	return prices, nil
}
