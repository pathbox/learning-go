package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	port = "9090"
)

func TestHandler(t *testing.T) {
	expected := []byte("Hello World!")

	req, err := http.NewRequest("GET", buildUrl("/"), nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	handler(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Response code was %v; want 200", res.Code)
	}

	if bytes.Compare(expected, res.Body.Bytes()) != 0 {
		t.Errorf("Response body was '%v'; want '%v'", expected, res.Body)
	}
}

func buildUrl(path string) string {
	return urlFor("http", port, path)
}

func urlFor(scheme string, serverPort string, path string) string {
	return scheme + "://localhost:" + serverPort + path
}
