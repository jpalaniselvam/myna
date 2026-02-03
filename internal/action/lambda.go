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

func runLambdaAction(cfg aws.Config, content []byte, payload []byte) (interface{}, error) {
	var action baseTypes.LambdaAction
	if err := parser.Unmarshal(content, &action); err != nil {
		return nil, fmt.Errorf("failed to parse lambda action: %w", err)
	}

	input, err := buildInvokeInput(action, payload)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("failed to invoke lambda: %w", err)
	}

	output, err := processLambdaOutput(resp)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func buildInvokeInput(action baseTypes.LambdaAction, payload []byte) (*lambda.InvokeInput, error) {

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

func processLambdaOutput(resp *lambda.InvokeOutput) (*baseTypes.LambdaOutput, error) {
	output := &baseTypes.LambdaOutput{
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
