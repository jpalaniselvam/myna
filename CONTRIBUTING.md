# Contributing to Myna

First off, thanks for taking the time to contribute! `myna` is a community-driven project, and we welcome contributions of all kinds, from bug fixes and feature implementations to documentation improvements.

## How Can I Contribute?

### Reporting Bugs

- **Search existing issues** to see if the bug has already connect reported.
- **Open a new issue** if it hasn't. Provide a clear title, description, and steps to reproduce. Include your OS, Go version, and `myna` version.

### Suggesting Features

- Open an issue with the **enhancement** label.
- Explain the behavior you would like to see and why it would be useful.

### Contributing Code

1.  **Fork** the repository.
2.  **Clone** your fork locally.
3.  Create a **new branch** for your feature or bug fix:
    ```bash
    git checkout -b feature/my-new-feature
    ```
4.  **Make your changes**. Write clean, maintainable, and well-documented code.
5.  **Run tests** to ensure your changes didn't break anything:
    ```bash
    go test ./...
    ```
6.  **Format your code**:
    ```bash
    go fmt ./...
    ```
7.  **Commit** your changes with a descriptive commit message.
8.  **Push** to your branch:
    ```bash
    git push origin feature/my-new-feature
    ```
9.  Open a **Pull Request**. Describe your changes and link to any relevant issues.

## Development Guide

### Prerequisites

- [Go](https://go.dev/dl/) (Latest stable version recommended)
- [Git](https://git-scm.com/downloads)

### Setting Up the Environment

1.  Clone the repo:
    ```bash
    git clone https://github.com/jpalaniselvam/myna.git
    cd myna
    ```
2.  Install dependencies:
    ```bash
    go mod download
    ```
3.  Build the binary:
    ```bash
    go build -o myna cmd/myna/main.go
    ```

### Code Style

- Follow standard Go conventions (Effective Go).
- We use `go fmt` for formatting.

## License

By contributing, you agree that your contributions will be licensed under the Apache 2.0 License.
