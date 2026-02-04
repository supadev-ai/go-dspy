package dspy

// Example represents a training example for optimization.
// It contains an input-output pair used for evaluation and optimization.
type Example[I any, O any] struct {
	Input  I
	Output O
}

// NewExample creates a new example with the given input and output.
func NewExample[I any, O any](input I, output O) Example[I, O] {
	return Example[I, O]{
		Input:  input,
		Output: output,
	}
}
