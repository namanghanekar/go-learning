package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type Handler struct {
	repo UserRepo
}

// API: /user?id=1
func (h Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	user := h.repo.GetUser(id)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, user)
}
