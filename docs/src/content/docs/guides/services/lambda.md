---
title: Lambda
description: Invoking AWS Lambda functions.
---

The `lambda.invoke` action allows you to execute an AWS Lambda function.

## Configuration

Set `kind = "lambda.invoke"` in your TOML file.

### `[lambda]` Block

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `function_name` | string | **Yes** | The name or ARN of the function to invoke. |
| `invocation_type` | string | No | `RequestResponse` (default), `Event`, or `DryRun`. |
| `qualifier` | string | No | Version or alias to invoke (e.g., `prod`, `1`). |
| `client_context` | map | No | Up to 3583 bytes of base64-encoded JSON context data. |

## Example

```toml
version = "1.0"
kind = "lambda.invoke"
description = "Sync invocation of my-function"

[lambda]
function_name = "my-function"
invocation_type = "RequestResponse"
qualifier = "prod"

[client_context]
custom_app_id = "myna-cli"

[payload]
data = """
{
  "key1": "value1",
  "key2": "{{dynamic_value}}"
}
"""
```

## Output

The action outputs the function's response payload, status code, and executed version.
