package main

import (
	"context"
	"fmt"
	"os"

	"github.com/supadev-ai/go-dspy/dspy"
	"github.com/supadev-ai/go-dspy/llm"
	"github.com/supadev-ai/go-dspy/optimizer"
)

// ClassificationInput represents input for text classification.
type ClassificationInput struct {
	Text string
}

// ClassificationOutput represents the classification result.
type ClassificationOutput struct {
	Label string
}

func main() {
	// Get API key from environment or use mock
	apiKey := os.Getenv("OPENAI_API_KEY")
	var client llm.Client

	if apiKey == "" {
		fmt.Println("No OPENAI_API_KEY found, using mock client")
		client = llm.NewMockClient().
			WithResponse("I love this product!", "Label: positive").
			WithResponse("This is terrible", "Label: negative").
			WithDefaultResponse("Label: neutral")
	} else {
		client = llm.NewOpenAIClient(apiKey)
	}

	// Create training examples
	examples := []dspy.Example[ClassificationInput, ClassificationOutput]{
		dspy.NewExample(
			ClassificationInput{Text: "I love this product!"},
			ClassificationOutput{Label: "positive"},
		),
		dspy.NewExample(
			ClassificationInput{Text: "This is terrible"},
			ClassificationOutput{Label: "negative"},
		),
		dspy.NewExample(
			ClassificationInput{Text: "It's okay, nothing special"},
			ClassificationOutput{Label: "neutral"},
		),
	}

	// Create initial predictor
	sig := dspy.NewSignature[ClassificationInput, ClassificationOutput](
		"TextClassification",
		"Classify the sentiment of the given text as positive, negative, or neutral.",
	)

	predictor := dspy.NewPredictor(sig, client)

	// Create optimizer
	opt := optimizer.NewBootstrapOptimizer[ClassificationInput, ClassificationOutput]().
		WithMaxIterations(5).
		WithMinScore(0.8)

	// Create metric
	metric := optimizer.ExactMatch[ClassificationInput, ClassificationOutput]()

	// Optimize the predictor
	fmt.Println("Optimizing predictor...")
	ctx := context.Background()
	optimized, err := opt.Optimize(ctx, predictor, examples, metric)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Optimization error: %v\n", err)
		os.Exit(1)
	}

	// Test the optimized predictor
	testInput := ClassificationInput{
		Text: "This is amazing!",
	}

	output, err := optimized.Forward(ctx, testInput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nTest Input: %s\n", testInput.Text)
	fmt.Printf("Predicted Label: %s\n", output.Label)
}
