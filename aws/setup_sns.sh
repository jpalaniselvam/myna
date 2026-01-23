#!/bin/bash
set -e

TOPIC_NAME="myna-test-topic"
FIFO_TOPIC_NAME="myna-test-topic.fifo"

# Helper function to get topic ARN
get_topic_arn() {
    aws sns list-topics --query "Topics[?contains(TopicArn, '$1')].TopicArn" --output text 2>/dev/null
}

echo "Starting SNS Setup..."

# 1. Standard Topic
echo "Checking Standard Topic: $TOPIC_NAME"
TOPIC_ARN=$(get_topic_arn "$TOPIC_NAME")
if [ -z "$TOPIC_ARN" ]; then
    echo "Creating Standard Topic..."
    TOPIC_ARN=$(aws sns create-topic --name "$TOPIC_NAME" --output text)
    echo "Created Standard Topic: $TOPIC_ARN"
else
    echo "Standard Topic exists: $TOPIC_ARN"
fi

# 2. FIFO Topic
echo "Checking FIFO Topic: $FIFO_TOPIC_NAME"
FIFO_TOPIC_ARN=$(get_topic_arn "$FIFO_TOPIC_NAME")
if [ -z "$FIFO_TOPIC_ARN" ]; then
    echo "Creating FIFO Topic..."
    
    cat <<EOF > fifo-topic-attrs.json
{
  "FifoTopic": "true",
  "ContentBasedDeduplication": "true"
}
EOF
    
    FIFO_TOPIC_ARN=$(aws sns create-topic --name "$FIFO_TOPIC_NAME" --attributes file://fifo-topic-attrs.json --output text)
    rm fifo-topic-attrs.json
    
    echo "Created FIFO Topic: $FIFO_TOPIC_ARN"
else
    echo "FIFO Topic exists: $FIFO_TOPIC_ARN"
fi

echo "SNS setup complete!"
