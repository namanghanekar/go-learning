package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// handler function
func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "pong"})
}

func TestStatusCodes(t *testing.T) {

	// set gin test mode
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/ping", pingHandler)

	// table-driven test cases
	tests := []struct {
		name   string
		method string
		url    string
		status int
	}{
		{"Valid route", "GET", "/ping", 200},
		{"Invalid route", "GET", "/wrong", 404},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatalf("error creating request: %v", err)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.status {
				t.Errorf("expected %d, got %d", tt.status, w.Code)
			}
		})
	}
}
