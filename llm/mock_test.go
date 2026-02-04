package llm

import (
	"context"
	"testing"
)

func TestNewMockClient(t *testing.T) {
	client := NewMockClient()
	if client == nil {
		t.Fatal("Expected mock client, got nil")
	}
}

func TestMockClient_Generate(t *testing.T) {
	client := NewMockClient().
		WithResponse("test prompt", "test response")

	ctx := context.Background()
	response, err := client.Generate(ctx, "test prompt")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response != "test response" {
		t.Errorf("Expected 'test response', got '%s'", response)
	}
}

func TestMockClient_Generate_DefaultResponse(t *testing.T) {
	client := NewMockClient().
		WithDefaultResponse("default response")

	ctx := context.Background()
	response, err := client.Generate(ctx, "unknown prompt")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response != "default response" {
		t.Errorf("Expected 'default response', got '%s'", response)
	}
}

func TestMockClient_Generate_NoResponse(t *testing.T) {
	// Create client without default response
	client := &MockClient{
		responses:       make(map[string]string),
		defaultResponse: "", // No default
	}

	ctx := context.Background()
	_, err := client.Generate(ctx, "unknown prompt")
	if err == nil {
		t.Error("Expected error for unconfigured prompt, got nil")
	}
}

func TestMockClient_Generate_PartialMatch(t *testing.T) {
	client := NewMockClient().
		WithResponse("test", "matched response")

	ctx := context.Background()
	response, err := client.Generate(ctx, "this is a test prompt")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response != "matched response" {
		t.Errorf("Expected 'matched response', got '%s'", response)
	}
}

func TestMockClient_GenerateWithOptions(t *testing.T) {
	client := NewMockClient().
		WithResponse("test", "response")

	ctx := context.Background()
	opts := DefaultGenerateOptions()
	response, err := client.GenerateWithOptions(ctx, "test", opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response != "response" {
		t.Errorf("Expected 'response', got '%s'", response)
	}
}

func TestDefaultGenerateOptions(t *testing.T) {
	opts := DefaultGenerateOptions()
	if opts == nil {
		t.Fatal("Expected options, got nil")
	}

	if opts.Temperature != 0.7 {
		t.Errorf("Expected temperature 0.7, got %f", opts.Temperature)
	}

	if opts.MaxTokens != 1000 {
		t.Errorf("Expected max tokens 1000, got %d", opts.MaxTokens)
	}
}
