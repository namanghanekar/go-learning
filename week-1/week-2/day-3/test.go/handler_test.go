package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock database
type MockRepo struct{}

func (m MockRepo) GetUser(id int) string {

	return "Mock User"
}

func TestGetUserHandler(t *testing.T) {
	mockRepo := MockRepo{}

	handler := Handler{
		repo: mockRepo,
	}

	req := httptest.NewRequest("GET", "/user?id=1", nil)
	w := httptest.NewRecorder()

	handler.GetUserHandler(w, req)

	res := w.Result()
	body := w.Body.String()

	// Check status code
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	// Check response body
	if body != "Mock User" {
		t.Errorf("expected Mock User, got %s", body)
	}
}
