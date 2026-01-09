package domain

import "time"

type Event struct {
	Timestamp string
	AgentID   string
	EventType string
	Severity  string
	User      string
	Process   string
	Message   string
	Raw       map[string]any
}

type DashboardStats struct {
	ActiveAgents  map[string]time.Time
	EventsByType  map[string]int
	SeverityDist  map[string]int
	TopUsers      map[string]int
	TopProcesses  map[string]int
	EventsPerHour map[int]int
	LastLogins    []map[string]any
}

type EventsPage struct {
	Data       []map[string]any
	Count      int
	Total      int
	Page       int
	Limit      int
	TotalPages int
}
