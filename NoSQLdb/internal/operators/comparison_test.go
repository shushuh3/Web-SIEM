package operators

import "testing"

func TestCompareEq(t *testing.T) {
	tests := []struct {
		name     string
		a, b     any
		expected bool
	}{
		{"equal strings", "hello", "hello", true},
		{"different strings", "hello", "world", false},
		{"equal ints", 42, 42, true},
		{"different ints", 42, 43, false},
		{"equal floats", 3.14, 3.14, true},
		{"int vs float", 42, 42.0, false}, // different types
		{"nil vs nil", nil, nil, true},
		{"nil vs value", nil, "value", false},
		{"equal slices", []int{1, 2, 3}, []int{1, 2, 3}, true},
		{"different slices", []int{1, 2, 3}, []int{1, 2, 4}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareEq(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("CompareEq(%v, %v) = %v, expected %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCompareGt(t *testing.T) {
	tests := []struct {
		name     string
		a, b     any
		expected bool
	}{
		{"int greater", 10, 5, true},
		{"int equal", 5, 5, false},
		{"int less", 5, 10, false},
		{"float greater", 3.14, 2.71, true},
		{"float less", 2.71, 3.14, false},
		{"int vs float greater", 10, 5.5, true},
		{"string comparison", "hello", "world", false}, // non-numeric
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareGt(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("CompareGt(%v, %v) = %v, expected %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCompareLt(t *testing.T) {
	tests := []struct {
		name     string
		a, b     any
		expected bool
	}{
		{"int less", 5, 10, true},
		{"int equal", 5, 5, false},
		{"int greater", 10, 5, false},
		{"float less", 2.71, 3.14, true},
		{"negative numbers", -10, -5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareLt(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("CompareLt(%v, %v) = %v, expected %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCompareIn(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		values   any
		expected bool
	}{
		{"string in slice", "apple", []any{"apple", "banana", "cherry"}, true},
		{"string not in slice", "grape", []any{"apple", "banana", "cherry"}, false},
		{"int in slice", 42, []any{1, 42, 100}, true},
		{"int not in slice", 50, []any{1, 42, 100}, false},
		{"empty slice", "anything", []any{}, false},
		{"invalid values type", "test", "not a slice", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareIn(tt.value, tt.values)
			if result != tt.expected {
				t.Errorf("CompareIn(%v, %v) = %v, expected %v", tt.value, tt.values, result, tt.expected)
			}
		})
	}
}

func TestCompareLike(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		pattern  any
		expected bool
	}{
		{"exact match", "hello", "hello", true},
		{"wildcard end", "hello world", "hello%", true},
		{"wildcard start", "hello world", "%world", true},
		{"wildcard both", "hello world", "%lo wo%", true},
		{"single char wildcard", "hello", "h_llo", true},
		{"no match", "hello", "world", false},
		{"non-string value", 123, "123", false},
		{"non-string pattern", "hello", 123, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareLike(tt.value, tt.pattern)
			if result != tt.expected {
				t.Errorf("CompareLike(%v, %v) = %v, expected %v", tt.value, tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expected  float64
		shouldErr bool
	}{
		{"float64", 3.14, 3.14, false},
		{"float32", float32(3.14), 3.14, false},
		{"int", 42, 42.0, false},
		{"int32", int32(42), 42.0, false},
		{"int64", int64(42), 42.0, false},
		{"uint", uint(42), 42.0, false},
		{"string", "42", 0, true},
		{"nil", nil, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := toFloat64(tt.input)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("expected error for input %v", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				// Allow small floating point differences for float32
				if result != tt.expected && (result-tt.expected > 0.01 || tt.expected-result > 0.01) {
					t.Errorf("toFloat64(%v) = %v, expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func BenchmarkCompareEq(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CompareEq("hello world", "hello world")
	}
}

func BenchmarkCompareLike(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CompareLike("hello world this is a test string", "%test%")
	}
}
