package service

import (
	"errors"
	"sync"
	"testing"
	"time"
)

type fakeRepo struct {
	data []map[string]any
	err  error
	mu   sync.Mutex
}

func (f *fakeRepo) FindAll(_ string, _ map[string]any) ([]map[string]any, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.data, nil
}

func TestGetEventsPaginationAndSort(t *testing.T) {
	repo := &fakeRepo{data: []map[string]any{
		{"_id": "1", "timestamp": "2024-01-02T10:00:00Z"},
		{"_id": "2", "timestamp": "2024-01-03T10:00:00Z"},
		{"_id": "1", "timestamp": "2024-01-03T10:00:00Z"},
	}}

	svc := NewSiemService(repo, "siem_events")

	page, err := svc.GetEvents(1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if page.Total != 3 || page.Count != 2 || page.Page != 1 || page.Limit != 2 || page.TotalPages != 2 {
		t.Fatalf("unexpected paging values: %+v", page)
	}

	firstID, _ := page.Data[0]["_id"].(string)
	secondID, _ := page.Data[1]["_id"].(string)
	if firstID != "2" || secondID != "1" {
		t.Fatalf("unexpected sort order: %v", []string{firstID, secondID})
	}
}

func TestGetStatsCounts(t *testing.T) {
	recent := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	older := time.Now().Add(-48 * time.Hour).UTC().Format(time.RFC3339)

	repo := &fakeRepo{data: []map[string]any{
		{"timestamp": recent, "agent_id": "a1", "event_type": "user_login", "severity": "low", "user": "alice", "process": "ssh"},
		{"timestamp": recent, "agent_id": "a2", "event_type": "auth_failure", "severity": "high", "user": "bob", "process": "sudo"},
		{"timestamp": older, "agent_id": "a1", "event_type": "file_access", "severity": "medium", "user": "alice", "process": "cat"},
	}}

	svc := NewSiemService(repo, "siem_events")

	stats, err := svc.GetStats()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(stats.ActiveAgents) != 2 {
		t.Fatalf("expected 2 active agents, got %d", len(stats.ActiveAgents))
	}

	if stats.EventsByType["user_login"] != 1 || stats.EventsByType["auth_failure"] != 1 {
		t.Fatalf("unexpected EventsByType: %+v", stats.EventsByType)
	}

	if stats.SeverityDist["low"] != 1 || stats.SeverityDist["high"] != 1 {
		t.Fatalf("unexpected SeverityDist: %+v", stats.SeverityDist)
	}

	if stats.TopUsers["alice"] != 1 || stats.TopUsers["bob"] != 1 {
		t.Fatalf("unexpected TopUsers: %+v", stats.TopUsers)
	}

	if stats.TopProcesses["ssh"] != 1 || stats.TopProcesses["sudo"] != 1 {
		t.Fatalf("unexpected TopProcesses: %+v", stats.TopProcesses)
	}

	if len(stats.LastLogins) != 2 {
		t.Fatalf("expected 2 last logins, got %d", len(stats.LastLogins))
	}
}

func TestGetEventsError(t *testing.T) {
	repo := &fakeRepo{err: errors.New("database error")}
	svc := NewSiemService(repo, "siem_events")

	_, err := svc.GetEvents(1, 10)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetEventsEmptyResult(t *testing.T) {
	repo := &fakeRepo{data: []map[string]any{}}
	svc := NewSiemService(repo, "siem_events")

	page, err := svc.GetEvents(1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if page.Total != 0 || page.Count != 0 {
		t.Errorf("expected empty page, got total=%d, count=%d", page.Total, page.Count)
	}
}

func TestGetEventsPageOutOfRange(t *testing.T) {
	repo := &fakeRepo{data: []map[string]any{
		{"_id": "1"},
		{"_id": "2"},
	}}
	svc := NewSiemService(repo, "siem_events")

	page, err := svc.GetEvents(100, 10) // Page way out of range
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if page.Count != 0 {
		t.Errorf("expected empty page for out of range, got count=%d", page.Count)
	}
	if page.Total != 2 {
		t.Errorf("expected total=2, got %d", page.Total)
	}
}

func TestGetEventsLimitClamping(t *testing.T) {
	repo := &fakeRepo{data: []map[string]any{}}
	svc := NewSiemService(repo, "siem_events")

	// Test that invalid page/limit values are handled
	page, _ := svc.GetEvents(0, 0) // Should default to page=1, limit=50
	if page.Page != 1 || page.Limit != 50 {
		t.Errorf("expected defaults page=1, limit=50, got page=%d, limit=%d", page.Page, page.Limit)
	}

	page, _ = svc.GetEvents(1, 300) // Should clamp limit to 200
	if page.Limit != 200 {
		t.Errorf("expected limit clamped to 200, got %d", page.Limit)
	}
}

func TestGetStatsError(t *testing.T) {
	repo := &fakeRepo{err: errors.New("database error")}
	svc := NewSiemService(repo, "siem_events")

	_, err := svc.GetStats()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetStatsCache(t *testing.T) {
	repo := &fakeRepo{data: []map[string]any{
		{"timestamp": time.Now().UTC().Format(time.RFC3339), "agent_id": "a1", "event_type": "test", "severity": "low"},
	}}

	svc := NewSiemService(repo, "siem_events").(*siemService)

	// First call
	_, _ = svc.GetStats()

	// Second call should use cache
	_, _ = svc.GetStats()

	// Cache should be used
	if svc.statsCache == nil {
		t.Error("expected stats to be cached")
	}
}

func TestExportAllEventsError(t *testing.T) {
	repo := &fakeRepo{err: errors.New("export error")}
	svc := NewSiemService(repo, "siem_events")

	_, err := svc.ExportAllEvents()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestExportAllEventsSorted(t *testing.T) {
	repo := &fakeRepo{data: []map[string]any{
		{"_id": "1", "timestamp": "2024-01-01T10:00:00Z"},
		{"_id": "2", "timestamp": "2024-01-03T10:00:00Z"},
		{"_id": "3", "timestamp": "2024-01-02T10:00:00Z"},
	}}
	svc := NewSiemService(repo, "siem_events")

	events, err := svc.ExportAllEvents()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}

	// Should be sorted by timestamp descending
	ts1 := events[0]["timestamp"].(string)
	ts2 := events[1]["timestamp"].(string)
	ts3 := events[2]["timestamp"].(string)

	if ts1 < ts2 || ts2 < ts3 {
		t.Errorf("events not sorted descending: %s, %s, %s", ts1, ts2, ts3)
	}
}

// Race condition test - concurrent access to service
func TestGetStatsConcurrent(t *testing.T) {
	recent := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	repo := &fakeRepo{data: []map[string]any{
		{"timestamp": recent, "agent_id": "a1", "event_type": "user_login", "severity": "low", "user": "alice", "process": "ssh"},
	}}

	svc := NewSiemService(repo, "siem_events")

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				_, err := svc.GetStats()
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		}()
	}
	wg.Wait()
}
