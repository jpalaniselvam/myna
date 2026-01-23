---
title: Actions
description: Defining AWS actions in TOML.
---

An **Action** is the fundamental unit of work in `myna`. It represents a single operation against an AWS service.

## Action Structure

All actions share a common structure defined by the `BaseAction` type.

```toml
version = "1.0"
kind = "<service>.<operation>"
description = "Description of what this action does"

# Optional Authentication Overrides
profile = "my-profile"
region = "us-east-1"
role_arn = "arn:aws:iam::..."

# Action-Specific Configuration
[<service>]
# Service specific fields...

# Pre-execution Variables
[pre]
user_id = "123"

# Request Payload (if applicable)
[payload]
data = """
{ ... }
"""
```

### Common Fields

| Field | Type | Description |
| :--- | :--- | :--- |
| `version` | string | Format version (e.g., "1.0") |
| `kind` | string | The type of action (e.g., `lambda.invoke`, `sqs.send_message`) |
| `description` | string | Human-readable description |
| `region` | string | (Optional) AWS Region override |
| `profile` | string | (Optional) AWS Named Profile override |
| `role_arn` | string | (Optional) IAM Role to assume |
| `timeout_ms` | int | (Optional) Client timeout in milliseconds |

### `[pre]` Block

Defines variables local to this action. These have the highest precedence and override Collection and Environment variables.

### `[payload]` Block

Defines the data to be sent with the request.

| Field | Description |
| :--- | :--- |
| `data` | Inline string data (often JSON). Supports variable interpolation `{{VAR}}`. |
| `file` | Path to a file containing the payload (relative to the action file). |

> **Note**: If both `data` and `file` are provided, `file` takes precedence (or behavior depends on implementation, check specific action docs).

## Supported Action Kinds

*   `lambda.invoke`: Invoke a Lambda function.
*   `sqs.send_message`: Send a message to an SQS queue.
*   `sqs.receive_message`: Receive messages from an SQS queue.
*   `sqs.delete_message`: Delete a message from an SQS queue.
*   `sqs.purge_queue`: Purge an SQS queue.
*   `sqs.start_message_move_task`: Start an SQS DLQ redrive task.
*   `sns.publish`: Publish a message to an SNS topic.
*   `eventbridge.put_events`: Put events to an EventBridge bus.
