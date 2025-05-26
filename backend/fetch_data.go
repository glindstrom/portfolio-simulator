package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"portfolio-simulator/backend/models"
	"portfolio-simulator/backend/repository"
	"portfolio-simulator/backend/service"
	"strconv"
	"strings"
	"time"
)

type AlphaVantageResponse struct {
	TimeSeries map[string]struct {
		Close string `json:"4. close"`
	} `json:"Time Series (Daily)"`
}
type CryptoResponse struct {
	TimeSeries map[string]struct {
		Close string `json:"4. close"`
	} `json:"Time Series (Digital Currency Daily)"`
}

func fetchData(symbol, apiKey, function, market string, store repository.DataStore) error {
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=%s&symbol=%s&market=%s&apikey=%s", function, symbol, market, apiKey)
	if function != "DIGITAL_CURRENCY_DAILY" {
		url = fmt.Sprintf("https://www.alphavantage.co/query?function=%s&symbol=%s&apikey=%s&outputsize=compact", function, symbol, apiKey)
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if strings.Contains(function, "DIGITAL_CURRENCY") {
		var data map[string]interface{}
		if err := json.Unmarshal(rawData, &data); err != nil {
			return err
		}
		log.Printf("BTC Raw Response: %+v\n", data)
		var cryptoData CryptoResponse
		if err := json.Unmarshal(rawData, &cryptoData); err != nil {
			return err
		}
		var prices []models.Price
		for date, price := range cryptoData.TimeSeries {
			t, err := time.Parse("2006-01-02", date)
			if err != nil {
				log.Printf("Error parsing date %s: %v", date, err)
				continue
			}
			close, err := strconv.ParseFloat(price.Close, 64)
			if err != nil {
				log.Printf("Error parsing price %s: %v", price.Close, err)
				continue
			}
			prices = append(prices, models.Price{Date: t, Close: close})
		}
		if len(prices) == 0 {
			log.Printf("No valid prices for %s", symbol)
		}
		return store.SavePrices(symbol, prices)
	}
	var data AlphaVantageResponse
	if err := json.Unmarshal(rawData, &data); err != nil {
		return err
	}
	var prices []models.Price
	for date, price := range data.TimeSeries {
		t, _ := time.Parse("2006-01-02", date)
		close, _ := strconv.ParseFloat(price.Close, 64)
		prices = append(prices, models.Price{Date: t, Close: close})
	}
	return store.SavePrices(symbol, prices)
}
func main() {
	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if apiKey == "" {
		log.Fatal("ALPHA_VANTAGE_API_KEY not set")
	}
	store := repository.NewCSVStore("../data")
	dataService := service.NewDataService(store)
	simService := service.NewSimulationService(dataService)
	symbols := []struct {
		symbol   string
		function string
		market   string
	}{
		{"BTC", "DIGITAL_CURRENCY_DAILY", "USD"},
		{"EUNL.DE", "TIME_SERIES_DAILY", ""},
	}
	for _, s := range symbols {
		if _, err := os.Stat("../data/" + s.symbol + ".csv"); os.IsNotExist(err) {
			if err := fetchData(s.symbol, apiKey, s.function, s.market, store); err != nil {
				log.Printf("Error fetching %s: %v", s.symbol, err)
			}
			time.Sleep(12 * time.Second)
		}
	}
	portfolio := models.Portfolio{
		InitialValue: 100000,
		Weights:      map[string]float64{"BTC": 0.05, "EUNL.DE": 0.95},
		SellRate:     0.05,
		Years:        10,
	}
	result, err := simService.RunSimulation(portfolio, 1000)
	if err != nil {
		log.Fatal(err)
	}
	mean := 0.0
	for _, values := range result.Values {
		mean += values[len(values)-1]
	}
	mean /= float64(len(result.Values))
	fmt.Printf("Mean portfolio value after %d years: %.2f EUR\n", portfolio.Years, mean)
}
