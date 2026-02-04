package memory

import "context"

// Store is an interface for storing and retrieving conversation history
// and intermediate results in DSPy pipelines.
type Store interface {
	// Put stores a value with the given key.
	Put(ctx context.Context, key string, value interface{}) error
	
	// Get retrieves a value by key.
	Get(ctx context.Context, key string) (interface{}, error)
	
	// Delete removes a value by key.
	Delete(ctx context.Context, key string) error
	
	// Clear removes all stored values.
	Clear(ctx context.Context) error
}
