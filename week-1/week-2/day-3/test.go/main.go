package main

import (
	"log"
	"net/http"
)

func main() {
	repo := RealRepo{} // database connection or real implementation

	handler := Handler{ // create the handler with the repo
		repo: repo, //assign the repo to the handler
	}

	http.HandleFunc("/user", handler.GetUserHandler)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
