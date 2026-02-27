#!/bin/bash
set -e

TABLE_USERS="myna-test-users"
TABLE_ORDERS="myna-test-orders"

echo "Starting DynamoDB Setup..."

# 1. Users Table
echo "Checking Table: $TABLE_USERS"
if ! aws dynamodb describe-table --table-name "$TABLE_USERS" >/dev/null 2>&1; then
    echo "Creating Users Table..."
    aws dynamodb create-table \
        --table-name "$TABLE_USERS" \
        --attribute-definitions \
            AttributeName=UserId,AttributeType=S \
        --key-schema \
            AttributeName=UserId,KeyType=HASH \
        --provisioned-throughput \
            ReadCapacityUnits=5,WriteCapacityUnits=5
    aws dynamodb wait table-exists --table-name "$TABLE_USERS"
    echo "Created Users Table."
else
    echo "Users Table exists."
fi

# 2. Orders Table
echo "Checking Table: $TABLE_ORDERS"
if ! aws dynamodb describe-table --table-name "$TABLE_ORDERS" >/dev/null 2>&1; then
    echo "Creating Orders Table..."
    aws dynamodb create-table \
        --table-name "$TABLE_ORDERS" \
        --attribute-definitions \
            AttributeName=UserId,AttributeType=S \
            AttributeName=OrderId,AttributeType=S \
        --key-schema \
            AttributeName=UserId,KeyType=HASH \
            AttributeName=OrderId,KeyType=RANGE \
        --provisioned-throughput \
            ReadCapacityUnits=5,WriteCapacityUnits=5
    aws dynamodb wait table-exists --table-name "$TABLE_ORDERS"
    echo "Created Orders Table."
else
    echo "Orders Table exists."
fi

echo "DynamoDB setup complete!"
