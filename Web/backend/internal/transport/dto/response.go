package dto

import "time"

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type EventsResponse struct {
	Status     string           `json:"status"`
	Count      int              `json:"count"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"totalPages"`
	Data       []map[string]any `json:"data"`
}

type StatsResponse struct {
	ActiveAgents  map[string]time.Time `json:"active_agents"`
	EventsByType  map[string]int       `json:"events_by_type"`
	SeverityDist  map[string]int       `json:"severity_distribution"`
	TopUsers      map[string]int       `json:"top_users"`
	TopProcesses  map[string]int       `json:"top_processes"`
	EventsPerHour map[int]int          `json:"events_per_hour"`
	LastLogins    []map[string]any     `json:"last_logins"`
}

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}
