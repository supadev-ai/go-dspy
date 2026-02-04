package optimizer

import (
	"context"
	"reflect"
)

// mockModule is a test helper that implements Module interface
type mockModule[I any, O any] struct {
	responses map[string]O
}

func (m *mockModule[I, O]) Forward(ctx context.Context, input I) (O, error) {
	var zero O
	// Extract key from input using reflection
	key := ""
	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return zero, nil
		}
		val = val.Elem()
	}
	if val.Kind() == reflect.Struct {
		// Try to find Question or Text field
		typ := val.Type()
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if field.Name == "Question" || field.Name == "Text" {
				fieldVal := val.Field(i)
				if fieldVal.Kind() == reflect.String {
					key = fieldVal.String()
					break
				}
			}
		}
	}
	if output, ok := m.responses[key]; ok {
		return output, nil
	}
	return zero, nil
}

// Helper to create a mock module that implements dspy.Module
func newMockModule[I any, O any]() *mockModule[I, O] {
	return &mockModule[I, O]{
		responses: make(map[string]O),
	}
}
