---
title: Basic Usage
description: Learn how to run your first AWS serverless action with myna.
---

`myna` is built around the concept of **Actions**. An action is a simple TOML file that describes an AWS operation, like invoking a Lambda function or sending an SQS message.

## The `run` Command

The primary way to use `myna` is the `run` command. It executes a single action file.

## Your First Action File

Actions are defined in TOML files. Here is a simple example of an action that invokes a Lambda function with a JSON payload:

```toml
# my-action.toml
version = "1.0"
kind = "lambda.invoke"
description = "Invoke a simple JSON handler"

[lambda]
function_name = "my-lambda-function"

[pre]
user_id = "user_123"

[payload]
data = """
{
  "user_id": "{{user_id}}",
  "status": "active"
}
"""
```

### Simple Execution

To run the action defined above:

```bash
myna run ./my-action.toml
```

## Execution Contexts

`myna` automatically detects the context in which an action is being run.

### 1. Standalone Context
If you run an action file that is not part of a collection, `myna` runs it with only the variables defined inside that file. This is perfect for quick tests.

### 2. Collection Context
If `myna` finds a `myna.toml` file in the current directory or any parent directory, it treats the action as part of a **Collection**. This allows you to:
- Share variables across multiple actions.
- Use environment-specific configurations.

## Using Environments

When working within a collection, you can switch between different environments (like `dev` or `prod`) using the `--env` flag:

```bash
myna run ./actions/invoke-lambda.toml --env dev
```

This will load variables from `environments/dev.toml` relative to your collection root.

## Variable Resolution

`myna` resolves variables in a specific order of priority:
1.  **Action Variables**: Defined in the `[pre]` block of your action TOML. (Highest priority)
2.  **Collection Variables**: Defined in `myna.toml`.
3.  **Environment Variables**: Defined in files like `environments/dev.toml`. (Lowest priority)

This hierarchy allows you to define defaults at the collection level and override them for specific actions or environments.

## Next Steps

- Explore [Design Overview](/myna/guides/design-overview/) to learn more about the project structure.
- Check out the [Action Guide](/myna/guides/actions/) to learn more about configuring actions.

