---
title: EventBridge
description: Putting events to Amazon EventBridge.
---

The `eventbridge.put_events` action allows you to ingest events into an EventBridge event bus.

## Configuration

Set `kind = "eventbridge.put_events"` in your TOML file.

### `[eventbridge]` Block

| Field | Type | Description |
| :--- | :--- | :--- |
| `bus_name` | string | Name or ARN of the event bus (default: `default`). |
| `source` | string | Identify the service/app sending the event. |
| `detail_type` | string | Free-form string identifying the type of event. |
| `resources` | list | AWS ARNs involved in the event. |
| `trace_header` | string | X-Ray trace header. |

## Example

```toml
version = "1.0"
kind = "eventbridge.put_events"

[eventbridge]
bus_name = "my-custom-bus"
source = "com.mycompany.app"
detail_type = "OrderPlaced"
resources = ["arn:aws:ec2:us-east-1:123:instance/i-12345"]

[payload]
data = """
{
  "order_id": "12345",
  "amount": 99.99,
  "currency": "USD"
}
"""
```

## Output

The action returns the `EventId` for each submitted event, or an error code if submission failed.
