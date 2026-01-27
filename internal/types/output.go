package types

import "time"

// ActionOutput represents the standardized output for all actions
type ActionOutput struct {
	Metadata ResponseMetadata `json:"metadata"`
	Data     interface{}      `json:"data,omitempty"`
	Error    *ActionError     `json:"error,omitempty"`
}

// ResponseMetadata contains metadata about the execution
type ResponseMetadata struct {
	Status    string    `json:"status"` // "success" or "failure"
	Action    Kind      `json:"action"`
	Timestamp time.Time `json:"timestamp"`
	Profile   string    `json:"profile,omitempty"`
	Region    string    `json:"region,omitempty"`
}

// ActionError represents a standardized error structure
type ActionError struct {
	Code    string `json:"code"` // e.g., "ResourceNotFound", "AccessDenied"
	Message string `json:"message"`
}
