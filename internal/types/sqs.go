package types

// SQSAction represents the structure of an SQS action TOML file
type SQSAction struct {
	BaseAction
	SQS SQSConfig `toml:"sqs" json:"sqs"`
}

// SQSConfig defines the SQS configuration parameters
type SQSConfig struct {
	QueueURL               string                       `toml:"queue_url" json:"queue_url"`
	DelaySeconds           int32                        `toml:"delay_seconds" json:"delay_seconds"`
	MessageGroupId         string                       `toml:"message_group_id,omitempty" json:"message_group_id,omitempty"`
	MessageDeduplicationId string                       `toml:"message_deduplication_id,omitempty" json:"message_deduplication_id,omitempty"`
	MaxNumberOfMessages    int32                        `toml:"max_number_of_messages" json:"max_number_of_messages"`
	WaitTimeSeconds        int32                        `toml:"wait_time_seconds" json:"wait_time_seconds"`
	VisibilityTimeout      int32                        `toml:"visibility_timeout" json:"visibility_timeout"`
	AttributeNames         []string                     `toml:"attribute_names,omitempty" json:"attribute_names,omitempty"`
	MessageAttributeNames  []string                     `toml:"message_attribute_names,omitempty" json:"message_attribute_names,omitempty"`
	ReceiptHandle          string                       `toml:"receipt_handle,omitempty" json:"receipt_handle,omitempty"`
	SourceARN              string                       `toml:"source_arn,omitempty" json:"source_arn,omitempty"`
	DestinationARN         string                       `toml:"destination_arn,omitempty" json:"destination_arn,omitempty"`
	MaxMessagesPerSecond   int32                        `toml:"max_messages_per_second" json:"max_messages_per_second"`
	MessageAttributes      map[string]SQSMessageAttrVal `toml:"message_attributes,omitempty" json:"message_attributes,omitempty"`
}

type SQSMessageAttrVal struct {
	DataType    string `toml:"DataType" json:"DataType"`
	StringValue string `toml:"StringValue" json:"StringValue"`
	BinaryValue []byte `toml:"BinaryValue" json:"BinaryValue"`
}
