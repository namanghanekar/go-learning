package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {

	// Step 1: Check method
	if r.Method != "PUT" {
		http.Error(w, "Only PUT allowed", http.StatusMethodNotAllowed)
		return
	}

	// Step 2: Get ID from query
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	// Step 3: Read JSON body
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Step 4: Create response
	response := map[string]string{
		"message": "User updated",
		"id":      id,
		"name":    user.Name,
		"role":    user.Role,
	}

	// Step 5: Send JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/user", updateUserHandler)
	http.ListenAndServe(":8080", nil)
}
