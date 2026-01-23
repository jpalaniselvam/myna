---
title: Payload Handling
description: How to define request payloads for your actions.
---

Most AWS actions require data to be sent, whether it's the JSON body for a Lambda function or the message body for an SQS queue. `myna` standardizes how you define this data using the `[payload]` block.

## Configuration

The `[payload]` block supports two mutually exclusive fields: `data` and `file`.

| Field  | Type   | Description                                                                 |
| :----- | :----- | :-------------------------------------------------------------------------- |
| `data` | string | Inline content (string or JSON). Supports variable interpolation.           |
| `file` | string | Path to a file containing the payload. Relative to the action file.         |

> **Precedence**: If both are present, `file` takes precedence.

## Inline Data

Use the `data` field for simple, self-contained payloads. Multi-line strings in TOML (`"""`) make this easy for JSON.

```toml
[payload]
data = """
{
  "user_id": "{{USER_ID}}",
  "action": "create"
}
"""
```

### Variable Substitution

Inline data automatically supports Handlebars-style variable substitution (`{{VAR}}`).

```toml
[pre]
ENV = "dev"

[payload]
data = "Processing in environment: {{ENV}}"
```

## File-Based Payloads

For larger or binary payloads, use the `file` field.

```toml
[payload]
file = "./payloads/large-event.json"
```

### Resolution Logic

1.  **Relative Path**: If the path is relative (e.g., `./data.json`), it is resolved relative to the **Action file's directory**.
2.  **Absolute Path**: You can provide a full absolute path (not recommended for sharing code).
3.  **Variable Substitution**: Content loaded from files **also supports variable substitution**.

```json
/* payloads/large-event.json */
{
  "complex_object": {
    "id": "{{ID_FROM_VARS}}"
  }
}
```
