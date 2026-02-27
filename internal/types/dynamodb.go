package types

// DynamoDBAction represents the structure of a DynamoDB action TOML file
type DynamoDBAction struct {
	BaseAction
	DynamoDB DynamoDBConfig `toml:"dynamodb" json:"dynamodb"`
}

// DynamoDBConfig defines the DynamoDB configuration parameters
type DynamoDBConfig struct {
	TableName                 string                 `toml:"table_name" json:"table_name"`
	Key                       map[string]interface{} `toml:"key" json:"key"`
	UpdateExpression          string                 `toml:"update_expression,omitempty" json:"update_expression,omitempty"`
	ConditionExpression       string                 `toml:"condition_expression,omitempty" json:"condition_expression,omitempty"`
	KeyConditionExpression    string                 `toml:"key_condition_expression,omitempty" json:"key_condition_expression,omitempty"`
	FilterExpression          string                 `toml:"filter_expression,omitempty" json:"filter_expression,omitempty"`
	ExpressionAttributeNames  map[string]string      `toml:"expression_attribute_names,omitempty" json:"expression_attribute_names,omitempty"`
	ExpressionAttributeValues map[string]interface{} `toml:"expression_attribute_values,omitempty" json:"expression_attribute_values,omitempty"`
	IndexName                 string                 `toml:"index_name,omitempty" json:"index_name,omitempty"`
	Limit                     int32                  `toml:"limit,omitempty" json:"limit,omitempty"`
	ConsistentRead            bool                   `toml:"consistent_read,omitempty" json:"consistent_read,omitempty"`
}
