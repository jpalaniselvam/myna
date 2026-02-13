package types

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/jpalaniselvam/myna/internal/parser"
)

/*
Summary of Test Cases:
1. BaseAction JSON:
   - Marshal: Verify correct JSON encoding.
   - Unmarshal: Verify correct JSON decoding, including type checking.

2. BaseAction TOML:
   - Marshal: Verify correct TOML encoding using parser package.
   - Unmarshal: Verify correct TOML decoding using parser package.

3. LambdaAction JSON & TOML:
   - Verify specific Lambda fields.
*/

func TestBaseAction_JSON(t *testing.T) {
	// Create a sample BaseAction
	original := BaseAction{
		Metadata: Metadata{
			Version:     "1.0",
			Kind:        KindLambdaInvoke,
			Description: "Test Lambda Invoke",
			TimeoutMs:   5000,
		},
		Credentials: Credentials{
			Region:  "us-east-1",
			Profile: "default",
		},
		Vars: Vars{
			Pre: map[string]interface{}{
				"key1": "value1",
				"key2": 123,
			},
		},
		Payload: Payload{
			Data: `{"foo":"bar"}`,
		},
	}

	t.Run("Marshal", func(t *testing.T) {
		jsonData, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v", err)
		}

		// Simple check that it's valid JSON
		if !json.Valid(jsonData) {
			t.Error("Marshaled data is not valid JSON")
		}
	})

	t.Run("Unmarshal", func(t *testing.T) {
		jsonData, _ := json.Marshal(original)
		var decoded BaseAction
		err := json.Unmarshal(jsonData, &decoded)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Check Metadata
		if decoded.Metadata != original.Metadata {
			t.Errorf("Metadata mismatch: got %+v, want %+v", decoded.Metadata, original.Metadata)
		}

		// Check Payload
		if decoded.Payload != original.Payload {
			t.Errorf("Payload mismatch: got %+v, want %+v", decoded.Payload, original.Payload)
		}

		// Check Vars manually due to int/float64 conversion in JSON
		if decoded.Vars.Pre["key1"] != original.Vars.Pre["key1"] {
			t.Errorf("Vars.Pre['key1'] mismatch: got %v, want %v", decoded.Vars.Pre["key1"], original.Vars.Pre["key1"])
		}

		val2, ok := decoded.Vars.Pre["key2"].(float64)
		if !ok || val2 != 123 {
			t.Errorf("Vars.Pre['key2'] mismatch: got %v (type %T), want 123 (float64)", decoded.Vars.Pre["key2"], decoded.Vars.Pre["key2"])
		}
	})
}

func TestBaseAction_TOML(t *testing.T) {
	// Create a sample BaseAction
	original := BaseAction{
		Metadata: Metadata{
			Version:     "1.0",
			Kind:        KindSQSSend,
			Description: "Test SQS Send",
		},
		Credentials: Credentials{
			Region:  "us-west-2",
			Profile: "default",
		},
		Vars: Vars{
			Pre: map[string]interface{}{
				"variable": "content",
			},
		},
		Payload: Payload{
			File: "payload.json",
		},
	}

	t.Run("Marshal", func(t *testing.T) {
		tomlData, err := parser.Marshal(original)
		if err != nil {
			t.Fatalf("Failed to marshal TOML: %v", err)
		}

		if len(tomlData) == 0 {
			t.Error("Marshaled TOML data is empty")
		}
	})

	t.Run("Unmarshal", func(t *testing.T) {
		tomlData, _ := parser.Marshal(original)
		var decoded BaseAction
		err := parser.Unmarshal(tomlData, &decoded)
		if err != nil {
			t.Fatalf("Failed to unmarshal TOML: %v", err)
		}

		// Verify equality
		if !reflect.DeepEqual(decoded, original) {
			t.Errorf("TOML mismatch:\nGot: %+v\nWant: %+v", decoded, original)
		}
	})
}

func TestLambdaAction(t *testing.T) {
	original := LambdaAction{
		BaseAction: BaseAction{
			Metadata: Metadata{
				Kind: KindLambdaInvoke,
			},
		},
		Lambda: LambdaConfig{
			FunctionName:   "my-function",
			InvocationType: "RequestResponse",
			Qualifier:      "prod",
		},
	}

	t.Run("JSON", func(t *testing.T) {
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("JSON marshal failed: %v", err)
		}

		var decoded LambdaAction
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("JSON unmarshal failed: %v", err)
		}

		if decoded.Lambda.FunctionName != original.Lambda.FunctionName {
			t.Errorf("JSON mismatch: got %v, want %v", decoded.Lambda, original.Lambda)
		}
	})

	t.Run("TOML", func(t *testing.T) {
		data, err := parser.Marshal(original)
		if err != nil {
			t.Fatalf("TOML marshal failed: %v", err)
		}

		var decoded LambdaAction
		if err := parser.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("TOML unmarshal failed: %v", err)
		}

		if decoded.Lambda.FunctionName != original.Lambda.FunctionName {
			t.Errorf("TOML mismatch: got %v, want %v", decoded.Lambda, original.Lambda)
		}
	})
}
