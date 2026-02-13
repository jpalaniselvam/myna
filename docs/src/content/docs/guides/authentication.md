---
title: Authentication
description: How to handle AWS authentication in myna.
---

`myna` relies on the standard AWS SDK credential chain. This means it supports the same authentication methods as the AWS CLI.

## Core Principles

1.  **No Credentials in Files**: `myna` does **not** store access keys or secrets in its TOML files. It only references profile names or role ARNs within a `[credentials]` block.
2.  **Standard Configuration**: It uses the standard `~/.aws/credentials` and `~/.aws/config` files.
3.  **Flexibility**: Different actions can run as different users/roles.

## Supported Methods

`myna` supports all authentication methods provided by the standard AWS credential chain, including:

1.  **Environment Variables**: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, etc.
2.  **Shared Configuration Files**: Profiles setup via `aws configure` or `aws configure sso`.
3.  **IAM Roles**: For EC2/ECS/Lambda instance profiles.

## Configuring Authentication

You can configure authentication at the **Action** level.

### Using a Named Profile

To specify which AWS profile to use, set the `profile` field in your action's metadata:

```toml
# my-action.toml
version = "1.0"
kind = "lambda.invoke"
description = "Invoke Lambda using prod profile"

[credentials]
profile = "prod-user"

[lambda]
function_name = "my-function"
```

### Assuming a Role

You can also assume a specific IAM role for an action using `role_arn`. This happens *after* the initial credentials (from default or named profile) are loaded.

```toml
# deploy-action.toml
version = "1.0"
kind = "lambda.invoke"
description = "Invoke using assumed role"

[credentials]
profile = "ci-user"
role_arn = "arn:aws:iam::123456789012:role/DeployRole"

[lambda]
function_name = "deploy-function"
```

## SSO Support

`myna` works seamlessly with AWS SSO (Identity Center).

1.  Login via the AWS CLI: `aws sso login --profile my-sso-profile`
2.  Reference the profile in your action:
    ```toml
    [credentials]
    profile = "my-sso-profile"
    ```

If your session has expired, `myna` will prompt you to login again.
