package action

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"

	"github.com/jpalaniselvam/myna/internal/parser"
	baseTypes "github.com/jpalaniselvam/myna/internal/types"
)

// SNSAction represents the structure of an sns.publish action TOML file
type SNSAction struct {
	baseTypes.BaseAction
	SNS SNSConfig `toml:"sns"`
}

// SNSConfig defines the SNS publication configuration
type SNSConfig struct {
	TopicARN               string                         `toml:"topic_arn"`
	TargetARN              string                         `toml:"target_arn"`
	PhoneNumber            string                         `toml:"phone_number"`
	Subject                string                         `toml:"subject"`
	MessageStructure       string                         `toml:"message_structure"`
	MessageGroupID         string                         `toml:"message_group_id"`
	MessageDeduplicationID string                         `toml:"message_deduplication_id"`
	MessageAttributes      map[string]SNSMessageAttribute `toml:"message_attributes"`
}

// SNSMessageAttribute defines a single message attribute
type SNSMessageAttribute struct {
	DataType    string `toml:"DataType"`
	StringValue string `toml:"StringValue"`
	BinaryValue []byte `toml:"BinaryValue"`
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

func runSNSAction(cfg aws.Config, content []byte, payload []byte) error {
	var action SNSAction
	if err := parser.Unmarshal(content, &action); err != nil {
		return fmt.Errorf("failed to parse sns action: %w", err)
	}

	input, err := buildPublishInput(action, payload)
	if err != nil {
		return err
	}

	client := sns.NewFromConfig(cfg)

	ctx := context.TODO()
	resp, err := client.Publish(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to publish to sns: %w", err)
	}

	output := &SNSOutput{
		MessageID:      aws.ToString(resp.MessageId),
		SequenceNumber: aws.ToString(resp.SequenceNumber),
	}

	fmt.Println(output.String())

	return nil
}

func buildPublishInput(action SNSAction, payload []byte) (*sns.PublishInput, error) {
	message := string(payload)

	input := &sns.PublishInput{
		Message: aws.String(message),
	}

	if action.SNS.TopicARN != "" {
		input.TopicArn = aws.String(action.SNS.TopicARN)
	}
	if action.SNS.TargetARN != "" {
		input.TargetArn = aws.String(action.SNS.TargetARN)
	}
	if action.SNS.PhoneNumber != "" {
		input.PhoneNumber = aws.String(action.SNS.PhoneNumber)
	}
	if action.SNS.Subject != "" {
		input.Subject = aws.String(action.SNS.Subject)
	}
	if action.SNS.MessageStructure != "" {
		input.MessageStructure = aws.String(action.SNS.MessageStructure)
	}
	if action.SNS.MessageGroupID != "" {
		input.MessageGroupId = aws.String(action.SNS.MessageGroupID)
	}
	if action.SNS.MessageDeduplicationID != "" {
		input.MessageDeduplicationId = aws.String(action.SNS.MessageDeduplicationID)
	}

	if len(action.SNS.MessageAttributes) > 0 {
		input.MessageAttributes = make(map[string]types.MessageAttributeValue)
		for key, attr := range action.SNS.MessageAttributes {
			val := types.MessageAttributeValue{
				DataType: aws.String(attr.DataType),
			}
			if attr.StringValue != "" {
				// Apply substitution to attribute values too
				val.StringValue = aws.String(attr.StringValue)
			}
			if len(attr.BinaryValue) > 0 {
				val.BinaryValue = attr.BinaryValue
			}
			input.MessageAttributes[key] = val
		}
	}

	return input, nil
}
