package tiingo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type PriceData struct {
	Date  time.Time `json:"date"`
	Close float64   `json:"close"`
}

func GetDailyPrices(symbol string, startDate, endDate time.Time) ([]PriceData, error) {
	apiKey := os.Getenv("TIINGO_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("TIINGO_API_KEY not set")
	}

	url := fmt.Sprintf("https://api.tiingo.com/tiingo/daily/%s/prices?startDate=%s&endDate=%s&token=%s",
		symbol,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
		apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var data []PriceData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return data, nil
}
