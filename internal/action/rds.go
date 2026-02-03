package action

import (
	"context"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"

	"github.com/jpalaniselvam/myna/internal/parser"
	baseTypes "github.com/jpalaniselvam/myna/internal/types"
)

func runRDSAction(cfg aws.Config, content []byte) (interface{}, error) {
	var action baseTypes.RDSAction
	if err := parser.Unmarshal(content, &action); err != nil {
		return nil, fmt.Errorf("failed to parse rds action: %w", err)
	}

	client := rds.NewFromConfig(cfg)
	ctx := context.TODO()

	switch action.Kind {
	case baseTypes.KindRDSDescribeDBInstances:
		return runRDSDescribeInstances(ctx, client, action)
	case baseTypes.KindRDSStartDBInstance:
		return runRDSStartInstance(ctx, client, action)
	case baseTypes.KindRDSStopDBInstance:
		return runRDSStopInstance(ctx, client, action)
	case baseTypes.KindRDSRebootDBInstance:
		return runRDSRebootInstance(ctx, client, action)
	default:
		return nil, fmt.Errorf("unsupported rds kind: %s", action.Kind)
	}
}

func runRDSDescribeInstances(ctx context.Context, client *rds.Client, action baseTypes.RDSAction) (interface{}, error) {
	input := &rds.DescribeDBInstancesInput{}

	if action.RDS.DBInstanceIdentifier != "" {
		input.DBInstanceIdentifier = aws.String(action.RDS.DBInstanceIdentifier)
	}

	if len(action.RDS.Filters) > 0 {
		var filters []types.Filter
		for k, v := range action.RDS.Filters {
			filters = append(filters, types.Filter{
				Name:   aws.String(k),
				Values: v,
			})
		}
		// Sort for deterministic behavior
		sort.Slice(filters, func(i, j int) bool {
			return *filters[i].Name < *filters[j].Name
		})
		input.Filters = filters
	}

	resp, err := client.DescribeDBInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe rds instances: %w", err)
	}

	instances := make([]map[string]interface{}, 0, len(resp.DBInstances))
	for _, instance := range resp.DBInstances {
		instances = append(instances, simplifyDBInstance(instance))
	}

	return map[string]interface{}{
		"db_instances": instances,
	}, nil
}

func runRDSStartInstance(ctx context.Context, client *rds.Client, action baseTypes.RDSAction) (interface{}, error) {
	input := &rds.StartDBInstanceInput{
		DBInstanceIdentifier: aws.String(action.RDS.DBInstanceIdentifier),
	}

	resp, err := client.StartDBInstance(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to start rds instance: %w", err)
	}

	return simplifyDBInstance(*resp.DBInstance), nil
}

func runRDSStopInstance(ctx context.Context, client *rds.Client, action baseTypes.RDSAction) (interface{}, error) {
	input := &rds.StopDBInstanceInput{
		DBInstanceIdentifier: aws.String(action.RDS.DBInstanceIdentifier),
	}

	if action.RDS.DBSnapshotIdentifier != "" {
		input.DBSnapshotIdentifier = aws.String(action.RDS.DBSnapshotIdentifier)
	}

	resp, err := client.StopDBInstance(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to stop rds instance: %w", err)
	}

	return simplifyDBInstance(*resp.DBInstance), nil
}

func runRDSRebootInstance(ctx context.Context, client *rds.Client, action baseTypes.RDSAction) (interface{}, error) {
	input := &rds.RebootDBInstanceInput{
		DBInstanceIdentifier: aws.String(action.RDS.DBInstanceIdentifier),
	}

	if action.RDS.ForceFailover {
		input.ForceFailover = aws.Bool(true)
	}

	resp, err := client.RebootDBInstance(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to reboot rds instance: %w", err)
	}

	return simplifyDBInstance(*resp.DBInstance), nil
}

func simplifyDBInstance(instance types.DBInstance) map[string]interface{} {
	out := map[string]interface{}{
		"db_instance_identifier": aws.ToString(instance.DBInstanceIdentifier),
		"db_instance_status":     aws.ToString(instance.DBInstanceStatus),
		"engine":                 aws.ToString(instance.Engine),
		"engine_version":         aws.ToString(instance.EngineVersion),
		"db_instance_class":      aws.ToString(instance.DBInstanceClass),
		"availability_zone":      aws.ToString(instance.AvailabilityZone),
		"master_username":        aws.ToString(instance.MasterUsername),
		"db_name":                aws.ToString(instance.DBName),
		"allocated_storage":      instance.AllocatedStorage,
	}

	if instance.Endpoint != nil {
		out["endpoint"] = map[string]interface{}{
			"address": aws.ToString(instance.Endpoint.Address),
			"port":    instance.Endpoint.Port,
		}
	}

	if instance.InstanceCreateTime != nil {
		out["instance_create_time"] = instance.InstanceCreateTime.String()
	}

	return out
}
