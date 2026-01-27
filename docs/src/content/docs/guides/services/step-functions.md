---
title: Step Functions
description: Orchestrating workflows with AWS Step Functions.
---

`myna` supports operations for managing and executing AWS Step Functions state machines.

## Supported Kinds

*   `sfn.list_state_machines`
*   `sfn.start_execution`
*   `sfn.describe_execution`
*   `sfn.stop_execution`

## List State Machines (`sfn.list_state_machines`)

Lists the existing state machines.

### Examples

```toml
version = "1.0"
kind = "sfn.list_state_machines"
```

## Start Execution (`sfn.start_execution`)

Starts a state machine execution.

### Configuration (`[sfn]`)

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `state_machine_arn` | string | **Yes** | The Amazon Resource Name (ARN) of the state machine to execute. |
| `name` | string | No | The name of the execution. |

The input for the execution is taken from the `[payload]` block (either `data` or `file`).

### Example

```toml
version = "1.0"
kind = "sfn.start_execution"

[sfn]
state_machine_arn = "arn:aws:states:us-east-1:123456789012:stateMachine:MyStateMachine"
name = "test-execution-1"

[payload]
data = """
{
  "input_key": "input_value"
}
"""
```

## Describe Execution (`sfn.describe_execution`)

Describes an execution.

### Configuration (`[sfn]`)

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `execution_arn` | string | **Yes** | The Amazon Resource Name (ARN) of the execution to describe. |

### Example

```toml
version = "1.0"
kind = "sfn.describe_execution"

[sfn]
execution_arn = "arn:aws:states:us-east-1:123456789012:execution:MyStateMachine:test-execution-1"
```

## Stop Execution (`sfn.stop_execution`)

Stops an execution.

### Configuration (`[sfn]`)

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `execution_arn` | string | **Yes** | The Amazon Resource Name (ARN) of the execution to stop. |
| `error` | string | No | The error code of the failure. |
| `cause` | string | No | A more detailed explanation of the cause of the failure. |

### Example

```toml
version = "1.0"
kind = "sfn.stop_execution"

[sfn]
execution_arn = "arn:aws:states:us-east-1:123456789012:execution:MyStateMachine:test-execution-1"
cause = "Timeout"
```
