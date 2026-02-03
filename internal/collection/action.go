package collection

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jpalaniselvam/myna/internal/types"
	"github.com/pelletier/go-toml/v2"
)

// CreateActionInput represents the input expected from the UI
type CreateActionInput struct {
	CollectionPath string          `json:"collection_path"`
	SubPath        string          `json:"sub_path"`
	FileName       string          `json:"file_name"`
	Data           json.RawMessage `json:"data"` // The action definition
}

// CreateAction creates a new action file
func CreateAction(input CreateActionInput) error {
	if input.CollectionPath == "" {
		return fmt.Errorf("collection path is required")
	}
	if input.FileName == "" {
		return fmt.Errorf("file name is required")
	}
	if !strings.HasSuffix(input.FileName, ".toml") {
		input.FileName += ".toml"
	}

	// 1. Unmarshal into BaseAction to find type
	var baseAction types.BaseAction
	if err := json.Unmarshal(input.Data, &baseAction); err != nil {
		return fmt.Errorf("failed to parse action data: %w", err)
	}

	if baseAction.Kind == "" {
		return fmt.Errorf("action kind is required")
	}

	// 2. Based on kind, unmarshal into specific struct
	var action interface{}

	switch baseAction.Kind {
	case types.KindLambdaInvoke:
		action = &types.LambdaAction{}
	case types.KindSQSSend, types.KindSQSReceive, types.KindSQSDelete, types.KindSQSPurge, types.KindSQSMoveTask:
		action = &types.SQSAction{}
	case types.KindSNSPublish:
		action = &types.SNSAction{}
	case types.KindEventBridgePutEvents:
		action = &types.EventBridgeAction{}
	case types.KindEC2DescribeInstances, types.KindEC2StartInstances, types.KindEC2StopInstances, types.KindEC2RebootInstances, types.KindEC2TerminateInstances:
		action = &types.EC2Action{}
	case types.KindRDSDescribeDBInstances, types.KindRDSStartDBInstance, types.KindRDSStopDBInstance, types.KindRDSRebootDBInstance:
		action = &types.RDSAction{}
	case types.KindSFNListStateMachines, types.KindSFNStartExecution, types.KindSFNDescribeExecution, types.KindSFNStopExecution:
		action = &types.SFNAction{}
	case types.KindSESSendEmail, types.KindSESVerifyEmailIdentity:
		action = &types.SESAction{}
	default:
		return fmt.Errorf("unsupported action kind: %s", baseAction.Kind)
	}

	if err := json.Unmarshal(input.Data, action); err != nil {
		return fmt.Errorf("failed to parse specific action structure: %w", err)
	}

	// 3. Construct full path
	actionDir := filepath.Join(input.CollectionPath, input.SubPath)
	if err := os.MkdirAll(actionDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fullPath := filepath.Join(actionDir, input.FileName)
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		return fmt.Errorf("action file already exists at %s", fullPath)
	}

	// 4. Marshal to TOML and write
	data, err := toml.Marshal(action)
	if err != nil {
		return fmt.Errorf("failed to marshal action to toml: %w", err)
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write action file: %w", err)
	}

	return nil
}
