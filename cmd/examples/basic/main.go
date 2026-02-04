package main

import (
	"context"
	"fmt"
	"os"

	"github.com/supadev-ai/go-dspy/dspy"
	"github.com/supadev-ai/go-dspy/llm"
)

// QAInput represents the input to a question-answering module.
type QAInput struct {
	Question string
}

// QAOutput represents the output from a question-answering module.
type QAOutput struct {
	Answer string
}

func main() {
	// Get API key from environment or use mock
	apiKey := os.Getenv("OPENAI_API_KEY")
	var client llm.Client

	if apiKey == "" {
		fmt.Println("No OPENAI_API_KEY found, using mock client")
		client = llm.NewMockClient().
			WithResponse("What is DSPy?", "DSPy is a framework for building LLM applications with automatic prompt optimization.").
			WithDefaultResponse("This is a mock answer.")
	} else {
		client = llm.NewOpenAIClient(apiKey)
	}

	// Create a predictor with a signature
	sig := dspy.NewSignature[QAInput, QAOutput](
		"QuestionAnswering",
		"Answer questions based on your knowledge.",
	)

	predictor := dspy.NewPredictor(sig, client)

	// Use the predictor
	ctx := context.Background()
	input := QAInput{
		Question: "What is DSPy?",
	}

	output, err := predictor.Forward(ctx, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Question: %s\n", input.Question)
	fmt.Printf("Answer: %s\n", output.Answer)
}
