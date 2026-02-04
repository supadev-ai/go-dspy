package dspy

import (
	"context"
	"testing"

	"github.com/supadev-ai/go-dspy/llm"
)

func TestNewPredictor(t *testing.T) {
	type Input struct {
		Question string
	}
	type Output struct {
		Answer string
	}

	sig := NewSignature[Input, Output]("Test", "Test description")
	client := llm.NewMockClient()

	predictor := NewPredictor(sig, client)

	if predictor.Signature.Name != "Test" {
		t.Errorf("Expected signature name 'Test', got '%s'", predictor.Signature.Name)
	}

	if predictor.Client == nil {
		t.Error("Expected client to be set, got nil")
	}
}

func TestPredictor_Forward_StringOutput(t *testing.T) {
	type Input struct {
		Question string
	}
	type Output struct {
		Answer string
	}

	client := llm.NewMockClient().
		WithResponse("Question: What is Go?", "Answer: Go is a programming language.")

	sig := NewSignature[Input, Output]("QA", "Answer questions")
	predictor := NewPredictor(sig, client)

	ctx := context.Background()
	input := Input{Question: "What is Go?"}

	output, err := predictor.Forward(ctx, input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if output.Answer == "" {
		t.Error("Expected answer to be populated, got empty string")
	}
}

func TestPredictor_Forward_StructOutput(t *testing.T) {
	type Input struct {
		Text string
	}
	type Output struct {
		Label string
		Score float64
	}

	client := llm.NewMockClient().
		WithResponse("Text: I love this!", "Label: positive\nScore: 0.9")

	sig := NewSignature[Input, Output]("Classification", "Classify text")
	predictor := NewPredictor(sig, client)

	ctx := context.Background()
	input := Input{Text: "I love this!"}

	output, err := predictor.Forward(ctx, input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if output.Label == "" {
		t.Error("Expected label to be populated, got empty string")
	}
}

func TestPredictor_Forward_JSONOutput(t *testing.T) {
	type Input struct {
		Text string
	}
	type Output struct {
		Label string
	}

	client := llm.NewMockClient().
		WithResponse("Text: test", `{"Label": "positive"}`)

	sig := NewSignature[Input, Output]("Classification", "Classify text")
	predictor := NewPredictor(sig, client)

	ctx := context.Background()
	input := Input{Text: "test"}

	output, err := predictor.Forward(ctx, input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if output.Label == "" {
		t.Error("Expected label to be populated from JSON, got empty string")
	}
}

func TestPredictor_Forward_ClientError(t *testing.T) {
	type Input struct {
		Question string
	}
	type Output struct {
		Answer string
	}

	// Client that will return an error
	client := llm.NewMockClient().
		WithDefaultResponse("") // Empty default will cause error for unconfigured prompts

	sig := NewSignature[Input, Output]("QA", "Answer questions")
	predictor := NewPredictor(sig, client)

	ctx := context.Background()
	input := Input{Question: "Unknown question"}

	_, err := predictor.Forward(ctx, input)
	if err == nil {
		t.Error("Expected error from client, got nil")
	}
}
