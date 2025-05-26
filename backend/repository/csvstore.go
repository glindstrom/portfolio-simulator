package repository

import (
	"encoding/csv"
	"os"
	"portfolio-simulator/backend/models"
	"strconv"
	"time"
)

type CSVStore struct {
	DataDir string
}

func NewCSVStore(dataDir string) *CSVStore {
	return &CSVStore{DataDir: dataDir}
}
func (s *CSVStore) SavePrices(symbol string, prices []models.Price) error {
	file, err := os.Create(s.DataDir + "/" + symbol + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"date", "close"})
	for _, p := range prices {
		writer.Write([]string{p.Date.Format("2006-01-02"), strconv.FormatFloat(p.Close, 'f', 2, 64)})
	}
	return nil
}
func (s *CSVStore) GetPrices(symbol string) ([]models.Price, error) {
	file, err := os.Open(s.DataDir + "/" + symbol + ".csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	var prices []models.Price
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		date, _ := time.Parse("2006-01-02", record[0])
		close, _ := strconv.ParseFloat(record[1], 64)
		prices = append(prices, models.Price{Date: date, Close: close})
	}
	return prices, nil
}
