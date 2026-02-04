package main

import (
	"context"
	"fmt"
	"os"

	"github.com/supadev-ai/go-dspy/dspy"
	"github.com/supadev-ai/go-dspy/llm"
)

// ReasoningInput represents input for a reasoning task.
type ReasoningInput struct {
	Problem string
}

// ReasoningOutput represents output with chain-of-thought reasoning.
type ReasoningOutput struct {
	Reasoning string
	Answer    string
}

func main() {
	// Get API key from environment or use mock
	apiKey := os.Getenv("OPENAI_API_KEY")
	var client llm.Client

	if apiKey == "" {
		fmt.Println("No OPENAI_API_KEY found, using mock client")
		client = llm.NewMockClient().
			WithResponse("Solve: 2+2", "Reasoning: I need to add 2 and 2 together. 2 + 2 = 4.\nAnswer: 4").
			WithDefaultResponse("Reasoning: [mock reasoning]\nAnswer: [mock answer]")
	} else {
		client = llm.NewOpenAIClient(apiKey)
	}

	// Create a chain-of-thought predictor
	sig := dspy.NewSignature[ReasoningInput, ReasoningOutput](
		"ChainOfThought",
		"Solve problems step by step. Show your reasoning process, then provide the final answer.",
	)

	predictor := dspy.NewPredictor(sig, client)

	// Use the predictor
	ctx := context.Background()
	input := ReasoningInput{
		Problem: "If a train travels 60 miles per hour, how long will it take to travel 120 miles?",
	}

	output, err := predictor.Forward(ctx, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Problem: %s\n", input.Problem)
	fmt.Printf("\nReasoning:\n%s\n", output.Reasoning)
	fmt.Printf("\nAnswer: %s\n", output.Answer)
}
