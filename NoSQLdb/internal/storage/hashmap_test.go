package storage

import (
	"fmt"
	"sync"
	"testing"
)

func TestHashMapBasicOperations(t *testing.T) {
	hm := NewHashMap()

	// Test Put and Get
	hm.Put("key1", "value1")
	hm.Put("key2", 42)
	hm.Put("key3", map[string]int{"nested": 1})

	val, ok := hm.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("expected value1, got %v", val)
	}

	val, ok = hm.Get("key2")
	if !ok || val != 42 {
		t.Errorf("expected 42, got %v", val)
	}

	val, ok = hm.Get("nonexistent")
	if ok {
		t.Error("expected key not found")
	}
}

func TestHashMapUpdate(t *testing.T) {
	hm := NewHashMap()

	hm.Put("key", "initial")
	hm.Put("key", "updated")

	val, ok := hm.Get("key")
	if !ok || val != "updated" {
		t.Errorf("expected updated, got %v", val)
	}

	if hm.Size != 1 {
		t.Errorf("expected size 1, got %d", hm.Size)
	}
}

func TestHashMapRemove(t *testing.T) {
	hm := NewHashMap()

	hm.Put("key1", "value1")
	hm.Put("key2", "value2")

	removed := hm.Remove("key1")
	if !removed {
		t.Error("expected key1 to be removed")
	}

	_, ok := hm.Get("key1")
	if ok {
		t.Error("key1 should not exist after removal")
	}

	if hm.Size != 1 {
		t.Errorf("expected size 1, got %d", hm.Size)
	}

	removed = hm.Remove("nonexistent")
	if removed {
		t.Error("should not remove nonexistent key")
	}
}

func TestHashMapResize(t *testing.T) {
	hm := NewHashMap()
	initialCapacity := hm.Capacity

	// Insert enough elements to trigger resize
	for i := 0; i < 20; i++ {
		hm.Put(fmt.Sprintf("key%d", i), i)
	}

	if hm.Capacity <= initialCapacity {
		t.Errorf("expected capacity to increase, got %d", hm.Capacity)
	}

	// Verify all elements are still accessible
	for i := 0; i < 20; i++ {
		val, ok := hm.Get(fmt.Sprintf("key%d", i))
		if !ok || val != i {
			t.Errorf("key%d: expected %d, got %v", i, i, val)
		}
	}
}

func TestHashMapItems(t *testing.T) {
	hm := NewHashMap()

	hm.Put("a", 1)
	hm.Put("b", 2)
	hm.Put("c", 3)

	items := hm.Items()

	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d", len(items))
	}

	if items["a"] != 1 || items["b"] != 2 || items["c"] != 3 {
		t.Errorf("unexpected items: %v", items)
	}
}

func TestHashMapCollision(t *testing.T) {
	hm := NewHashMap()

	// These may or may not collide, but we test chain handling
	keys := []string{"abc", "bca", "cab", "acb", "bac", "cba"}
	for i, key := range keys {
		hm.Put(key, i)
	}

	for i, key := range keys {
		val, ok := hm.Get(key)
		if !ok || val != i {
			t.Errorf("key %s: expected %d, got %v", key, i, val)
		}
	}
}

// Race condition test
func TestHashMapConcurrentAccess(t *testing.T) {
	hm := NewHashMap()
	var wg sync.WaitGroup
	iterations := 100

	// Note: This test will likely fail with -race flag
	// because HashMap is not thread-safe by design.
	// This demonstrates the need for external synchronization.

	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			hm.Put(fmt.Sprintf("writer1_%d", i), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			hm.Put(fmt.Sprintf("writer2_%d", i), i)
		}
	}()

	wg.Wait()

	// Just verify it completes without panic
	// Actual count may vary due to race conditions
	t.Logf("Final size: %d (expected ~%d)", hm.Size, iterations*2)
}

func BenchmarkHashMapPut(b *testing.B) {
	hm := NewHashMap()
	for i := 0; i < b.N; i++ {
		hm.Put(fmt.Sprintf("key%d", i), i)
	}
}

func BenchmarkHashMapGet(b *testing.B) {
	hm := NewHashMap()
	for i := 0; i < 1000; i++ {
		hm.Put(fmt.Sprintf("key%d", i), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hm.Get(fmt.Sprintf("key%d", i%1000))
	}
}
