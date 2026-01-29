#!/bin/bash
set -e

if [ -z "$1" ]; then
    echo "Error: Email identity argument is required."
    echo "Usage: $0 <email_identity>"
    exit 1
fi

EMAIL_IDENTITY="$1"

echo "Starting SES Setup..."
echo "Using email identity: $EMAIL_IDENTITY"

# Check if identity exists
if aws ses list-identities --identity-type EmailAddress --query "Identities" --output text | grep -q "$EMAIL_IDENTITY"; then
    echo "Identity $EMAIL_IDENTITY already exists."
else
    echo "Creating/Verifying Email Identity: $EMAIL_IDENTITY"
    aws ses verify-email-identity --email-address "$EMAIL_IDENTITY"
    echo "Verification email sent to $EMAIL_IDENTITY. Please verify it to send emails in Sandbox mode."
fi

echo "SES setup complete!"
