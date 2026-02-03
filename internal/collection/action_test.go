package collection

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/jpalaniselvam/myna/internal/types"
	"github.com/pelletier/go-toml/v2"
)

/*
Summary of Test Cases:
1. CreateAction:
   - SuccessLambda: Verify creation of a Lambda action file with valid configuration.
   - NestedPath: Verify creation of action files in nested directories.
   - UnknownKind: Error when an unsupported action kind is provided.
2. UpdateAction:
   - SuccessUpdate: Verify updating an existing action file updates the content correctly.
   - NonExistent: Error when trying to update a non-existent action.
3. DeleteAction:
   - SuccessDelete: Verify deletion of an existing action file.
   - NonExistent: Error when trying to delete a non-existent action.
4. GetAction:
   - SuccessGet: Verify reading an action file returns the correct struct.
   - NonExistent: Error when trying to get a non-existent action.
*/

func TestCreateAction(t *testing.T) {
	tmpDir := t.TempDir()
	collName := "test-col"
	collPath := filepath.Join(tmpDir, collName)
	Create(tmpDir, collName, "desc")

	t.Run("SuccessLambda", func(t *testing.T) {
		lambdaAction := types.LambdaAction{
			BaseAction: types.BaseAction{
				Metadata: types.Metadata{
					Kind:        types.KindLambdaInvoke,
					Description: "Invoke a lambda",
				},
			},
			Lambda: types.LambdaConfig{
				FunctionName: "my-func",
			},
		}

		actionData, _ := json.Marshal(lambdaAction)

		input := CreateActionInput{
			CollectionPath: collPath,
			SubPath:        "",
			FileName:       "my-lambda.toml",
			Data:           actionData,
		}

		err := CreateAction(input)
		if err != nil {
			t.Fatalf("CreateAction failed: %v", err)
		}

		actionPath := filepath.Join(collPath, "my-lambda.toml")
		data, _ := os.ReadFile(actionPath)
		var m map[string]interface{}
		toml.Unmarshal(data, &m)

		if m["kind"] != "lambda.invoke" {
			t.Errorf("expected kind lambda.invoke, got %v", m["kind"])
		}
		lambdaConfig := m["lambda"].(map[string]interface{})
		if lambdaConfig["function_name"] != "my-func" {
			t.Errorf("expected function_name my-func, got %v", lambdaConfig["function_name"])
		}
	})

	t.Run("NestedPath", func(t *testing.T) {
		sqsAction := types.SQSAction{
			BaseAction: types.BaseAction{
				Metadata: types.Metadata{
					Kind:        types.KindSQSSend,
					Description: "Nested SQS",
				},
			},
			SQS: types.SQSConfig{
				QueueURL: "https://sqs.us-east-1.amazonaws.com/123/queue",
			},
		}

		actionData, _ := json.Marshal(sqsAction)

		input := CreateActionInput{
			CollectionPath: collPath,
			SubPath:        "subdir",
			FileName:       "nested.toml",
			Data:           actionData,
		}

		err := CreateAction(input)
		if err != nil {
			t.Fatalf("Nested CreateAction failed: %v", err)
		}

		actionPath := filepath.Join(collPath, "subdir", "nested.toml")
		if _, err := os.Stat(actionPath); os.IsNotExist(err) {
			t.Error("expected nested action file to exist")
		}
	})

	t.Run("UnknownKind", func(t *testing.T) {
		baseAction := types.BaseAction{
			Metadata: types.Metadata{
				Kind: types.Kind("unknown.kind"),
			},
		}
		actionData, _ := json.Marshal(baseAction)

		input := CreateActionInput{
			CollectionPath: collPath,
			SubPath:        "",
			FileName:       "unknown.toml",
			Data:           actionData,
		}
		if err := CreateAction(input); err == nil {
			t.Error("expected error for unknown kind")
		}
	})
}

func TestUpdateAction(t *testing.T) {
	tmpDir := t.TempDir()
	collName := "test-col-update"
	collPath := filepath.Join(tmpDir, collName)
	Create(tmpDir, collName, "desc")

	// Create initial action
	lambdaAction := types.LambdaAction{
		BaseAction: types.BaseAction{
			Metadata: types.Metadata{
				Kind:        types.KindLambdaInvoke,
				Description: "Initial description",
			},
		},
		Lambda: types.LambdaConfig{
			FunctionName: "initial-func",
		},
	}
	actionData, _ := json.Marshal(lambdaAction)
	createInput := CreateActionInput{
		CollectionPath: collPath,
		FileName:       "update-test.toml",
		Data:           actionData,
	}
	CreateAction(createInput)

	t.Run("SuccessUpdate", func(t *testing.T) {
		updatedAction := lambdaAction
		updatedAction.Metadata.Description = "Updated description"
		updatedAction.Lambda.FunctionName = "updated-func"

		updatedData, _ := json.Marshal(updatedAction)

		updateInput := UpdateActionInput{
			CollectionPath: collPath,
			FileName:       "update-test.toml",
			Data:           updatedData,
		}

		err := UpdateAction(updateInput)
		if err != nil {
			t.Fatalf("UpdateAction failed: %v", err)
		}

		// Verify
		actionPath := filepath.Join(collPath, "update-test.toml")
		data, _ := os.ReadFile(actionPath)
		var m map[string]interface{}
		toml.Unmarshal(data, &m)
		if m["description"] != "Updated description" {
			t.Errorf("expected description 'Updated description', got %v", m["description"])
		}
		lambdaConfig := m["lambda"].(map[string]interface{})
		if lambdaConfig["function_name"] != "updated-func" {
			t.Errorf("expected function_name 'updated-func', got %v", lambdaConfig["function_name"])
		}
	})

	t.Run("NonExistent", func(t *testing.T) {
		updateInput := UpdateActionInput{
			CollectionPath: collPath,
			FileName:       "non-existent.toml",
			Data:           actionData,
		}
		if err := UpdateAction(updateInput); err == nil {
			t.Error("expected error when updating non-existent file")
		}
	})
}

func TestDeleteAction(t *testing.T) {
	tmpDir := t.TempDir()
	collName := "test-col-delete"
	collPath := filepath.Join(tmpDir, collName)
	Create(tmpDir, collName, "desc")

	// Create action
	lambdaAction := types.LambdaAction{
		BaseAction: types.BaseAction{
			Metadata: types.Metadata{
				Kind: types.KindLambdaInvoke,
			},
		},
	}
	actionData, _ := json.Marshal(lambdaAction)
	createInput := CreateActionInput{
		CollectionPath: collPath,
		FileName:       "delete-test.toml",
		Data:           actionData,
	}
	CreateAction(createInput)

	t.Run("SuccessDelete", func(t *testing.T) {
		deleteInput := GetActionInput{
			CollectionPath: collPath,
			FileName:       "delete-test.toml",
		}
		err := DeleteAction(deleteInput)
		if err != nil {
			t.Fatalf("DeleteAction failed: %v", err)
		}

		if _, err := os.Stat(filepath.Join(collPath, "delete-test.toml")); !os.IsNotExist(err) {
			t.Error("expected file to be deleted")
		}
	})

	t.Run("NonExistent", func(t *testing.T) {
		deleteInput := GetActionInput{
			CollectionPath: collPath,
			FileName:       "non-existent.toml",
		}
		if err := DeleteAction(deleteInput); err == nil {
			t.Error("expected error when deleting non-existent file")
		}
	})
}

func TestGetAction(t *testing.T) {
	tmpDir := t.TempDir()
	collName := "test-col-get"
	collPath := filepath.Join(tmpDir, collName)
	Create(tmpDir, collName, "desc")

	// Create action
	lambdaAction := types.LambdaAction{
		BaseAction: types.BaseAction{
			Metadata: types.Metadata{
				Kind:        types.KindLambdaInvoke,
				Description: "Get description",
			},
		},
		Lambda: types.LambdaConfig{
			FunctionName: "get-func",
		},
	}
	actionData, _ := json.Marshal(lambdaAction)
	createInput := CreateActionInput{
		CollectionPath: collPath,
		FileName:       "get-test.toml",
		Data:           actionData,
	}
	CreateAction(createInput)

	t.Run("SuccessGet", func(t *testing.T) {
		getInput := GetActionInput{
			CollectionPath: collPath,
			FileName:       "get-test.toml",
		}

		action, err := GetAction(getInput)
		if err != nil {
			t.Fatalf("GetAction failed: %v", err)
		}

		lAction, ok := action.(*types.LambdaAction)
		if !ok {
			t.Fatalf("expected *types.LambdaAction, got %T", action)
		}

		if lAction.Metadata.Description != "Get description" {
			t.Errorf("expected description 'Get description', got %s", lAction.Metadata.Description)
		}
		if lAction.Lambda.FunctionName != "get-func" {
			t.Errorf("expected function name 'get-func', got %s", lAction.Lambda.FunctionName)
		}
	})

	t.Run("NonExistent", func(t *testing.T) {
		getInput := GetActionInput{
			CollectionPath: collPath,
			FileName:       "non-existent.toml",
		}
		if _, err := GetAction(getInput); err == nil {
			t.Error("expected error when getting non-existent file")
		}
	})
}
