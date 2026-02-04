package dspy

// Signature defines the input/output contract for a DSPy module.
// It provides type-safe interfaces for LLM pipelines.
type Signature[I any, O any] struct {
	Name        string
	Description string
}

// NewSignature creates a new signature with the given name and description.
func NewSignature[I any, O any](name, description string) Signature[I, O] {
	return Signature[I, O]{
		Name:        name,
		Description: description,
	}
}
