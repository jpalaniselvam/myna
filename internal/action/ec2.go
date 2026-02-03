package action

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/jpalaniselvam/myna/internal/parser"
	baseTypes "github.com/jpalaniselvam/myna/internal/types"
)

func runEC2Action(cfg aws.Config, content []byte) (interface{}, error) {
	var action baseTypes.EC2Action
	if err := parser.Unmarshal(content, &action); err != nil {
		return nil, fmt.Errorf("failed to parse ec2 action: %w", err)
	}

	client := ec2.NewFromConfig(cfg)
	ctx := context.TODO()

	var output interface{}
	var err error

	switch action.Kind {
	case baseTypes.KindEC2DescribeInstances:
		output, err = runEC2DescribeInstances(ctx, client, action)
	case baseTypes.KindEC2StartInstances:
		output, err = runEC2StartInstances(ctx, client, action)
	case baseTypes.KindEC2StopInstances:
		output, err = runEC2StopInstances(ctx, client, action)
	case baseTypes.KindEC2RebootInstances:
		output, err = runEC2RebootInstances(ctx, client, action)
	case baseTypes.KindEC2TerminateInstances:
		output, err = runEC2TerminateInstances(ctx, client, action)
	default:
		return nil, fmt.Errorf("unsupported ec2 kind: %s", action.Kind)
	}

	if err != nil {
		return nil, err
	}

	return output, nil
}

func runEC2DescribeInstances(ctx context.Context, client *ec2.Client, action baseTypes.EC2Action) (interface{}, error) {
	input := &ec2.DescribeInstancesInput{}

	if len(action.EC2.InstanceIds) > 0 {
		input.InstanceIds = action.EC2.InstanceIds
	}

	if len(action.EC2.Filters) > 0 {
		var filters []types.Filter
		for k, v := range action.EC2.Filters {
			filters = append(filters, types.Filter{
				Name:   aws.String(k),
				Values: v,
			})
		}
		input.Filters = filters
	}

	if action.EC2.DryRun {
		input.DryRun = aws.Bool(true)
	}

	resp, err := client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %w", err)
	}

	// Create a simplified output structure
	type InstanceOutput struct {
		InstanceId       string            `json:"instance_id"`
		InstanceType     string            `json:"instance_type"`
		State            string            `json:"state"`
		PublicIpAddress  string            `json:"public_ip_address,omitempty"`
		PrivateIpAddress string            `json:"private_ip_address,omitempty"`
		Tags             map[string]string `json:"tags,omitempty"`
		LaunchTime       string            `json:"launch_time"`
	}

	var instances []InstanceOutput

	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			tags := make(map[string]string)
			for _, t := range instance.Tags {
				tags[aws.ToString(t.Key)] = aws.ToString(t.Value)
			}

			instances = append(instances, InstanceOutput{
				InstanceId:       aws.ToString(instance.InstanceId),
				InstanceType:     string(instance.InstanceType),
				State:            string(instance.State.Name),
				PublicIpAddress:  aws.ToString(instance.PublicIpAddress),
				PrivateIpAddress: aws.ToString(instance.PrivateIpAddress),
				Tags:             tags,
				LaunchTime:       instance.LaunchTime.String(),
			})
		}
	}

	return map[string]interface{}{
		"instances": instances,
	}, nil
}

func runEC2StartInstances(ctx context.Context, client *ec2.Client, action baseTypes.EC2Action) (interface{}, error) {
	if len(action.EC2.InstanceIds) == 0 {
		return nil, fmt.Errorf("instance_ids is required for start_instances")
	}

	input := &ec2.StartInstancesInput{
		InstanceIds: action.EC2.InstanceIds,
	}

	if action.EC2.DryRun {
		input.DryRun = aws.Bool(true)
	}

	resp, err := client.StartInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to start instances: %w", err)
	}

	return map[string]interface{}{
		"starting_instances": resp.StartingInstances,
	}, nil
}

func runEC2StopInstances(ctx context.Context, client *ec2.Client, action baseTypes.EC2Action) (interface{}, error) {
	if len(action.EC2.InstanceIds) == 0 {
		return nil, fmt.Errorf("instance_ids is required for stop_instances")
	}

	input := &ec2.StopInstancesInput{
		InstanceIds: action.EC2.InstanceIds,
	}

	if action.EC2.Hibernate {
		input.Hibernate = aws.Bool(true)
	}
	if action.EC2.Force {
		input.Force = aws.Bool(true)
	}
	if action.EC2.DryRun {
		input.DryRun = aws.Bool(true)
	}

	resp, err := client.StopInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to stop instances: %w", err)
	}

	return map[string]interface{}{
		"stopping_instances": resp.StoppingInstances,
	}, nil
}

func runEC2RebootInstances(ctx context.Context, client *ec2.Client, action baseTypes.EC2Action) (interface{}, error) {
	if len(action.EC2.InstanceIds) == 0 {
		return nil, fmt.Errorf("instance_ids is required for reboot_instances")
	}

	input := &ec2.RebootInstancesInput{
		InstanceIds: action.EC2.InstanceIds,
	}

	if action.EC2.DryRun {
		input.DryRun = aws.Bool(true)
	}

	_, err := client.RebootInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to reboot instances: %w", err)
	}

	return map[string]string{"status": "success"}, nil
}

func runEC2TerminateInstances(ctx context.Context, client *ec2.Client, action baseTypes.EC2Action) (interface{}, error) {
	if len(action.EC2.InstanceIds) == 0 {
		return nil, fmt.Errorf("instance_ids is required for terminate_instances")
	}

	input := &ec2.TerminateInstancesInput{
		InstanceIds: action.EC2.InstanceIds,
	}

	if action.EC2.DryRun {
		input.DryRun = aws.Bool(true)
	}

	resp, err := client.TerminateInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to terminate instances: %w", err)
	}

	return map[string]interface{}{
		"terminating_instances": resp.TerminatingInstances,
	}, nil
}
