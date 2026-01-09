package model

type DBRequest struct {
	Database string           `json:"database"`
	Command  string           `json:"operation"`
	Data     []map[string]any `json:"data,omitempty"`
	Query    map[string]any   `json:"query,omitempty"`
}

type DBResponse struct {
	Status  string           `json:"status"`
	Message string           `json:"message,omitempty"`
	Data    []map[string]any `json:"data,omitempty"`
	Count   int              `json:"count,omitempty"`
}
