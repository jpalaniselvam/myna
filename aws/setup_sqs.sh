#!/bin/bash
set -e

QUEUE_NAME="myna-test-queue"
DLQ_NAME="myna-test-dlq"
FIFO_QUEUE_NAME="myna-test-fifo.fifo"
FIFO_DLQ_NAME="myna-test-fifo-dlq.fifo"

# Helper function to get queue URL
get_queue_url() {
    aws sqs get-queue-url --queue-name "$1" --output text 2>/dev/null || echo ""
}

# Helper function to get queue ARN
get_queue_arn() {
    local url="$1"
    if [ -z "$url" ]; then return; fi
    aws sqs get-queue-attributes --queue-url "$url" --attribute-names QueueArn --query Attributes.QueueArn --output text
}

echo "Starting SQS Setup..."

# 1. Standard DLQ
echo "Checking Standard DLQ: $DLQ_NAME"
DLQ_URL=$(get_queue_url "$DLQ_NAME")
if [ -z "$DLQ_URL" ]; then
    echo "Creating Standard DLQ..."
    DLQ_URL=$(aws sqs create-queue --queue-name "$DLQ_NAME" --output text)
    echo "Created DLQ: $DLQ_URL"
else
    echo "Standard DLQ exists."
fi
DLQ_ARN=$(get_queue_arn "$DLQ_URL")

# 2. Standard Queue
echo "Checking Standard Queue: $QUEUE_NAME"
QUEUE_URL=$(get_queue_url "$QUEUE_NAME")
if [ -z "$QUEUE_URL" ]; then
    echo "Creating Standard Queue..."
    
    # Construct attributes using valid JSON format
    # Escaping quotes for the inner JSON string of RedrivePolicy
    cat <<EOF > std-queue-attrs.json
{
  "RedrivePolicy": "{\"deadLetterTargetArn\":\"$DLQ_ARN\",\"maxReceiveCount\":\"3\"}"
}
EOF
    
    QUEUE_URL=$(aws sqs create-queue --queue-name "$QUEUE_NAME" --attributes file://std-queue-attrs.json --output text)
    rm std-queue-attrs.json
    
    echo "Created Standard Queue: $QUEUE_URL"
else
    echo "Standard Queue exists."
fi

# 3. FIFO DLQ
echo "Checking FIFO DLQ: $FIFO_DLQ_NAME"
FIFO_DLQ_URL=$(get_queue_url "$FIFO_DLQ_NAME")
if [ -z "$FIFO_DLQ_URL" ]; then
    echo "Creating FIFO DLQ..."
    
    cat <<EOF > fifo-dlq-attrs.json
{
  "FifoQueue": "true"
}
EOF
    
    FIFO_DLQ_URL=$(aws sqs create-queue --queue-name "$FIFO_DLQ_NAME" --attributes file://fifo-dlq-attrs.json --output text)
    rm fifo-dlq-attrs.json
    
    echo "Created FIFO DLQ: $FIFO_DLQ_URL"
else
    echo "FIFO DLQ exists."
fi
FIFO_DLQ_ARN=$(get_queue_arn "$FIFO_DLQ_URL")

# 4. FIFO Queue
echo "Checking FIFO Queue: $FIFO_QUEUE_NAME"
FIFO_QUEUE_URL=$(get_queue_url "$FIFO_QUEUE_NAME")
if [ -z "$FIFO_QUEUE_URL" ]; then
    echo "Creating FIFO Queue..."

    cat <<EOF > fifo-queue-attrs.json
{
  "FifoQueue": "true",
  "ContentBasedDeduplication": "true",
  "RedrivePolicy": "{\"deadLetterTargetArn\":\"$FIFO_DLQ_ARN\",\"maxReceiveCount\":\"3\"}"
}
EOF
    
    FIFO_QUEUE_URL=$(aws sqs create-queue --queue-name "$FIFO_QUEUE_NAME" --attributes file://fifo-queue-attrs.json --output text)
    rm fifo-queue-attrs.json
    
    echo "Created FIFO Queue: $FIFO_QUEUE_URL"
else
    echo "FIFO Queue exists."
fi

echo "SQS setup complete!"
