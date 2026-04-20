package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Health endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status": "OK",
	}
	json.NewEncoder(w).Encode(response)
}

// Greet endpoint
func greetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	message := fmt.Sprintf("Hello Naman %s!", name)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message": message,
	}
	json.NewEncoder(w).Encode(response)
}
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/greet", greetHandler)

	fmt.Println("server started on port 8080")
	http.ListenAndServe(":8080", mux)
}
