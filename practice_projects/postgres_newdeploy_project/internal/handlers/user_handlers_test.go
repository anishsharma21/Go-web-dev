package handlers_test

import (
	"net/http"
	"testing"
)

func TestAddUser(t *testing.T) {
	req, err := http.NewRequest("POST", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to perform the request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}
}
