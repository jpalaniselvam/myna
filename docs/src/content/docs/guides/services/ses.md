---
title: SES
description: Sending emails with Amazon Simple Email Service (SES).
---

`myna` supports sending emails and managing identities in Amazon SES.

## Supported Kinds

*   `ses.send_email`
*   `ses.verify_email_identity`

## Send Email (`ses.send_email`)

Composes an email message and immediately queues it for sending.

### Configuration (`[ses]`)

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `source` | string | **Yes** | The email address that is sending the email. |
| `to_addresses` | list | **Yes** | The recipients' email addresses. |
| `cc_addresses` | list | No | The CC recipients. |
| `bcc_addresses` | list | No | The BCC recipients. |
| `subject` | string | **Yes** | The subject of the message. |
| `html_body` | string | No | The HTML body of the email. |
| `text_body` | string | No | The text body of the email. |

At least one of `html_body` or `text_body` must be specified.

### Example

```toml
version = "1.0"
kind = "ses.send_email"

[ses]
source = "sender@example.com"
to_addresses = ["recipient@example.com"]
subject = "Hello from Myna"
text_body = "This is a test email sent via Myna CLI."
html_body = "<h1>Hello from Myna</h1><p>This is a test email.</p>"
```

## Verify Email Identity (`ses.verify_email_identity`)

Verifies an email identity. This triggers a verification email to the specified address.

### Configuration (`[ses]`)

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `email_address` | string | **Yes** | The email address to verify. |

### Example

```toml
version = "1.0"
kind = "ses.verify_email_identity"

[ses]
email_address = "new-sender@example.com"
```
