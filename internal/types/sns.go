package types

import "encoding/json"

// SNSAction represents the structure of an sns.publish action TOML file
type SNSAction struct {
	BaseAction
	SNS SNSConfig `toml:"sns" json:"sns"`
}

// SNSConfig defines the SNS publication configuration
type SNSConfig struct {
	TopicARN               string                         `toml:"topic_arn,omitempty" json:"topic_arn,omitempty"`
	TargetARN              string                         `toml:"target_arn,omitempty" json:"target_arn,omitempty"`
	PhoneNumber            string                         `toml:"phone_number,omitempty" json:"phone_number,omitempty"`
	Subject                string                         `toml:"subject,omitempty" json:"subject,omitempty"`
	MessageStructure       string                         `toml:"message_structure,omitempty" json:"message_structure,omitempty"`
	MessageGroupID         string                         `toml:"message_group_id,omitempty" json:"message_group_id,omitempty"`
	MessageDeduplicationID string                         `toml:"message_deduplication_id,omitempty" json:"message_deduplication_id,omitempty"`
	MessageAttributes      map[string]SNSMessageAttribute `toml:"message_attributes,omitempty" json:"message_attributes,omitempty"`
}

// SNSMessageAttribute defines a single message attribute
type SNSMessageAttribute struct {
	DataType    string `toml:"DataType" json:"DataType"`
	StringValue string `toml:"StringValue" json:"StringValue"`
	BinaryValue []byte `toml:"BinaryValue" json:"BinaryValue"`
}

// SNSOutput represents the output of an SNS publish execution
type SNSOutput struct {
	MessageID      string `json:"message_id"`
	SequenceNumber string `json:"sequence_number,omitempty"`
}

func (s *SNSOutput) String() string {
	b, _ := json.MarshalIndent(s, "", "  ")
	return string(b)
}
