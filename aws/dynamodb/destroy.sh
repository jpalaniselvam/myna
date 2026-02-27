#!/bin/bash
set -e

TABLE_USERS="myna-test-users"
TABLE_ORDERS="myna-test-orders"

echo "Destroying DynamoDB resources..."

aws dynamodb delete-table --table-name "$TABLE_USERS" 2>/dev/null || true
aws dynamodb delete-table --table-name "$TABLE_ORDERS" 2>/dev/null || true

echo "DynamoDB resources destroyed."
