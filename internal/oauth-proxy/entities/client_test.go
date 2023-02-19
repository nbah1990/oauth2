package entities

import (
	"golang.org/x/exp/slices"
	"testing"
)

func TestHasScope(t *testing.T) {
	client := Client{Scopes: []string{"read", "write", "delete"}}

	// Test case 1: client has the requested scope
	if !client.HasScope("write") {
		t.Errorf("Expected HasScope to return true for the scope present in the client")
	}

	// Test case 2: client doesn't have the requested scope
	if client.HasScope("admin") {
		t.Errorf("Expected HasScope to return false for the scope not present in the client")
	}
}

func TestCheckSecret(t *testing.T) {
	client := Client{Secret: "secret-123"}

	// Test case 1: correct secret
	if !client.CheckSecret("secret-123") {
		t.Errorf("Expected CheckSecret to return true for correct secret")
	}

	// Test case 2: incorrect secret
	if client.CheckSecret("incorrect-secret") {
		t.Errorf("Expected CheckSecret to return false for incorrect secret")
	}
}

func TestCreateClient(t *testing.T) {
	// Test case: client is created with given name and scopes
	client := CreateClient("test-client", []string{"read", "write"})
	if client.Name != "test-client" || !slices.Equal(client.Scopes, []string{"read", "write"}) {
		t.Errorf("Expected CreateClient to create client with given name and scopes")
	}
}
