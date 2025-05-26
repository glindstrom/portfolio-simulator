package tiingo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type PriceData struct {
	Date  time.Time
	Close float64
}

// TiingoService fetches price data from Tiingo.
type TiingoService struct {
	APIKey string
	Client *http.Client
}

// NewTiingoService returns a new instance of TiingoService.
func NewTiingoService() *TiingoService {
	return &TiingoService{
		APIKey: os.Getenv("TIINGO_API_KEY"),
		Client: http.DefaultClient,
	}
}

// GetMonthlyReturns fetches monthly closing prices and returns percentage returns.
func (s *TiingoService) GetMonthlyReturns(ticker string) ([]float64, error) {
	prices, err := s.GetMonthlyPrices(ticker)
	if err != nil {
		return nil, err
	}
	return ToMonthlyReturns(prices), nil
}

// GetMonthlyPrices fetches monthly close prices for a given ticker.
func (s *TiingoService) GetMonthlyPrices(ticker string) ([]PriceData, error) {
	start := "2000-01-01"
	end := time.Now().Format("2006-01-02")
	url := fmt.Sprintf("https://api.tiingo.com/tiingo/daily/%s/prices?startDate=%s&endDate=%s&resampleFreq=monthly&token=%s",
		ticker, start, end, s.APIKey)

	resp, err := s.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch prices: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var raw []struct {
		Date  string  `json:"date"`
		Close float64 `json:"adjClose"`
	}

	err = json.NewDecoder(resp.Body).Decode(&raw)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var prices []PriceData
	for _, r := range raw {
		t, err := time.Parse(time.RFC3339, r.Date)
		if err != nil {
			continue
		}
		prices = append(prices, PriceData{
			Date:  t,
			Close: r.Close,
		})
	}
	return prices, nil
}
