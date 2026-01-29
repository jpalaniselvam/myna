---
title: Environments
description: Managing environment-specific configurations.
---

Environments allow you to define configuration sets for different stages of your application lifecycle (e.g., `local`, `dev`, `prod`).

## Defining Environments

Environments are defined as TOML files within the `environments/` directory of your collection.

```text
my-collection/
├── myna.toml
└── environments/
    ├── dev.toml
    ├── prod.toml
    └── local.toml
```

## Environment File Structure

Each environment file is a simple TOML file that defines a `[pre]` table. These variables are loaded when you specify the environment flag.

```toml
# environments/dev.toml
[pre]
ENV = "dev"
AWS_REGION = "us-east-1"
API_URL = "https://api-dev.example.com"
```

## Using Environments

To run an action using a specific environment's configuration, use the `--env` flag with the `run` command:

```bash
myna run ./actions/invoke.toml --env dev
```

This command will:
1.  Load variables from `environments/dev.toml`.
2.  Merge them with variables from `myna.toml` (Collection).
3.  Merge them with variables from the action file.

## Best Practices

*   **Keep it minimal**: Use environment files for values that *change* between environments (URLs, Account IDs, Table Names).
*   **Use Placeholders**: Define structured naming conventions in your Collection `[pre]` (e.g., `TABLE_NAME = "users-${ENV}"`) and just define `ENV` in your environment files.
