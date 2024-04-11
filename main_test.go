package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := []byte("&8fdsfh dsfkjskdfj hsa -+=/{}[`]")
	data := []byte("GqnfP58Gby8zxk84ffSKjeav0rmLDiF8ezM02KQybripyLFtVuSpQecBUACH6AAb6qDn1Gk72sHiLIItS5qVDpdwl4T57So2i6EhdPW9JgsfDJYNmNDagqyNOwTLU81W")
	ciphertext, err := Encrypt(key, data)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	plaintext, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}
	if string(plaintext) != string(data) {
		t.Fatalf("Decrypt failed: %v", err)
	}
}

// Mock and test httpRequest function
func TestHttpRequest(t *testing.T) {
	// Create a mock HTTP server
	body := "blarblarblarimthebody\nbutbewareimmultiline\n"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle the request and send a response
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, body)
	}))
	defer server.Close()

	// Make a request to the mock server using the httpRequest function
	response, err := httpRequest("GET", server.URL, "X-blar", "secretsecret")
	if err != nil {
		t.Fatalf("httpRequest failed: %v", err)
	}

	// Check if the response matches the expected result
	if response != body {
		t.Fatalf("httpRequest returned unexpected response: got %q, want %q", response, body)
	}
}
