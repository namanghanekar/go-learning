package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func userHandler(w http.ResponseWriter, r *http.Request) {

	// Step 1: Check method
	if r.Method != "GET" {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	// Step 2: Get query parameter
	id := r.URL.Query().Get("id")

	// Step 3: Validate
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	// Step 4: Create response
	response := map[string]string{
		"id":   id,
		"name": "Naman",
	}

	// Step 5: Set header
	w.Header().Set("Content-Type", "application/json")

	// Step 6: Send response
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/user", userHandler)

	fmt.Println("Server running ")
	http.ListenAndServe(":8080", nil)
}
