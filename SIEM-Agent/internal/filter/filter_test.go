package filter

import (
	"testing"
	"time"

	"github.com/Narotan/SIEM-Agent/internal/domain"
)

func TestEventFilterSeverityThreshold(t *testing.T) {
	cfg := Config{
		SeverityThreshold: "medium",
	}

	f, err := NewFilter(cfg)
	if err != nil {
		t.Fatalf("failed to create filter: %v", err)
	}

	tests := []struct {
		name     string
		severity string
		expected bool
	}{
		{"low severity blocked", "low", false},
		{"info severity blocked", "info", false},
		{"medium severity passed", "medium", true},
		{"high severity passed", "high", true},
		{"critical severity passed", "critical", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := domain.Event{
				Timestamp: time.Now(),
				Severity:  tt.severity,
				Source:    "test",
			}

			result := f.Match(event)
			if result != tt.expected {
				t.Errorf("Match() = %v, expected %v for severity %s", result, tt.expected, tt.severity)
			}
		})
	}
}

func TestEventFilterExcludeSource(t *testing.T) {
	cfg := Config{
		ExcludeSources: []string{"noisy_source", "spam"},
	}

	f, err := NewFilter(cfg)
	if err != nil {
		t.Fatalf("failed to create filter: %v", err)
	}

	tests := []struct {
		source   string
		expected bool
	}{
		{"noisy_source", false},
		{"spam", false},
		{"important", true},
		{"auth", true},
	}

	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			event := domain.Event{
				Timestamp: time.Now(),
				Source:    tt.source,
			}

			result := f.Match(event)
			if result != tt.expected {
				t.Errorf("Match() = %v, expected %v for source %s", result, tt.expected, tt.source)
			}
		})
	}
}

func TestEventFilterExcludePatterns(t *testing.T) {
	cfg := Config{
		ExcludePatterns: []string{
			"CRON.*session",
			"systemd.*Started",
		},
	}

	f, err := NewFilter(cfg)
	if err != nil {
		t.Fatalf("failed to create filter: %v", err)
	}

	tests := []struct {
		name     string
		rawLog   string
		expected bool
	}{
		{"cron session excluded", "CRON[123]: pam_unix(cron:session): session opened", false},
		{"systemd started excluded", "systemd[1]: Started Daily apt download", false},
		{"auth failure passed", "sshd[456]: Failed password for root", true},
		{"login passed", "sshd[789]: Accepted publickey for user", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := domain.Event{
				Timestamp: time.Now(),
				RawLog:    tt.rawLog,
				Source:    "test",
			}

			result := f.Match(event)
			if result != tt.expected {
				t.Errorf("Match() = %v, expected %v for log: %s", result, tt.expected, tt.rawLog)
			}
		})
	}
}

func TestEventFilterIncludePatterns(t *testing.T) {
	cfg := Config{
		IncludePatterns: []string{
			"ssh",
			"sudo",
			"authentication",
		},
	}

	f, err := NewFilter(cfg)
	if err != nil {
		t.Fatalf("failed to create filter: %v", err)
	}

	tests := []struct {
		name     string
		rawLog   string
		expected bool
	}{
		{"ssh included", "sshd[123]: connection from 192.168.1.1", true},
		{"sudo included", "sudo: user : command=/bin/ls", true},
		{"auth included", "pam: authentication failure", true},
		{"unrelated excluded", "kernel: CPU temperature 45C", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := domain.Event{
				Timestamp: time.Now(),
				RawLog:    tt.rawLog,
				Source:    "test",
			}

			result := f.Match(event)
			if result != tt.expected {
				t.Errorf("Match() = %v, expected %v for log: %s", result, tt.expected, tt.rawLog)
			}
		})
	}
}

func TestEventFilterCombined(t *testing.T) {
	cfg := Config{
		SeverityThreshold: "medium",
		ExcludeSources:    []string{"cron"},
		ExcludePatterns:   []string{"session (opened|closed)"},
		IncludePatterns:   []string{"ssh|sudo|auth"},
	}

	f, err := NewFilter(cfg)
	if err != nil {
		t.Fatalf("failed to create filter: %v", err)
	}

	event := domain.Event{
		Timestamp: time.Now(),
		Severity:  "high",
		Source:    "auth",
		RawLog:    "sshd: Failed password for root",
		EventType: "auth_failure",
	}

	if !f.Match(event) {
		t.Error("expected high severity auth failure to pass combined filter")
	}

	// Low severity should be blocked
	event.Severity = "low"
	if f.Match(event) {
		t.Error("expected low severity to be blocked")
	}
}

func TestEventFilterInvalidPattern(t *testing.T) {
	cfg := Config{
		ExcludePatterns: []string{"[invalid regex"},
	}

	_, err := NewFilter(cfg)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestNoOpFilter(t *testing.T) {
	f := NewNoOpFilter()

	event := domain.Event{
		Timestamp: time.Now(),
		Severity:  "low",
		Source:    "any",
		RawLog:    "anything",
	}

	if !f.Match(event) {
		t.Error("NoOpFilter should pass all events")
	}
}

func TestSeverityToInt(t *testing.T) {
	tests := []struct {
		severity string
		expected int
	}{
		{"low", 1},
		{"info", 1},
		{"medium", 2},
		{"warning", 2},
		{"high", 3},
		{"error", 3},
		{"critical", 3},
		{"unknown", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.severity, func(t *testing.T) {
			result := severityToInt(tt.severity)
			if result != tt.expected {
				t.Errorf("severityToInt(%s) = %d, expected %d", tt.severity, result, tt.expected)
			}
		})
	}
}

func BenchmarkEventFilterMatch(b *testing.B) {
	cfg := Config{
		SeverityThreshold: "low",
		ExcludePatterns:   []string{"CRON", "systemd.*Started"},
		IncludePatterns:   []string{"ssh|sudo|auth"},
	}

	f, _ := NewFilter(cfg)

	event := domain.Event{
		Timestamp: time.Now(),
		Severity:  "high",
		Source:    "auth",
		RawLog:    "sshd: Failed password for root from 192.168.1.100",
		EventType: "auth_failure",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Match(event)
	}
}
