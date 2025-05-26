// backend/cmd/server/main.go

package main

import (
	"log"
	"net/http"

	"portfolio-simulator/backend/internal/api"
	"portfolio-simulator/backend/internal/data/tiingo"
)

func main() {
	tiingoSvc := tiingo.NewTiingoService()
	handler := &api.Handler{Tiingo: tiingoSvc}

	http.HandleFunc("/api/simulate", handler.RunSimulation)

	log.Println("Server running on :8085")
	log.Fatal(http.ListenAndServe(":8085", nil))
}
