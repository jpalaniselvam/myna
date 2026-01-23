---
title: Design Overview
description: High-level architecture and concepts of myna.
---

This document outlines the overall design of the tool, including how AWS requests are managed and stored, as well as details on environment variables, templating, authentication, and CLI commands.

## Key Concepts

1. **Collections**: Group related actions together.
2. **Environments**: Define variables for different deployment stages (local, dev, prod).
3. **Actions**: Individual AWS serverless tasks (Lambda invoke, SQS send, etc.).
4. **Variables**: Dynamic value replacement within actions.
5. **Authentication**: Seamless AWS credential management.
6. **CLI Commands**: Powerful terminal interface for execution and management.

## Storage Format

Designed to be **Git-friendly**, the tool stores collections in a directory structure. All metadata and actions are stored in **TOML** files.

### Collections

Collections are the top-level entity.

- Represented as **directories** (the directory name is the collection name).
- Contain a `myna.toml` file with collection metadata.
- Can store a hierarchy of actions in subdirectories or directly within the collection directory.

### Environment Variables

Environment variables are managed using `.toml` configuration files:
All environment files are stored under environments directory of the collection root.

- `dev.toml`
- `local.toml`
- `prod.toml`

## Example Collection Structure

```text
my-collection
├── myna.toml
├── environments
│   ├── dev.toml
│   ├── local.toml
│   └── prod.toml
├── s3-upload.toml
├── lambda-invoke.toml
├── sqs-send-message.toml
├── sns-publish.toml
├── eventbridge-send-event.toml
├── coreapis
│   ├── create-lambda.toml
│   ├── update-lambda.toml
│   ├── delete-lambda.toml
│   └── invoke-lambda.toml
├── eventtesting
│   ├── action2.toml
│   └── action3.toml
└── cleanup
    └── action3.toml
```

#### Collection Metadata

All collection metadata is stored in a `myna.toml` file in the collection directory.
The collection metadata contains:

- version
- kind
- description
- vars
- aws

#### Example myna.toml

```toml
version = "1.0"
kind = "collection"
description = "My collection description"

[vars]
MY_VARIABLE = "my_value"
MY_OTHER_VARIABLE = "my_other_value"
```
