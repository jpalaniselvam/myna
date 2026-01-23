package action

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"

	"github.com/jpalaniselvam/myna/internal/parser"
	baseTypes "github.com/jpalaniselvam/myna/internal/types"
)

// LambdaAction represents the structure of a lambda.invoke action TOML file
type LambdaAction struct {
	baseTypes.BaseAction
	Lambda LambdaConfig `toml:"lambda"`
}

// LambdaConfig defines the lambda function configuration
type LambdaConfig struct {
	FunctionName   string                 `toml:"function_name"`
	InvocationType string                 `toml:"invocation_type"` // RequestResponse (default), Event, DryRun
	Qualifier      string                 `toml:"qualifier"`
	ClientContext  map[string]interface{} `toml:"client_context"`
}

// LambdaOutput represents the output of a lambda execution
type LambdaOutput struct {
	StatusCode      int                    `json:"status_code"`
	ExecutedVersion string                 `json:"executed_version"`
	Payload         map[string]interface{} `json:"payload"`
	Logs            string                 `json:"logs,omitempty"`
	RawPayload      []byte                 `json:"-"` // Internal use, for when response is not JSON
}

func (l *LambdaOutput) String() string {
	b, _ := json.MarshalIndent(l, "", "  ")
	return string(b)
}

func runLambdaAction(cfg aws.Config, content []byte, payload []byte) error {
	var action LambdaAction
	if err := parser.Unmarshal(content, &action); err != nil {
		return fmt.Errorf("failed to parse lambda action: %w", err)
	}

	input, err := buildInvokeInput(action, payload)
	if err != nil {
		return err
	}

	client := lambda.NewFromConfig(cfg)

	ctx := context.TODO()
	if action.TimeoutMs > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(action.TimeoutMs)*time.Millisecond)
		defer cancel()
	}

	resp, err := client.Invoke(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to invoke lambda: %w", err)
	}

	output, err := processLambdaOutput(resp)
	if err != nil {
		return err
	}

	fmt.Println(output.String())

	return nil
}

func buildInvokeInput(action LambdaAction, payload []byte) (*lambda.InvokeInput, error) {

	var invocationType = types.InvocationType(action.Lambda.InvocationType)

	input := &lambda.InvokeInput{
		FunctionName:   aws.String(action.Lambda.FunctionName),
		Payload:        payload,
		InvocationType: invocationType,
		Qualifier:      aws.String(action.Lambda.Qualifier),
	}

	if invocationType == types.InvocationTypeRequestResponse {
		input.LogType = types.LogTypeTail
	}

	if action.Lambda.Qualifier == "" {
		input.Qualifier = nil
	}

	if len(action.Lambda.ClientContext) > 0 {
		ccBytes, err := json.Marshal(action.Lambda.ClientContext)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal client context: %w", err)
		}
		encodedCC := base64.StdEncoding.EncodeToString(ccBytes)
		input.ClientContext = aws.String(encodedCC)
	}

	return input, nil
}

func processLambdaOutput(resp *lambda.InvokeOutput) (*LambdaOutput, error) {
	output := &LambdaOutput{
		StatusCode:      int(resp.StatusCode),
		ExecutedVersion: aws.ToString(resp.ExecutedVersion),
		RawPayload:      resp.Payload,
	}

	if resp.LogResult != nil {
		// LogResult is base64 encoded
		b, err := base64.StdEncoding.DecodeString(*resp.LogResult)
		if err == nil {
			output.Logs = string(b)
		}
	}

	var respPayload map[string]interface{}
	// Only try to unmarshal if we have a payload
	if len(resp.Payload) > 0 {
		if err := json.Unmarshal(resp.Payload, &respPayload); err == nil {
			output.Payload = respPayload
		}
	}

	return output, nil
}
