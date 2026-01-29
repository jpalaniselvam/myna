#!/bin/bash
set -e

if [ -z "$1" ]; then
    echo "Error: Email identity argument is required."
    echo "Usage: $0 <email_identity>"
    exit 1
fi

EMAIL_IDENTITY="$1"

echo "Starting SES Cleanup..."
echo "Deleting email identity: $EMAIL_IDENTITY"

aws ses delete-identity --identity "$EMAIL_IDENTITY"

echo "SES cleanup complete!"
