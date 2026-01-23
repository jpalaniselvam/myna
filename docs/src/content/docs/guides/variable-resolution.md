---
title: Variable Resolution
description: Understanding how myna resolves variables.
---

`myna` uses a powerful variable substitution system allowing you to parameterize your actions.

## Resolution Hierarchy

When a variable like `{{MY_VAR}}` is encountered, `myna` looks for its value in the following order (Highest to Lowest Priority):

1.  **Action `[pre]`**: Variables defined in the action file itself.
2.  **Collection `[pre]`**: Variables defined in `myna.toml`.
3.  **Environment `[pre]`**: Variables defined in `environments/<env>.toml`.

### Example

**Action (`action.toml`)**:
```toml
[pre]
TIMEOUT = 500
```

**Collection (`myna.toml`)**:
```toml
[pre]
TIMEOUT = 1000
REGION = "us-east-1"
```

**Environment (`dev.toml`)**:
```toml
[pre]
REGION = "eu-west-1"
STAGE = "dev"
```

**Resulting Values**:
*   `TIMEOUT`: **500** (Action overrides Collection)
*   `REGION`: **"us-east-1"** (Collection overrides Environment)
*   `STAGE`: **"dev"** (From Environment)

## Syntax

### Handlebars `{{VAR}}`
Used inside string values (like `payload.data` or `function_name`). This is a direct text replacement.

```toml
[payload]
data = """
{
  "id": "{{USER_ID}}"
}
"""
```

### Shell Style `${VAR}`
Used mostly when defining variables in TOML tables to reference *other* variables available at that time/level.

## Missing Variables

If a variable is referenced but not found in any of the 3 scopes, `myna` will:
1.  Log a warning to the console.
2.  Leave the template string as-is (e.g., `{{MISSING_VAR}}`).
