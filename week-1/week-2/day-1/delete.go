package main

import (
	"encoding/json"
	"net/http"
)

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {

	// Step 1: Check method
	if r.Method != "DELETE" {
		http.Error(w, "Only DELETE allowed", http.StatusMethodNotAllowed)
		return
	}

	// Step 2: Get ID from query
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	// Step 3: Perform delete logic (dummy here)

	// Step 4: Send response
	response := map[string]string{
		"message": "User deleted",
		"id":      id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/user", deleteUserHandler)
	http.ListenAndServe(":9090", nil)
}
