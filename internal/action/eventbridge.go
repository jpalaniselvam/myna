package action

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"

	"github.com/jpalaniselvam/myna/internal/parser"
	baseTypes "github.com/jpalaniselvam/myna/internal/types"
)

// EventBridgeAction represents the structure of an eventbridge.put_events action TOML file
type EventBridgeAction struct {
	baseTypes.BaseAction
	EventBridge EventBridgeConfig `toml:"eventbridge"`
}

// EventBridgeConfig defines the EventBridge specific configuration
type EventBridgeConfig struct {
	BusName     string   `toml:"bus_name"`
	Source      string   `toml:"source"`
	DetailType  string   `toml:"detail_type"`
	Resources   []string `toml:"resources"`
	TraceHeader string   `toml:"trace_header"`
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

func runEventBridgeAction(cfg aws.Config, content []byte, payload []byte) error {
	var action EventBridgeAction
	if err := parser.Unmarshal(content, &action); err != nil {
		return fmt.Errorf("failed to parse eventbridge action: %w", err)
	}

	input, err := buildPutEventsInput(action, payload)
	if err != nil {
		return err
	}

	client := eventbridge.NewFromConfig(cfg)

	ctx := context.TODO()
	if action.TimeoutMs > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(action.TimeoutMs)*time.Millisecond)
		defer cancel()
	}

	resp, err := client.PutEvents(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put events: %w", err)
	}

	output := processEventBridgeOutput(resp)
	fmt.Println(output.String())

	return nil
}

func buildPutEventsInput(action EventBridgeAction, payload []byte) (*eventbridge.PutEventsInput, error) {
	var detail string
	if payload != nil {
		detail = string(payload)
	} else {
		// Valid to send empty detail? Usually JSON object. Let's send empty JSON object if internal requirement.
		// However, API allows string. We'll pass empty string if nil, though likely user error or specific use case.
		detail = "{}"
	}

	var resources []string
	for _, r := range action.EventBridge.Resources {
		resources = append(resources, r)
	}

	entry := types.PutEventsRequestEntry{
		Detail:       aws.String(detail),
		DetailType:   aws.String(action.EventBridge.DetailType),
		Source:       aws.String(action.EventBridge.Source),
		EventBusName: aws.String(action.EventBridge.BusName),
		Resources:    resources,
	}

	if action.EventBridge.TraceHeader != "" {
		entry.TraceHeader = aws.String(action.EventBridge.TraceHeader)
	}

	return &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{entry},
	}, nil
}

func processEventBridgeOutput(resp *eventbridge.PutEventsOutput) *EventBridgeOutput {
	output := &EventBridgeOutput{
		FailedEntryCount: resp.FailedEntryCount,
		Entries:          make([]EventBridgeEntryOutput, len(resp.Entries)),
	}

	for i, entry := range resp.Entries {
		output.Entries[i] = EventBridgeEntryOutput{
			EventId:      entry.EventId,
			ErrorCode:    entry.ErrorCode,
			ErrorMessage: entry.ErrorMessage,
		}
	}
	return output
}
