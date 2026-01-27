---
title: Installation
description: How to install myna on Windows, macOS, and Linux.
---

`myna` is distributed as a single binary, making it easy to install on any platform.

## Supported Platforms

Based on our build configuration, `myna` supports the following:

- **Windows**: x86_64, i386, ARM64
- **macOS (Darwin)**: x86_64, ARM64 (Apple Silicon)
- **Linux**: x86_64, i386, ARM64, ARMv6/v7

## Installation Methods

### Windows (winget) `coming soon`

The easiest way to install `myna` on Windows is via `winget`:

```powershell
winget install jpalaniselvam.myna
```

### Homebrew (macOS) `coming soon`

If you use Homebrew on macOS, you can install `myna` from our tap:

```bash
brew install jpalaniselvam/tap/myna
```

### Direct Download

For all other platforms, including **Linux** (where package managers like `apt` or `yum` are not yet supported), please download the binary directly:

1.  Go to the [GitHub Releases](https://github.com/jpalaniselvam/myna/releases) page.
2.  Download the archive for your OS and architecture (e.g., `myna_Linux_x86_64.tar.gz` or `myna_Windows_arm64.zip`).
3.  Extract the archive.
4.  Move the `myna` binary to a directory in your `PATH` (e.g., `/usr/local/bin` on Linux/macOS or a dedicated bin folder on Windows).

### Verification

After installation, verify that `myna` is working by running:

```bash
myna --version
```
