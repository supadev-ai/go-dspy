package memory

import (
	"context"
	"fmt"
	"sync"
)

// InMemoryStore is a thread-safe in-memory implementation of Store.
type InMemoryStore struct {
	mu    sync.RWMutex
	store map[string]interface{}
}

// NewInMemoryStore creates a new in-memory store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		store: make(map[string]interface{}),
	}
}

// Put implements the Store interface.
func (s *InMemoryStore) Put(ctx context.Context, key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		s.store[key] = value
		return nil
	}
}

// Get implements the Store interface.
func (s *InMemoryStore) Get(ctx context.Context, key string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		value, ok := s.store[key]
		if !ok {
			return nil, fmt.Errorf("key not found: %s", key)
		}
		return value, nil
	}
}

// Delete implements the Store interface.
func (s *InMemoryStore) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		delete(s.store, key)
		return nil
	}
}

// Clear implements the Store interface.
func (s *InMemoryStore) Clear(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		s.store = make(map[string]interface{})
		return nil
	}
}
