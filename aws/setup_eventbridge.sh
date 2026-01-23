#!/bin/bash
set -e

BUS_NAME="myna-test-bus"

# Helper function to get event bus ARN
get_event_bus_arn() {
    aws events describe-event-bus --name "$1" --query "Arn" --output text 2>/dev/null || echo ""
}

echo "Starting EventBridge Setup..."

# 1. Custom Event Bus
echo "Checking Event Bus: $BUS_NAME"
BUS_ARN=$(get_event_bus_arn "$BUS_NAME")
if [ -z "$BUS_ARN" ]; then
    echo "Creating Event Bus..."
    BUS_ARN=$(aws events create-event-bus --name "$BUS_NAME" --query "EventBusArn" --output text)
    echo "Created Event Bus: $BUS_ARN"
else
    echo "Event Bus exists."
fi

echo "EventBridge setup complete!"
