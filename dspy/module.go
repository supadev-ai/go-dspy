package dspy

import "context"

// Module is the core interface for all DSPy components.
// It represents a composable unit that transforms input I to output O.
type Module[I any, O any] interface {
	// Forward executes the module with the given input and returns the output.
	Forward(ctx context.Context, input I) (O, error)
}
