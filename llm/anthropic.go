package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AnthropicClient implements the Client interface for Anthropic's Claude API.
type AnthropicClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	defaultOpts *GenerateOptions
}

// NewAnthropicClient creates a new Anthropic client with the given API key.
func NewAnthropicClient(apiKey string) *AnthropicClient {
	return &AnthropicClient{
		apiKey:  apiKey,
		baseURL: "https://api.anthropic.com/v1",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		defaultOpts: &GenerateOptions{
			Temperature: 0.7,
			MaxTokens:   1000,
			Model:       "claude-3-sonnet-20240229",
		},
	}
}

// WithBaseURL sets a custom base URL.
func (c *AnthropicClient) WithBaseURL(url string) *AnthropicClient {
	c.baseURL = url
	return c
}

// WithTimeout sets a custom HTTP timeout.
func (c *AnthropicClient) WithTimeout(timeout time.Duration) *AnthropicClient {
	c.httpClient.Timeout = timeout
	return c
}

// WithDefaultOptions sets default generation options.
func (c *AnthropicClient) WithDefaultOptions(opts *GenerateOptions) *AnthropicClient {
	c.defaultOpts = opts
	return c
}

// Generate implements the Client interface.
func (c *AnthropicClient) Generate(ctx context.Context, prompt string) (string, error) {
	return c.GenerateWithOptions(ctx, prompt, c.defaultOpts)
}

// GenerateWithOptions implements the Client interface with custom options.
func (c *AnthropicClient) GenerateWithOptions(ctx context.Context, prompt string, opts *GenerateOptions) (string, error) {
	if c.apiKey == "" {
		return "", &ErrClientNotConfigured{Provider: "Anthropic"}
	}

	model := opts.Model
	if model == "" {
		model = c.defaultOpts.Model
	}

	reqBody := map[string]interface{}{
		"model":     model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": opts.Temperature,
		"max_tokens":  opts.MaxTokens,
	}

	if len(opts.Stop) > 0 {
		reqBody["stop_sequences"] = opts.Stop
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/messages", bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(response.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return response.Content[0].Text, nil
}
