---
title: EC2
description: Managing Amazon EC2 instances.
---

`myna` supports common EC2 instance lifecycle operations and information retrieval.

## Supported Kinds

*   `ec2.describe_instances`
*   `ec2.start_instances`
*   `ec2.stop_instances`
*   `ec2.reboot_instances`
*   `ec2.terminate_instances`

## Describe Instances (`ec2.describe_instances`)

Retrieves detailed information about one or more instances.

### Configuration (`[ec2]`)

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `instance_ids` | list | No | A list of instance IDs to describe. |
| `filters` | map[string]list | No | Filters to scope the results. Keys are filter names (e.g., `instance-state-name`), values are lists of allowed strings. |

### Example

```toml
version = "1.0"
kind = "ec2.describe_instances"

[ec2]
filters = { "instance-state-name" = ["running", "pending"], "tag:Environment" = ["production"] }
```

## Lifecycle Operations

Manage the power state of your instances.

### Supported Kinds
*   `ec2.start_instances`: Starts an instance that uses an Amazon EBS volume as its root device.
*   `ec2.stop_instances`: Stops an Amazon EBS-backed instance.
*   `ec2.reboot_instances`: Requests a reboot of the specified instances.
*   `ec2.terminate_instances`: Shuts down the specified instances. This operation is idempotent.

### Configuration (`[ec2]`)

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `instance_ids` | list | **Yes** | The IDs of the instances. |
| `hibernate` | bool | No | *(Stop only)* Hibernates the instance if the instance was enabled for hibernation at launch. |
| `force` | bool | No | *(Stop only)* Forces the instances to stop. The instances do not have an opportunity to flush file system caches or file system metadata. |
| `dry_run` | bool | No | Checks whether you have the required permissions for the action, without actually making the request. |

### Example

```toml
version = "1.0"
kind = "ec2.stop_instances"

[ec2]
instance_ids = ["i-1234567890abcdef0", "i-0987654321fedcba0"]
hibernate = true
```
