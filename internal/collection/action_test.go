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
