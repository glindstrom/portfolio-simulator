package main

import (
	"log"
	"net/http"

	"portfolio-simulator/backend/internal/api"
	"portfolio-simulator/backend/internal/data/tiingo"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // byt ut vid behov
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	tiingoSvc := tiingo.NewTiingoService()
	handler := &api.Handler{Tiingo: tiingoSvc}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/simulate", handler.RunSimulation)

	log.Println("Server running on :8085")
	log.Fatal(http.ListenAndServe(":8085", corsMiddleware(mux)))
}
