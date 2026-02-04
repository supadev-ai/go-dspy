package llm

import (
	"context"
	"fmt"
	"strings"
)

// MockClient is a test-friendly implementation of the Client interface.
// It returns predictable responses based on the prompt.
type MockClient struct {
	responses map[string]string
	defaultResponse string
}

// NewMockClient creates a new mock client.
func NewMockClient() *MockClient {
	return &MockClient{
		responses: make(map[string]string),
		defaultResponse: "This is a mock response.",
	}
}

// WithResponse sets a specific response for a given prompt.
func (m *MockClient) WithResponse(prompt, response string) *MockClient {
	m.responses[prompt] = response
	return m
}

// WithDefaultResponse sets the default response when no specific response is configured.
func (m *MockClient) WithDefaultResponse(response string) *MockClient {
	m.defaultResponse = response
	return m
}

// Generate implements the Client interface.
func (m *MockClient) Generate(ctx context.Context, prompt string) (string, error) {
	// Check for exact match
	if response, ok := m.responses[prompt]; ok {
		return response, nil
	}

	// Check for partial matches (contains)
	for key, response := range m.responses {
		if strings.Contains(prompt, key) {
			return response, nil
		}
	}

	// Return default response if set
	if m.defaultResponse != "" {
		return m.defaultResponse, nil
	}

	// No response configured
	return "", fmt.Errorf("no response configured for prompt: %s", prompt)
}

// GenerateWithOptions implements the Client interface.
func (m *MockClient) GenerateWithOptions(ctx context.Context, prompt string, opts *GenerateOptions) (string, error) {
	// Check for exact match
	if response, ok := m.responses[prompt]; ok {
		return response, nil
	}

	// Check for partial matches (contains)
	for key, response := range m.responses {
		if strings.Contains(prompt, key) {
			return response, nil
		}
	}

	// Return default response if set
	if m.defaultResponse != "" {
		return m.defaultResponse, nil
	}

	// No response configured
	return "", fmt.Errorf("no response configured for prompt: %s", prompt)
}
