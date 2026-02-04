package memory

import (
	"context"
	"testing"
)

func TestNewInMemoryStore(t *testing.T) {
	store := NewInMemoryStore()
	if store == nil {
		t.Fatal("Expected store, got nil")
	}
}

func TestInMemoryStore_PutAndGet(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()

	key := "test-key"
	value := "test-value"

	err := store.Put(ctx, key, value)
	if err != nil {
		t.Fatalf("Expected no error on Put, got %v", err)
	}

	retrieved, err := store.Get(ctx, key)
	if err != nil {
		t.Fatalf("Expected no error on Get, got %v", err)
	}

	if retrieved != value {
		t.Errorf("Expected value '%s', got '%v'", value, retrieved)
	}
}

func TestInMemoryStore_Get_NotFound(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()

	_, err := store.Get(ctx, "non-existent-key")
	if err == nil {
		t.Error("Expected error for non-existent key, got nil")
	}
}

func TestInMemoryStore_Delete(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()

	key := "test-key"
	value := "test-value"

	store.Put(ctx, key, value)
	
	err := store.Delete(ctx, key)
	if err != nil {
		t.Fatalf("Expected no error on Delete, got %v", err)
	}

	_, err = store.Get(ctx, key)
	if err == nil {
		t.Error("Expected error after deletion, got nil")
	}
}

func TestInMemoryStore_Clear(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()

	store.Put(ctx, "key1", "value1")
	store.Put(ctx, "key2", "value2")

	err := store.Clear(ctx)
	if err != nil {
		t.Fatalf("Expected no error on Clear, got %v", err)
	}

	_, err = store.Get(ctx, "key1")
	if err == nil {
		t.Error("Expected error after Clear, got nil")
	}

	_, err = store.Get(ctx, "key2")
	if err == nil {
		t.Error("Expected error after Clear, got nil")
	}
}

func TestInMemoryStore_ConcurrentAccess(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()

	// Simple concurrent test
	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 100; i++ {
			store.Put(ctx, "key1", "value1")
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			store.Put(ctx, "key2", "value2")
		}
		done <- true
	}()

	<-done
	<-done

	// Verify both keys exist
	val1, err1 := store.Get(ctx, "key1")
	val2, err2 := store.Get(ctx, "key2")

	if err1 != nil || val1 != "value1" {
		t.Error("Concurrent access failed for key1")
	}

	if err2 != nil || val2 != "value2" {
		t.Error("Concurrent access failed for key2")
	}
}

func TestInMemoryStore_ContextCancellation(t *testing.T) {
	store := NewInMemoryStore()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := store.Put(ctx, "key", "value")
	if err == nil {
		t.Error("Expected error on cancelled context, got nil")
	}

	_, err = store.Get(ctx, "key")
	if err == nil {
		t.Error("Expected error on cancelled context, got nil")
	}
}
