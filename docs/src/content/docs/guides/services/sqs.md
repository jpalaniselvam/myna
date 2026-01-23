---
title: SQS
description: Interacting with Amazon Simple Queue Service (SQS).
---

`myna` supports multiple SQS operations, including sending, receiving, and managing messages.

## Supported Kinds

*   `sqs.send_message`
*   `sqs.receive_message`
*   `sqs.delete_message`
*   `sqs.purge_queue`
*   `sqs.start_message_move_task`

## Send Message (`sqs.send_message`)

Sends a message to a queue.

### Configuration (`[sqs]`)

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `queue_url` | string | **Yes** | The URL of the SQS queue. |
| `delay_seconds` | int | No | Delivery delay (0-900 seconds). |
| `message_group_id` | string | No | Required for FIFO queues. |
| `message_deduplication_id` | string | No | Required for FIFO queues (unless content-based dedup is on). |
| `message_attributes` | map | No | Custom metadata/attributes. |

### Example

```toml
version = "1.0"
kind = "sqs.send_message"

[sqs]
queue_url = "https://sqs.us-east-1.amazonaws.com/12345/my-queue"
delay_seconds = 10

[sqs.message_attributes.MyAttr]
DataType = "String"
StringValue = "MyValue"

[payload]
data = "Hello World"
```

## Receive Message (`sqs.receive_message`)

Polls a queue for messages.

### Configuration (`[sqs]`)

| Field | Type | Description |
| :--- | :--- | :--- |
| `max_number_of_messages` | int | Max messages to retrieve (1-10). |
| `wait_time_seconds` | int | Long polling wait time (0-20). |
| `visibility_timeout` | int | Time to keep messages hidden after receive. |
| `attribute_names` | list | System attributes to retrieve (e.g., `["All"]`). |
| `message_attribute_names` | list | Message attributes to retrieve. |

## Other Operations

### Purge Queue (`sqs.purge_queue`)
Deletes all messages in a queue. Requires `queue_url`.

### Delete Message (`sqs.delete_message`)
Deletes a specific message. Requires `queue_url` and `receipt_handle`.

### Start Redrive (`sqs.start_message_move_task`)
Moves messages from a DLQ back to source. Requires `source_arn` and optional `destination_arn`.
