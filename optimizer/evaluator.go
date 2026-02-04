package optimizer

import (
	"context"
	"strings"

	"github.com/supadev-ai/go-dspy/dspy"
)

// Metric is a function that evaluates the quality of a prediction.
// It returns a score where higher is better.
type Metric[I any, O any] func(predicted, expected O) float64

// ExactMatch returns 1.0 if predicted exactly matches expected, 0.0 otherwise.
func ExactMatch[I any, O comparable]() Metric[I, O] {
	return func(predicted, expected O) float64 {
		if predicted == expected {
			return 1.0
		}
		return 0.0
	}
}

// StringContains returns 1.0 if predicted string contains expected, 0.0 otherwise.
func StringContains[I any]() Metric[I, string] {
	return func(predicted, expected string) float64 {
		// Simple contains check
		if len(expected) == 0 {
			return 1.0
		}
		if len(predicted) == 0 {
			return 0.0
		}
		// Check if predicted contains expected (case-insensitive)
		if strings.Contains(strings.ToLower(predicted), strings.ToLower(expected)) {
			return 1.0
		}
		return 0.0
	}
}

// Evaluate runs a module on examples and returns the average metric score.
func Evaluate[I any, O any](
	module dspy.Module[I, O],
	examples []dspy.Example[I, O],
	metric Metric[I, O],
) (float64, error) {
	if len(examples) == 0 {
		return 0.0, nil
	}

	ctx := context.Background()
	var totalScore float64
	for _, ex := range examples {
		predicted, err := module.Forward(ctx, ex.Input)
		if err != nil {
			return 0.0, err
		}
		score := metric(predicted, ex.Output)
		totalScore += score
	}

	return totalScore / float64(len(examples)), nil
}
