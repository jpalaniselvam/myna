package types

import "encoding/json"

// EventBridgeAction represents the structure of an eventbridge.put_events action TOML file
type EventBridgeAction struct {
	BaseAction
	EventBridge EventBridgeConfig `toml:"eventbridge" json:"eventbridge"`
}

// EventBridgeConfig defines the EventBridge specific configuration
type EventBridgeConfig struct {
	BusName     string   `toml:"bus_name" json:"bus_name"`
	Source      string   `toml:"source" json:"source"`
	DetailType  string   `toml:"detail_type" json:"detail_type"`
	Resources   []string `toml:"resources,omitempty" json:"resources,omitempty"`
	TraceHeader string   `toml:"trace_header,omitempty" json:"trace_header,omitempty"`
}

// EventBridgeOutput represents the output of an EventBridge put_events execution
type EventBridgeOutput struct {
	FailedEntryCount int32                    `json:"FailedEntryCount"`
	Entries          []EventBridgeEntryOutput `json:"Entries"`
}

type EventBridgeEntryOutput struct {
	EventId      *string `json:"EventId"`
	ErrorCode    *string `json:"ErrorCode,omitempty"`
	ErrorMessage *string `json:"ErrorMessage,omitempty"`
}

func (e *EventBridgeOutput) String() string {
	b, _ := json.MarshalIndent(e, "", "  ")
	return string(b)
}
