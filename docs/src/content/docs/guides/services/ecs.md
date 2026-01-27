---
title: ECS
description: Managing Amazon ECS clusters and tasks.
---

`myna` supports operations for running tasks and managing services on Amazon ECS.

## Supported Kinds

*   `ecs.list_clusters`
*   `ecs.run_task`
*   `ecs.update_service`

## List Clusters (`ecs.list_clusters`)

Returns a list of existing clusters.

### Example

```toml
version = "1.0"
kind = "ecs.list_clusters"
```

## Run Task (`ecs.run_task`)

Starts a new task using the specified task definition.

### Configuration (`[ecs]`)

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `cluster` | string | No | The short name or full Amazon Resource Name (ARN) of the cluster on which to run your task. |
| `task_definition` | string | **Yes** | The family and revision (family:revision) or full ARN of the task definition to run. |
| `count` | int | No | The number of instantiations of the specified task to place on your cluster. Default is 1. |
| `launch_type` | string | No | The launch type on which to run your task (`EC2` or `FARGATE`). |
| `subnets` | list | No | (Fargate/Awsvpc) The IDs of the subnets associated with the task or service. |
| `security_groups` | list | No | (Fargate/Awsvpc) The IDs of the security groups associated with the task or service. |
| `assign_public_ip` | string | No | (Fargate/Awsvpc) Whether the task's elastic network interface receives a public IP address (`ENABLED` or `DISABLED`). |

### Example

```toml
version = "1.0"
kind = "ecs.run_task"

[ecs]
cluster = "my-cluster"
task_definition = "my-app:1"
launch_type = "FARGATE"
subnets = ["subnet-12345678", "subnet-87654321"]
security_groups = ["sg-12345678"]
assign_public_ip = "ENABLED"
```

## Update Service (`ecs.update_service`)

Modifies the parameters of a service.

### Configuration (`[ecs]`)

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `cluster` | string | No | The short name or full Amazon Resource Name (ARN) of the cluster your service to update is running on. |
| `service` | string | **Yes** | The name of the service to update. |
| `task_definition` | string | No | The family and revision (family:revision) or full ARN of the task definition to run in your service. |
| `desired_count` | int | No | The number of instantiations of the task to place and keep running in your service. |
| `force_new_deployment` | bool | No | Whether to force a new deployment of the service. |

### Example

```toml
version = "1.0"
kind = "ecs.update_service"

[ecs]
cluster = "my-cluster"
service = "my-service"
task_definition = "my-app:2"
force_new_deployment = true
```
