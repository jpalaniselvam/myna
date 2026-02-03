package action

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"

	"github.com/jpalaniselvam/myna/internal/parser"
	baseTypes "github.com/jpalaniselvam/myna/internal/types"
)

func runSFNAction(cfg aws.Config, content []byte, payload []byte) (interface{}, error) {
	var action baseTypes.SFNAction
	if err := parser.Unmarshal(content, &action); err != nil {
		return nil, fmt.Errorf("failed to parse sfn action: %w", err)
	}

	client := sfn.NewFromConfig(cfg)
	ctx := context.TODO()

	var output interface{}
	var err error

	switch action.Kind {
	case baseTypes.KindSFNListStateMachines:
		output, err = runSFNListStateMachines(ctx, client)
	case baseTypes.KindSFNStartExecution:
		output, err = runSFNStartExecution(ctx, client, action, payload)
	case baseTypes.KindSFNDescribeExecution:
		output, err = runSFNDescribeExecution(ctx, client, action)
	case baseTypes.KindSFNStopExecution:
		output, err = runSFNStopExecution(ctx, client, action)
	default:
		return nil, fmt.Errorf("unsupported sfn kind: %s", action.Kind)
	}

	if err != nil {
		return nil, err
	}

	return output, nil
}

// runSFNListStateMachines handles sfn.list_state_machines
func runSFNListStateMachines(ctx context.Context, client *sfn.Client) (interface{}, error) {
	input := &sfn.ListStateMachinesInput{}

	resp, err := client.ListStateMachines(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list state machines: %w", err)
	}

	type StateMachineOutput struct {
		StateMachineArn string `json:"state_machine_arn"`
		Name            string `json:"name"`
		Type            string `json:"type"`
		CreationDate    string `json:"creation_date"`
	}

	stateMachines := make([]StateMachineOutput, 0, len(resp.StateMachines))
	for _, sm := range resp.StateMachines {
		stateMachines = append(stateMachines, StateMachineOutput{
			StateMachineArn: aws.ToString(sm.StateMachineArn),
			Name:            aws.ToString(sm.Name),
			Type:            string(sm.Type),
			CreationDate:    sm.CreationDate.String(),
		})
	}

	return map[string]interface{}{
		"state_machines": stateMachines,
	}, nil
}

// runSFNStartExecution handles sfn.start_execution
func runSFNStartExecution(ctx context.Context, client *sfn.Client, action baseTypes.SFNAction, payload []byte) (interface{}, error) {
	input := &sfn.StartExecutionInput{
		StateMachineArn: aws.String(action.SFN.StateMachineARN),
	}

	if action.SFN.Name != "" {
		input.Name = aws.String(action.SFN.Name)
	}

	if len(payload) > 0 {
		input.Input = aws.String(string(payload))
	}

	resp, err := client.StartExecution(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to start execution: %w", err)
	}

	return map[string]string{
		"execution_arn": aws.ToString(resp.ExecutionArn),
		"start_date":    resp.StartDate.String(),
	}, nil
}

// runSFNDescribeExecution handles sfn.describe_execution
func runSFNDescribeExecution(ctx context.Context, client *sfn.Client, action baseTypes.SFNAction) (interface{}, error) {
	input := &sfn.DescribeExecutionInput{
		ExecutionArn: aws.String(action.SFN.ExecutionARN),
	}

	resp, err := client.DescribeExecution(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe execution: %w", err)
	}

	output := map[string]interface{}{
		"execution_arn":     aws.ToString(resp.ExecutionArn),
		"state_machine_arn": aws.ToString(resp.StateMachineArn),
		"name":              aws.ToString(resp.Name),
		"status":            string(resp.Status),
		"start_date":        resp.StartDate.String(),
	}

	if resp.StopDate != nil {
		output["stop_date"] = resp.StopDate.String()
	}
	if resp.Input != nil {
		output["input"] = aws.ToString(resp.Input)
	}
	if resp.Output != nil {
		output["output"] = aws.ToString(resp.Output)
	}
	if resp.Error != nil {
		output["error"] = aws.ToString(resp.Error)
	}
	if resp.Cause != nil {
		output["cause"] = aws.ToString(resp.Cause)
	}

	return output, nil
}

// runSFNStopExecution handles sfn.stop_execution
func runSFNStopExecution(ctx context.Context, client *sfn.Client, action baseTypes.SFNAction) (interface{}, error) {
	input := &sfn.StopExecutionInput{
		ExecutionArn: aws.String(action.SFN.ExecutionARN),
	}

	if action.SFN.Error != "" {
		input.Error = aws.String(action.SFN.Error)
	}
	if action.SFN.Cause != "" {
		input.Cause = aws.String(action.SFN.Cause)
	}

	resp, err := client.StopExecution(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to stop execution: %w", err)
	}

	return map[string]string{
		"stop_date": resp.StopDate.String(),
	}, nil
}
