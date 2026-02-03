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

// UpdateActionInput represents the input for updating an action
type UpdateActionInput struct {
	CollectionPath string          `json:"collection_path"`
	SubPath        string          `json:"sub_path"`
	FileName       string          `json:"file_name"`
	Data           json.RawMessage `json:"data"` // The action definition
}

// GetActionInput defines the location of an action file
type GetActionInput struct {
	CollectionPath string `json:"collection_path"`
	SubPath        string `json:"sub_path"`
	FileName       string `json:"file_name"`
}

// Helper to get specific action struct
func getActionStruct(kind types.Kind) (interface{}, error) {
	switch kind {
	case types.KindLambdaInvoke:
		return &types.LambdaAction{}, nil
	case types.KindSQSSend, types.KindSQSReceive, types.KindSQSDelete, types.KindSQSPurge, types.KindSQSMoveTask:
		return &types.SQSAction{}, nil
	case types.KindSNSPublish:
		return &types.SNSAction{}, nil
	case types.KindEventBridgePutEvents:
		return &types.EventBridgeAction{}, nil
	case types.KindEC2DescribeInstances, types.KindEC2StartInstances, types.KindEC2StopInstances, types.KindEC2RebootInstances, types.KindEC2TerminateInstances:
		return &types.EC2Action{}, nil
	case types.KindRDSDescribeDBInstances, types.KindRDSStartDBInstance, types.KindRDSStopDBInstance, types.KindRDSRebootDBInstance:
		return &types.RDSAction{}, nil
	case types.KindSFNListStateMachines, types.KindSFNStartExecution, types.KindSFNDescribeExecution, types.KindSFNStopExecution:
		return &types.SFNAction{}, nil
	case types.KindSESSendEmail, types.KindSESVerifyEmailIdentity:
		return &types.SESAction{}, nil
	default:
		return nil, fmt.Errorf("unsupported action kind: %s", kind)
	}
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

	// 2. Based on kind, get specific struct
	action, err := getActionStruct(baseAction.Kind)
	if err != nil {
		return err
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

// UpdateAction updates an existing action file
func UpdateAction(input UpdateActionInput) error {
	if input.CollectionPath == "" {
		return fmt.Errorf("collection path is required")
	}
	if input.FileName == "" {
		return fmt.Errorf("file name is required")
	}
	// Ensure filename ends with .toml
	if !strings.HasSuffix(input.FileName, ".toml") {
		input.FileName += ".toml"
	}

	fullPath := filepath.Join(input.CollectionPath, input.SubPath, input.FileName)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("action file does not exist at %s", fullPath)
	}

	// Unmarshal JSON to get Kind
	var baseAction types.BaseAction
	if err := json.Unmarshal(input.Data, &baseAction); err != nil {
		return fmt.Errorf("failed to parse action data: %w", err)
	}
	if baseAction.Kind == "" {
		return fmt.Errorf("action kind is required")
	}

	// Get specific struct
	action, err := getActionStruct(baseAction.Kind)
	if err != nil {
		return err
	}

	// Unmarshal full data
	if err := json.Unmarshal(input.Data, action); err != nil {
		return fmt.Errorf("failed to parse specific action structure: %w", err)
	}

	// Marshal to TOML and write
	data, err := toml.Marshal(action)
	if err != nil {
		return fmt.Errorf("failed to marshal action to toml: %w", err)
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write action file: %w", err)
	}

	return nil
}

// DeleteAction deletes an action file
func DeleteAction(input GetActionInput) error {
	if input.CollectionPath == "" || input.FileName == "" {
		return fmt.Errorf("collection path and file name are required")
	}
	if !strings.HasSuffix(input.FileName, ".toml") {
		input.FileName += ".toml"
	}

	fullPath := filepath.Join(input.CollectionPath, input.SubPath, input.FileName)
	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("action file not found: %s", fullPath)
		}
		return fmt.Errorf("failed to delete action file: %w", err)
	}
	return nil
}

// GetAction reads an action file and returns it as a strong typed struct (or interface{})
func GetAction(input GetActionInput) (interface{}, error) {
	if input.CollectionPath == "" || input.FileName == "" {
		return nil, fmt.Errorf("collection path and file name are required")
	}
	if !strings.HasSuffix(input.FileName, ".toml") {
		input.FileName += ".toml"
	}

	fullPath := filepath.Join(input.CollectionPath, input.SubPath, input.FileName)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	// Unmarshal into BaseAction to get Kind
	var baseAction types.BaseAction
	if err := toml.Unmarshal(data, &baseAction); err != nil {
		return nil, fmt.Errorf("failed to parse TOML metadata: %w", err)
	}

	// Get specific struct
	action, err := getActionStruct(baseAction.Kind)
	if err != nil {
		return nil, err
	}

	// Unmarshal fully
	if err := toml.Unmarshal(data, action); err != nil {
		return nil, fmt.Errorf("failed to parse specific action structure: %w", err)
	}

	return action, nil
}
