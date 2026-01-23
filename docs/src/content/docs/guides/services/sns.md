---
title: SNS
description: Publishing to Amazon Simple Notification Service (SNS).
---

The `sns.publish` action allows you to send messages to an SNS topic, phone number, or mobile endpoint.

## Configuration

Set `kind = "sns.publish"` in your TOML file.

### `[sns]` Block

| Field | Type | Description |
| :--- | :--- | :--- |
| `topic_arn` | string | ARN of the topic to publish to. |
| `target_arn` | string | ARN of specific target (e.g. mobile endpoint). |
| `phone_number` | string | E.164 phone number for SMS. |
| `subject` | string | Subject line (for Email/Email-JSON). |
| `message_structure` | string | Set to `json` if sending different messages for different protocols. |
| `message_group_id` | string | Required for FIFO topics. |
| `message_deduplication_id` | string | Required for FIFO topics. |
| `message_attributes` | map | content metadata. |

> **Note**: You must provide exactly one of `topic_arn`, `target_arn`, or `phone_number`.

## Example

```toml
version = "1.0"
kind = "sns.publish"

[sns]
topic_arn = "arn:aws:sns:us-east-1:123456789012:my-topic"
subject = "System Alert"

[sns.message_attributes.Priority]
DataType = "String"
StringValue = "High"

[payload]
data = "Server load is high!"
```
