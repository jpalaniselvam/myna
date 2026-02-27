package action

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/jpalaniselvam/myna/internal/parser"
	baseTypes "github.com/jpalaniselvam/myna/internal/types"
)

func runDynamoDBAction(cfg aws.Config, content []byte, payload []byte) (interface{}, error) {
	var action baseTypes.DynamoDBAction
	if err := parser.Unmarshal(content, &action); err != nil {
		return nil, fmt.Errorf("failed to parse dynamodb action: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	ctx := context.TODO()

	switch action.Kind {
	case baseTypes.KindDynamoDBListTables:
		return runDynamoDBListTables(ctx, client, action)
	case baseTypes.KindDynamoDBPutItem:
		return runDynamoDBPutItem(ctx, client, action, payload)
	case baseTypes.KindDynamoDBGetItem:
		return runDynamoDBGetItem(ctx, client, action)
	case baseTypes.KindDynamoDBUpdateItem:
		return runDynamoDBUpdateItem(ctx, client, action)
	case baseTypes.KindDynamoDBDeleteItem:
		return runDynamoDBDeleteItem(ctx, client, action)
	case baseTypes.KindDynamoDBQuery:
		return runDynamoDBQuery(ctx, client, action)
	case baseTypes.KindDynamoDBScan:
		return runDynamoDBScan(ctx, client, action)
	default:
		return nil, fmt.Errorf("unsupported dynamodb kind: %s", action.Kind)
	}
}

func runDynamoDBListTables(ctx context.Context, client *dynamodb.Client, action baseTypes.DynamoDBAction) (interface{}, error) {
	resp, err := client.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}
	return map[string]interface{}{
		"table_names": resp.TableNames,
	}, nil
}

func runDynamoDBPutItem(ctx context.Context, client *dynamodb.Client, action baseTypes.DynamoDBAction, payload []byte) (interface{}, error) {
	var item map[string]interface{}
	if err := json.Unmarshal(payload, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal item to attribute values: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(action.DynamoDB.TableName),
		Item:      av,
	}

	if action.DynamoDB.ConditionExpression != "" {
		input.ConditionExpression = aws.String(action.DynamoDB.ConditionExpression)
	}

	_, err = client.PutItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to put item: %w", err)
	}

	return map[string]string{"status": "success"}, nil
}

func runDynamoDBGetItem(ctx context.Context, client *dynamodb.Client, action baseTypes.DynamoDBAction) (interface{}, error) {
	key, err := attributevalue.MarshalMap(action.DynamoDB.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal key: %w", err)
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(action.DynamoDB.TableName),
		Key:       key,
	}

	if action.DynamoDB.ConsistentRead {
		input.ConsistentRead = aws.Bool(true)
	}

	resp, err := client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	if resp.Item == nil {
		return nil, nil
	}

	var result map[string]interface{}
	if err := attributevalue.UnmarshalMap(resp.Item, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return result, nil
}

func runDynamoDBUpdateItem(ctx context.Context, client *dynamodb.Client, action baseTypes.DynamoDBAction) (interface{}, error) {
	key, err := attributevalue.MarshalMap(action.DynamoDB.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal key: %w", err)
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(action.DynamoDB.TableName),
		Key:       key,
	}

	if action.DynamoDB.UpdateExpression != "" {
		input.UpdateExpression = aws.String(action.DynamoDB.UpdateExpression)
	}
	if action.DynamoDB.ConditionExpression != "" {
		input.ConditionExpression = aws.String(action.DynamoDB.ConditionExpression)
	}

	if len(action.DynamoDB.ExpressionAttributeNames) > 0 {
		input.ExpressionAttributeNames = action.DynamoDB.ExpressionAttributeNames
	}

	if len(action.DynamoDB.ExpressionAttributeValues) > 0 {
		avs, err := attributevalue.MarshalMap(action.DynamoDB.ExpressionAttributeValues)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal expression attribute values: %w", err)
		}
		input.ExpressionAttributeValues = avs
	}

	_, err = client.UpdateItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return map[string]string{"status": "success"}, nil
}

func runDynamoDBDeleteItem(ctx context.Context, client *dynamodb.Client, action baseTypes.DynamoDBAction) (interface{}, error) {
	key, err := attributevalue.MarshalMap(action.DynamoDB.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal key: %w", err)
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(action.DynamoDB.TableName),
		Key:       key,
	}

	if action.DynamoDB.ConditionExpression != "" {
		input.ConditionExpression = aws.String(action.DynamoDB.ConditionExpression)
	}

	_, err = client.DeleteItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to delete item: %w", err)
	}

	return map[string]string{"status": "success"}, nil
}

func runDynamoDBQuery(ctx context.Context, client *dynamodb.Client, action baseTypes.DynamoDBAction) (interface{}, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(action.DynamoDB.TableName),
		KeyConditionExpression: aws.String(action.DynamoDB.KeyConditionExpression),
	}

	if action.DynamoDB.IndexName != "" {
		input.IndexName = aws.String(action.DynamoDB.IndexName)
	}

	if action.DynamoDB.Limit > 0 {
		input.Limit = aws.Int32(action.DynamoDB.Limit)
	}

	if len(action.DynamoDB.ExpressionAttributeNames) > 0 {
		input.ExpressionAttributeNames = action.DynamoDB.ExpressionAttributeNames
	}

	if len(action.DynamoDB.ExpressionAttributeValues) > 0 {
		avs, err := attributevalue.MarshalMap(action.DynamoDB.ExpressionAttributeValues)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal expression attribute values: %w", err)
		}
		input.ExpressionAttributeValues = avs
	}

	resp, err := client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	var items []map[string]interface{}
	if err := attributevalue.UnmarshalListOfMaps(resp.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	return map[string]interface{}{
		"items": items,
		"count": resp.Count,
	}, nil
}

func runDynamoDBScan(ctx context.Context, client *dynamodb.Client, action baseTypes.DynamoDBAction) (interface{}, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(action.DynamoDB.TableName),
	}

	if action.DynamoDB.FilterExpression != "" {
		input.FilterExpression = aws.String(action.DynamoDB.FilterExpression)
	}

	if action.DynamoDB.Limit > 0 {
		input.Limit = aws.Int32(action.DynamoDB.Limit)
	}

	if len(action.DynamoDB.ExpressionAttributeNames) > 0 {
		input.ExpressionAttributeNames = action.DynamoDB.ExpressionAttributeNames
	}

	if len(action.DynamoDB.ExpressionAttributeValues) > 0 {
		avs, err := attributevalue.MarshalMap(action.DynamoDB.ExpressionAttributeValues)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal expression attribute values: %w", err)
		}
		input.ExpressionAttributeValues = avs
	}

	resp, err := client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	var items []map[string]interface{}
	if err := attributevalue.UnmarshalListOfMaps(resp.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	return map[string]interface{}{
		"items": items,
		"count": resp.Count,
	}, nil
}
