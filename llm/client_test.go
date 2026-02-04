package llm

import (
	"context"
	"testing"
)

func TestErrClientNotConfigured(t *testing.T) {
	err := &ErrClientNotConfigured{Provider: "TestProvider"}
	expected := "LLM client not configured: TestProvider"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestClientInterface(t *testing.T) {
	// Test that MockClient implements Client interface
	// Create client without default response to test error case
	var client Client = &MockClient{
		responses:       make(map[string]string),
		defaultResponse: "",
	}
	if client == nil {
		t.Fatal("MockClient should implement Client interface")
	}

	ctx := context.Background()
	_, err := client.Generate(ctx, "test")
	// Error is expected since no response is configured and no default is set
	if err == nil {
		t.Error("Expected error for unconfigured client, got nil")
	}

	// Test that it works with default response
	client2 := NewMockClient()
	response, err := client2.Generate(ctx, "test")
	if err != nil {
		t.Errorf("Expected no error with default response, got %v", err)
	}
	if response == "" {
		t.Error("Expected default response, got empty string")
	}
}
