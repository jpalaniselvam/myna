---
title: Why myna?
description: The motivation behind building myna - the Postman for AWS serverless.
---

As a user of the AWS Console, I found it difficult to switch between different services and regions to perform basic operations such as invoking a Lambda function, putting data into S3, or sending messages to SQS, SNS, or EventBridge for testing.

I wanted a tool like **Postman** but for AWS services, where I can perform all these operations in a single place. For example, I want to configure a Lambda request with different input payloads (or no payload), run it, and see the output. I also want to be able to configure and run requests for S3, SQS, SNS, and EventBridge in the same way.

The goal is to do this in a **repeatable** way, just like HTTP requests in Postman.

## The Problem with Current Tools

Modern serverless systems are built from events, queues, and functions, not just HTTP endpoints. However, most developer tools are still API-first, requiring you to:
- Navigate complex AWS Console UIs for simple tests.
- Write one-off scripts just to send a message to a queue.
- Manually copy-paste payloads across different service tabs.

## The myna Solution

`myna` bridges this gap by providing:
- **Git-first workflows**: Keep your test cases in version control next to your code.
- **Repeatable execution**: Save and replay events across different environments.
- **Unified Interface**: Use a consistent TOML-based configuration for Lambda, S3, SQS, SNS, and EventBridge.
