package action

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"

	"github.com/jpalaniselvam/myna/internal/parser"
	baseTypes "github.com/jpalaniselvam/myna/internal/types"
)

func runSNSAction(cfg aws.Config, content []byte, payload []byte) (interface{}, error) {
	var action baseTypes.SNSAction
	if err := parser.Unmarshal(content, &action); err != nil {
		return nil, fmt.Errorf("failed to parse sns action: %w", err)
	}

	input, err := buildPublishInput(action, payload)
	if err != nil {
		return nil, err
	}

	client := sns.NewFromConfig(cfg)

	ctx := context.TODO()
	resp, err := client.Publish(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to publish to sns: %w", err)
	}

	output := &baseTypes.SNSOutput{
		MessageID:      aws.ToString(resp.MessageId),
		SequenceNumber: aws.ToString(resp.SequenceNumber),
	}

	return output, nil
}

func buildPublishInput(action baseTypes.SNSAction, payload []byte) (*sns.PublishInput, error) {
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
