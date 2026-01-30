package index

import (
	"fmt"
	"sync"
	"testing"
)

func TestBTreeInsertAndSearch(t *testing.T) {
	tree := NewBPlusTree(3)

	tree.Insert(Key("apple"), Value("fruit1"))
	tree.Insert(Key("banana"), Value("fruit2"))
	tree.Insert(Key("cherry"), Value("fruit3"))

	values := tree.Search(Key("banana"))
	if len(values) != 1 || string(values[0]) != "fruit2" {
		t.Errorf("expected [fruit2], got %v", values)
	}

	values = tree.Search(Key("nonexistent"))
	if len(values) != 0 {
		t.Errorf("expected empty result, got %v", values)
	}
}

func TestBTreeDuplicateKeys(t *testing.T) {
	tree := NewBPlusTree(3)

	tree.Insert(Key("key"), Value("value1"))
	tree.Insert(Key("key"), Value("value2"))
	tree.Insert(Key("key"), Value("value3"))

	values := tree.Search(Key("key"))
	if len(values) != 3 {
		t.Errorf("expected 3 values, got %d", len(values))
	}
}

func TestBTreeManyInserts(t *testing.T) {
	tree := NewBPlusTree(3)

	for i := 0; i < 100; i++ {
		key := Key(fmt.Sprintf("key%03d", i))
		value := Value(fmt.Sprintf("value%d", i))
		tree.Insert(key, value)
	}

	// Verify all keys are searchable
	for i := 0; i < 100; i++ {
		key := Key(fmt.Sprintf("key%03d", i))
		values := tree.Search(key)
		if len(values) != 1 {
			t.Errorf("key%03d: expected 1 value, got %d", i, len(values))
		}
	}
}

func TestBTreeSplitLeaf(t *testing.T) {
	tree := NewBPlusTree(2) // Small order to force splits

	// Insert enough to trigger leaf split
	for i := 0; i < 10; i++ {
		tree.Insert(Key(fmt.Sprintf("k%d", i)), Value(fmt.Sprintf("v%d", i)))
	}

	// Verify all values still accessible
	for i := 0; i < 10; i++ {
		values := tree.Search(Key(fmt.Sprintf("k%d", i)))
		if len(values) == 0 {
			t.Errorf("k%d not found after split", i)
		}
	}
}

func TestBTreeOrdering(t *testing.T) {
	tree := NewBPlusTree(3)

	// Insert in random order
	keys := []string{"delta", "alpha", "charlie", "bravo", "echo"}
	for _, k := range keys {
		tree.Insert(Key(k), Value(k+"_val"))
	}

	// All should be findable
	for _, k := range keys {
		values := tree.Search(Key(k))
		if len(values) == 0 {
			t.Errorf("%s not found", k)
		}
	}
}

// Race condition test - B+Tree without synchronization
func TestBTreeConcurrentInsert(t *testing.T) {
	tree := NewBPlusTree(3)
	var wg sync.WaitGroup
	iterations := 50

	// Note: This test demonstrates that B+Tree is not thread-safe
	// It may fail with -race flag

	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			tree.Insert(Key(fmt.Sprintf("a%03d", i)), Value(fmt.Sprintf("val_a%d", i)))
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			tree.Insert(Key(fmt.Sprintf("b%03d", i)), Value(fmt.Sprintf("val_b%d", i)))
		}
	}()

	wg.Wait()

	t.Log("Concurrent insert completed (may have race conditions)")
}

func BenchmarkBTreeInsert(b *testing.B) {
	tree := NewBPlusTree(32)
	for i := 0; i < b.N; i++ {
		tree.Insert(Key(fmt.Sprintf("key%d", i)), Value(fmt.Sprintf("val%d", i)))
	}
}

func BenchmarkBTreeSearch(b *testing.B) {
	tree := NewBPlusTree(32)
	for i := 0; i < 10000; i++ {
		tree.Insert(Key(fmt.Sprintf("key%05d", i)), Value(fmt.Sprintf("val%d", i)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Search(Key(fmt.Sprintf("key%05d", i%10000)))
	}
}
