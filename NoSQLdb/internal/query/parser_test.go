package query

import (
	"testing"
)

func TestParseEmptyQuery(t *testing.T) {
	q, err := Parse("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if q.Conditions == nil {
		t.Error("expected non-nil conditions")
	}

	if len(q.Conditions) != 0 {
		t.Errorf("expected empty conditions, got %v", q.Conditions)
	}
}

func TestParseSimpleQuery(t *testing.T) {
	q, err := Parse(`{"name": "Alice", "age": 25}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if q.Conditions["name"] != "Alice" {
		t.Errorf("expected name=Alice, got %v", q.Conditions["name"])
	}

	// JSON numbers are float64
	if q.Conditions["age"] != 25.0 {
		t.Errorf("expected age=25, got %v", q.Conditions["age"])
	}
}

func TestParseQueryWithOperators(t *testing.T) {
	q, err := Parse(`{"age": {"$gt": 18}, "status": {"$in": ["active", "pending"]}}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ageCondition, ok := q.Conditions["age"].(map[string]any)
	if !ok {
		t.Fatalf("expected age to be map, got %T", q.Conditions["age"])
	}

	if ageCondition["$gt"] != 18.0 {
		t.Errorf("expected $gt=18, got %v", ageCondition["$gt"])
	}
}

func TestParseInvalidJSON(t *testing.T) {
	_, err := Parse(`{invalid json}`)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParseDocument(t *testing.T) {
	doc, err := ParseDocument(`{"_id": "123", "name": "Test", "values": [1, 2, 3]}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if doc["_id"] != "123" {
		t.Errorf("expected _id=123, got %v", doc["_id"])
	}

	if doc["name"] != "Test" {
		t.Errorf("expected name=Test, got %v", doc["name"])
	}

	values, ok := doc["values"].([]any)
	if !ok {
		t.Fatalf("expected values to be slice, got %T", doc["values"])
	}

	if len(values) != 3 {
		t.Errorf("expected 3 values, got %d", len(values))
	}
}

func TestParseDocumentInvalidJSON(t *testing.T) {
	_, err := ParseDocument(`not json`)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParseNestedQuery(t *testing.T) {
	q, err := Parse(`{"$or": [{"status": "active"}, {"priority": {"$gt": 5}}]}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	orConditions, ok := q.Conditions["$or"].([]any)
	if !ok {
		t.Fatalf("expected $or to be slice, got %T", q.Conditions["$or"])
	}

	if len(orConditions) != 2 {
		t.Errorf("expected 2 conditions in $or, got %d", len(orConditions))
	}
}

func BenchmarkParse(b *testing.B) {
	query := `{"name": "test", "age": {"$gt": 18}, "status": {"$in": ["a", "b", "c"]}}`
	for i := 0; i < b.N; i++ {
		Parse(query)
	}
}

func BenchmarkParseDocument(b *testing.B) {
	doc := `{"_id": "123", "name": "Test User", "email": "test@example.com", "age": 25, "tags": ["a", "b", "c"]}`
	for i := 0; i < b.N; i++ {
		ParseDocument(doc)
	}
}
