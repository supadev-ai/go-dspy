package llm

import (
	"context"
	"fmt"
)

// Client is the interface for LLM providers.
// It abstracts over different LLM APIs (OpenAI, Anthropic, etc.)
type Client interface {
	// Generate sends a prompt to the LLM and returns the generated text.
	Generate(ctx context.Context, prompt string) (string, error)
	
	// GenerateWithOptions allows for more control over generation parameters.
	GenerateWithOptions(ctx context.Context, prompt string, opts *GenerateOptions) (string, error)
}

// GenerateOptions provides configuration for LLM generation.
type GenerateOptions struct {
	Temperature float64
	MaxTokens   int
	Model       string
	Stop        []string
}

// DefaultGenerateOptions returns sensible defaults for generation.
func DefaultGenerateOptions() *GenerateOptions {
	return &GenerateOptions{
		Temperature: 0.7,
		MaxTokens:   1000,
		Model:       "",
		Stop:        nil,
	}
}

// ErrClientNotConfigured is returned when a client is not properly initialized.
type ErrClientNotConfigured struct {
	Provider string
}

func (e *ErrClientNotConfigured) Error() string {
	return fmt.Sprintf("LLM client not configured: %s", e.Provider)
}
