package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	addr := ":8085"
	log.Printf("Server running at %s...", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
