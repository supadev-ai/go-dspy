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

// OpenAIClient implements the Client interface for OpenAI's API.
type OpenAIClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	defaultOpts *GenerateOptions
}

// NewOpenAIClient creates a new OpenAI client with the given API key.
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		apiKey:  apiKey,
		baseURL: "https://api.openai.com/v1",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		defaultOpts: &GenerateOptions{
			Temperature: 0.7,
			MaxTokens:   1000,
			Model:       "gpt-3.5-turbo",
		},
	}
}

// WithBaseURL sets a custom base URL (useful for Azure OpenAI or proxies).
func (c *OpenAIClient) WithBaseURL(url string) *OpenAIClient {
	c.baseURL = url
	return c
}

// WithTimeout sets a custom HTTP timeout.
func (c *OpenAIClient) WithTimeout(timeout time.Duration) *OpenAIClient {
	c.httpClient.Timeout = timeout
	return c
}

// WithDefaultOptions sets default generation options.
func (c *OpenAIClient) WithDefaultOptions(opts *GenerateOptions) *OpenAIClient {
	c.defaultOpts = opts
	return c
}

// Generate implements the Client interface.
func (c *OpenAIClient) Generate(ctx context.Context, prompt string) (string, error) {
	return c.GenerateWithOptions(ctx, prompt, c.defaultOpts)
}

// GenerateWithOptions implements the Client interface with custom options.
func (c *OpenAIClient) GenerateWithOptions(ctx context.Context, prompt string, opts *GenerateOptions) (string, error) {
	if c.apiKey == "" {
		return "", &ErrClientNotConfigured{Provider: "OpenAI"}
	}

	model := opts.Model
	if model == "" {
		model = c.defaultOpts.Model
	}

	reqBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": opts.Temperature,
		"max_tokens":  opts.MaxTokens,
	}

	if len(opts.Stop) > 0 {
		reqBody["stop"] = opts.Stop
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

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
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return response.Choices[0].Message.Content, nil
}
