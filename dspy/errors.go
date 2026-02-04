package dspy

import "fmt"

// Error types for DSPy operations
type Error struct {
	Op  string
	Err error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

// Common error constructors
func ErrModuleExecution(op string, err error) *Error {
	return &Error{Op: op, Err: err}
}

func ErrInvalidSignature(op string, err error) *Error {
	return &Error{Op: op, Err: err}
}

func ErrOptimizationFailed(op string, err error) *Error {
	return &Error{Op: op, Err: err}
}
