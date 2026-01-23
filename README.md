## Why?

As a user of the AWS Console, I found it difficult to switch between different services and regions to perform basic operations such as invoking a Lambda function, putting data into S3, or sending messages to SQS, SNS, or EventBridge for testing.

I wanted a tool like **Postman** but for AWS services, where I can perform all these operations in a single place. For example, I want to configure a Lambda request with different input payloads (or no payload), run it, and see the output. I also want to be able to configure and run requests for S3, SQS, SNS, and EventBridge in the same way.

The goal is to do this in a **repeatable** way, just like HTTP requests in Postman.

# myna

A Git-first CLI for executing and replaying AWS serverless actions.

`myna` lets you invoke Lambda functions, send SQS messages, publish SNS events, emit EventBridge events, upload to S3, and run serverless workflows using simple, version-controlled TOML files.

Think **Postman for AWS serverless**

Modern serverless systems are built from events, queues, and functions, not HTTP endpoints.

But most developer tools are still API-first.

`myna` is designed for:

* Executing AWS serverless primitives directly
* Replaying real events deterministically
* Chaining actions without writing glue code
* Keeping everything in Git, not a database
* Working locally, in CI, or in production environments

## Important Design Principles

1. **Configurable**: Collections can be configured with variables, environments, and actions.
2. **Repeatable**: Collections can be run multiple times with different variables and environments.
3. **Git-Friendly**: Collections are stored in a way that is easy to version control.
4. **CLI-First Design**: The tool is designed to be used from the command line

## Non-Goals

* Abstracting AWS APIs
* Replacing infrastructure as code tools
* Hiding request or response payloads
* Providing a hosted service

`myna` is a developer tool, not a platform.

## Contributing

Contributions are welcome.

* Issues for bugs and feature requests
* PRs for new action kinds and improvements
* Discussions for design decisions

## License

Apache 2.0 License
