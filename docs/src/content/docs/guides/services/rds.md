---
title: RDS
description: Managing Amazon RDS instances.
---

`myna` supports operations for managing the lifecycle of Amazon RDS DB instances.

## Supported Kinds

*   `rds.describe_db_instances`
*   `rds.start_db_instance`
*   `rds.stop_db_instance`
*   `rds.reboot_db_instance`

## Describe DB Instances (`rds.describe_db_instances`)

Returns information about provisioned RDS instances. This API supports pagination.

### Configuration (`[rds]`)

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `db_instance_identifier` | string | No | The user-supplied instance identifier. If this parameter is specified, information from only the specific DB instance is returned. |
| `filters` | map[string]list | No | A filter that specifies one or more DB instances to describe. |

### Example

```toml
version = "1.0"
kind = "rds.describe_db_instances"

[rds]
db_instance_identifier = "my-db-instance"
```

## Lifecycle Operations

Manage the power state of your DB instances.

### Supported Kinds
*   `rds.start_db_instance`: Starts an Amazon RDS DB instance that was stopped.
*   `rds.stop_db_instance`: Stops an Amazon RDS DB instance.
*   `rds.reboot_db_instance`: You might need to reboot your DB instance, usually for maintenance reasons.

### Configuration (`[rds]`)

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `db_instance_identifier` | string | **Yes** | The DB instance identifier. |
| `db_snapshot_identifier` | string | No | *(Stop only)* The user-supplied instance identifier of the DB Snapshot created immediately before the DB instance is stopped. |
| `force_failover` | bool | No | *(Reboot only)* When `true`, the reboot is conducted through a Multi-AZ failover. |

### Example

```toml
version = "1.0"
kind = "rds.reboot_db_instance"

[rds]
db_instance_identifier = "my-db-instance"
force_failover = false
```
