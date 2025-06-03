package main

import (
	"log"
	"net/http"

	"portfolio-simulator/backend/internal/api"
	"portfolio-simulator/backend/internal/data/tiingo"
)

// corsMiddleware applies development CORS headers.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Make allowed origin configurable for different environments.
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
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
	priceFetcherSvc := tiingo.NewService()

	apiHandler := &api.Handler{
		Fetcher: priceFetcherSvc,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/simulate", apiHandler.RunSimulation)

	log.Println("Server starting on http://localhost:8085")
	if err := http.ListenAndServe(":8085", corsMiddleware(mux)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
