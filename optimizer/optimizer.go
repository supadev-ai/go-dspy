package optimizer

import (
	"context"

	"github.com/supadev-ai/go-dspy/dspy"
)

// Optimizer is an interface for optimizing DSPy modules.
type Optimizer[I any, O any] interface {
	// Optimize improves a module's performance on the given examples.
	Optimize(
		ctx context.Context,
		module dspy.Module[I, O],
		examples []dspy.Example[I, O],
		metric Metric[I, O],
	) (dspy.Module[I, O], error)
}
