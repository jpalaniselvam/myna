---
title: Collections
description: Organizing actions into collections.
---

In `myna`, a **Collection** is a directory that groups related actions together. It provides a shared context for variables and configuration.

## Anatomy of a Collection

A collection is defined by the presence of a `myna.toml` file in a directory.

```text
my-collection/
├── myna.toml          # Collection definition
├── environments/      # Environment configurations
│   ├── dev.toml
│   └── prod.toml
├── create-user.toml   # Action
└── get-user.toml      # Action
```

## Collection Metadata (`myna.toml`)

The `myna.toml` file contains metadata about the collection and shared variables.

```toml
# myna.toml
version = "1.0"
kind = "collection"
description = "User Management API"

[credentials]
region = "us-east-1"
profile = "dev-profile"

[pre]
# Variables available to all actions in this collection
API_VERSION = "v1"
DEFAULT_TIMEOUT = 5000
TABLE_NAME = "users-${ENV}"
```

### Fields

| Field | Description |
| :--- | :--- |
| `version` | The version of the collection format (e.g., "1.0") |
| `kind` | Must be `"collection"` |
| `description` | A human-readable description |
| `[credentials]` | (Optional) Shared AWS credentials (region, profile) |
| `[pre]` | A table of shared variables (strings, numbers, booleans) |

## Variable Inheritance

Variables defined in `[pre]` of `myna.toml` are available to all actions within the collection directory (and subdirectories). They can be overridden by Environment variables or Action-specific variables.

See [Variable Resolution](/guides/variable-resolution) for more details.
