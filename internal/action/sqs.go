package action

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"github.com/jpalaniselvam/myna/internal/parser"
	baseTypes "github.com/jpalaniselvam/myna/internal/types"
)

func runSQSAction(cfg aws.Config, content []byte, payload []byte) (interface{}, error) {
	var action baseTypes.SQSAction
	if err := parser.Unmarshal(content, &action); err != nil {
		return nil, fmt.Errorf("failed to parse sqs action: %w", err)
	}

	client := sqs.NewFromConfig(cfg)
	ctx := context.TODO()

	var output interface{}
	var err error

	switch action.Kind {
	case baseTypes.KindSQSSend:
		output, err = runSQSSend(ctx, client, action, payload)
	case baseTypes.KindSQSReceive:
		output, err = runSQSReceive(ctx, client, action)
	case baseTypes.KindSQSDelete:
		output, err = runSQSDelete(ctx, client, action)
	case baseTypes.KindSQSPurge:
		output, err = runSQSPurge(ctx, client, action)
	case baseTypes.KindSQSMoveTask:
		output, err = runSQSMoveTask(ctx, client, action)
	default:
		return nil, fmt.Errorf("unsupported sqs kind: %s", action.Kind)
	}

	if err != nil {
		return nil, err
	}

	return output, nil
}

// runSQSSend handles sqs.send_message
func runSQSSend(ctx context.Context, client *sqs.Client, action baseTypes.SQSAction, payload []byte) (interface{}, error) {
	body := string(payload)

	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(action.SQS.QueueURL),
		MessageBody: aws.String(body),
	}

	if action.SQS.DelaySeconds > 0 {
		input.DelaySeconds = action.SQS.DelaySeconds
	}
	if action.SQS.MessageGroupId != "" {
		input.MessageGroupId = aws.String(action.SQS.MessageGroupId)
	}
	if action.SQS.MessageDeduplicationId != "" {
		input.MessageDeduplicationId = aws.String(action.SQS.MessageDeduplicationId)
	}

	if len(action.SQS.MessageAttributes) > 0 {
		input.MessageAttributes = make(map[string]types.MessageAttributeValue)
		for k, v := range action.SQS.MessageAttributes {
			mav := types.MessageAttributeValue{
				DataType: aws.String(v.DataType),
			}
			if v.StringValue != "" {
				mav.StringValue = aws.String(v.StringValue)
			}
			if len(v.BinaryValue) > 0 {
				mav.BinaryValue = v.BinaryValue
			}
			input.MessageAttributes[k] = mav
		}
	}

	resp, err := client.SendMessage(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to send sqs message: %w", err)
	}

	return map[string]string{
		"message_id":                aws.ToString(resp.MessageId),
		"md5_of_body":               aws.ToString(resp.MD5OfMessageBody),
		"md5_of_message_attributes": aws.ToString(resp.MD5OfMessageAttributes),
		"sequence_number":           aws.ToString(resp.SequenceNumber),
	}, nil
}

// runSQSReceive handles sqs.receive_message
func runSQSReceive(ctx context.Context, client *sqs.Client, action baseTypes.SQSAction) (interface{}, error) {

	input := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(action.SQS.QueueURL),
	}

	if action.SQS.MaxNumberOfMessages > 0 {
		input.MaxNumberOfMessages = action.SQS.MaxNumberOfMessages
	}
	if action.SQS.WaitTimeSeconds > 0 {
		input.WaitTimeSeconds = action.SQS.WaitTimeSeconds
	}
	if action.SQS.VisibilityTimeout > 0 {
		input.VisibilityTimeout = action.SQS.VisibilityTimeout
	}
	if len(action.SQS.AttributeNames) > 0 {
		input.AttributeNames = make([]types.QueueAttributeName, len(action.SQS.AttributeNames))
		for i, v := range action.SQS.AttributeNames {
			input.AttributeNames[i] = types.QueueAttributeName(v)
		}
	}
	if len(action.SQS.MessageAttributeNames) > 0 {
		input.MessageAttributeNames = action.SQS.MessageAttributeNames
	}

	resp, err := client.ReceiveMessage(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to receive sqs messages: %w", err)
	}

	type MessageOutput struct {
		MessageId         string                                      `json:"message_id"`
		ReceiptHandle     string                                      `json:"receipt_handle"`
		Body              string                                      `json:"body"`
		Attributes        map[string]string                           `json:"attributes"`
		MessageAttributes map[string]struct{ DataType, Value string } `json:"message_attributes"`
	}

	messages := make([]MessageOutput, 0, len(resp.Messages))
	for _, m := range resp.Messages {
		attributes := make(map[string]string)
		for k, v := range m.Attributes {
			attributes[string(k)] = v
		}

		msgAttrs := make(map[string]struct{ DataType, Value string })
		for k, v := range m.MessageAttributes {
			val := ""
			if v.StringValue != nil {
				val = *v.StringValue
			}
			msgAttrs[k] = struct{ DataType, Value string }{
				DataType: aws.ToString(v.DataType),
				Value:    val,
			}
		}

		messages = append(messages, MessageOutput{
			MessageId:         aws.ToString(m.MessageId),
			ReceiptHandle:     aws.ToString(m.ReceiptHandle),
			Body:              aws.ToString(m.Body),
			Attributes:        attributes,
			MessageAttributes: msgAttrs,
		})
	}

	return map[string]interface{}{
		"messages": messages,
	}, nil
}

// runSQSDelete handles sqs.delete_message
func runSQSDelete(ctx context.Context, client *sqs.Client, action baseTypes.SQSAction) (interface{}, error) {
	var _, err = client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(action.SQS.QueueURL),
		ReceiptHandle: aws.String(action.SQS.ReceiptHandle),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to delete sqs message: %w", err)
	}

	return map[string]string{"status": "success"}, nil
}

// runSQSPurge handles sqs.purge_queue
func runSQSPurge(ctx context.Context, client *sqs.Client, action baseTypes.SQSAction) (interface{}, error) {
	var _, err = client.PurgeQueue(ctx, &sqs.PurgeQueueInput{
		QueueUrl: aws.String(action.SQS.QueueURL),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to purge sqs queue: %w", err)
	}

	return map[string]string{"status": "success"}, nil
}

// runSQSMoveTask handles sqs.start_message_move_task
func runSQSMoveTask(ctx context.Context, client *sqs.Client, action baseTypes.SQSAction) (interface{}, error) {
	input := &sqs.StartMessageMoveTaskInput{
		SourceArn: aws.String(action.SQS.SourceARN),
	}

	if action.SQS.DestinationARN != "" {
		input.DestinationArn = aws.String(action.SQS.DestinationARN)
	}
	if action.SQS.MaxMessagesPerSecond > 0 {
		input.MaxNumberOfMessagesPerSecond = &action.SQS.MaxMessagesPerSecond
	}

	resp, err := client.StartMessageMoveTask(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to start message move task: %w", err)
	}

	return map[string]string{
		"task_handle": aws.ToString(resp.TaskHandle),
	}, nil
}
