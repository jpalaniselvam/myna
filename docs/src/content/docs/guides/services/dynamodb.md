---
title: DynamoDB
description: interacting with Amazon DynamoDB.
---

`myna` supports common DynamoDB operations for managing data.

## Supported Kinds

*   `dynamodb.list_tables`
*   `dynamodb.put_item`
*   `dynamodb.get_item`
*   `dynamodb.update_item`
*   `dynamodb.delete_item`
*   `dynamodb.query`
*   `dynamodb.scan`

## List Tables (`dynamodb.list_tables`)

Returns an array of table names associated with the current account and endpoint.

### Example

```toml
version = "1.0"
kind = "dynamodb.list_tables"
```

## Data Operations

### Put Item (`dynamodb.put_item`)

Creates a new item, or replaces an old item with a new item.

#### Configuration (`[dynamodb]`)

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `table_name` | string | **Yes** | The name of the table to contain the item. |
| `condition_expression` | string | No | A condition that must be satisfied in order for a conditional PutItem operation to succeed. |

The item to put is taken from the `[payload]` block (either `data` or `file`). The payload should be a standard JSON object. The tool handles the conversion to DynamoDB JSON format.

#### Example

```toml
version = "1.0"
kind = "dynamodb.put_item"

[dynamodb]
table_name = "Users"

[payload]
data = """
{
  "UserId": "user-123",
  "Name": "Alice",
  "Age": 30,
  "Active": true
}
"""
```

### Get Item (`dynamodb.get_item`)

The `GetItem` operation returns a set of attributes for the item with the given primary key.

#### Configuration (`[dynamodb]`)

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `table_name` | string | **Yes** | The name of the table containing the requested item. |
| `key` | map | **Yes** | A map representing the primary key of the item to retrieve. |
| `consistent_read` | bool | No | Determines the read consistency model: If set to `true`, then the operation uses strongly consistent reads; otherwise, the operation uses eventually consistent reads. |

#### Example

```toml
version = "1.0"
kind = "dynamodb.get_item"

[dynamodb]
table_name = "Users"
key = { "UserId" = "user-123" }
```

### Delete Item (`dynamodb.delete_item`)

Deletes a single item in a table by primary key.

#### Configuration (`[dynamodb]`)

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `table_name` | string | **Yes** | The name of the table from which to delete the item. |
| `key` | map | **Yes** | A map representing the primary key of the item to delete. |
| `condition_expression` | string | No | A condition that must be satisfied in order for a conditional DeleteItem to succeed. |

#### Example

```toml
version = "1.0"
kind = "dynamodb.delete_item"

[dynamodb]
table_name = "Users"
key = { "UserId" = "user-123" }
```

### Update Item (`dynamodb.update_item`)

Edits an existing item's attributes, or adds a new item to the table if it does not already exist.

#### Configuration (`[dynamodb]`)

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `table_name` | string | **Yes** | The name of the table. |
| `key` | map | **Yes** | The primary key of the item to be updated. |
| `update_expression` | string | No | An expression that defines one or more attributes to be updated, action to be performed, and new value(s). |
| `condition_expression` | string | No | A condition that must be satisfied in order for a conditional update to succeed. |
| `expression_attribute_names` | map | No | Substitution tokens for attribute names in an expression. |
| `expression_attribute_values` | map | No | Substitution tokens for attribute values in an expression. |

#### Example

```toml
version = "1.0"
kind = "dynamodb.update_item"

[dynamodb]
table_name = "Users"
key = { "UserId" = "user-123" }
update_expression = "SET Age = :newAge"
expression_attribute_values = { ":newAge" = 31 }
```

### Query (`dynamodb.query`)

The `Query` operation finds items based on primary key values.

#### Configuration (`[dynamodb]`)

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `table_name` | string | **Yes** | The name of the table. |
| `key_condition_expression` | string | **Yes** | The condition that specifies the key values for items to be retrieved. |
| `expression_attribute_values` | map | No | Substitution tokens for attribute values. |
| `index_name` | string | No | The name of an index to query. |
| `limit` | int | No | The maximum number of items to evaluate (not necessarily the number of matching items). |

#### Example

```toml
version = "1.0"
kind = "dynamodb.query"

[dynamodb]
table_name = "Orders"
key_condition_expression = "OrderId = :oid"
expression_attribute_values = { ":oid" = "ord-555" }
```

### Scan (`dynamodb.scan`)

The `Scan` operation returns one or more items and item attributes by accessing every item in a table or a secondary index.

#### Configuration (`[dynamodb]`)

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `table_name` | string | **Yes** | The name of the table. |
| `filter_expression` | string | No | A string that contains conditions that DynamoDB applies after the Scan operation, but before the data is returned to you. |
| `expression_attribute_values` | map | No | Substitution tokens for attribute values. |
| `limit` | int | No | The maximum number of items to evaluate. |

#### Example

```toml
version = "1.0"
kind = "dynamodb.scan"

[dynamodb]
table_name = "Users"
limit = 10
```
