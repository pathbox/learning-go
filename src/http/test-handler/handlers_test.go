// handlers_test.go
package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
  req, err := http.NewRequest("GET", "/health-check", nil)
  if err != nil {
    t.Fatal(err)
  }

  rr := httptest.NewRecorder()
  handler := http.HandlerFunc(HealthCheckHandler)
  handler.ServeHTTP(rr, req)

  if status := rr.Code; status != http.StatusOK {
    t.Error("handler returned wrong status code: got %v want %v",
    status, http.StatusOK)
  }

  expected := `{"alive": true}`
  if rr.Body.String() != expected {
    t.Errorf("handler returned unexpected body: got %v want %v",
    rr.Body.String(), expected)
  }
}
