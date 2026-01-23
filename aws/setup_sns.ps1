$ErrorActionPreference = "Stop"

$TOPIC_NAME = "myna-test-topic"
$FIFO_TOPIC_NAME = "myna-test-topic.fifo"

# Helper function to get topic ARN
function Get-TopicArn ($name) {
    try {
        $arn = aws sns list-topics --query "Topics[?contains(TopicArn, '$name')].TopicArn" --output text 2>$null
        if ($LASTEXITCODE -eq 0) { return $arn }
        return $null
    } catch {
        return $null
    }
}

Write-Host "Starting SNS Setup..." -ForegroundColor Cyan

# 1. Standard Topic
Write-Host "Checking Standard Topic: $TOPIC_NAME"
$topicArn = Get-TopicArn $TOPIC_NAME
if (-not $topicArn) {
    Write-Host "Creating Standard Topic..."
    $topicArn = aws sns create-topic --name $TOPIC_NAME --output text
    Write-Host "Created Standard Topic: $topicArn" -ForegroundColor Green
} else {
    Write-Host "Standard Topic exists: $topicArn"
}

# 2. FIFO Topic
Write-Host "Checking FIFO Topic: $FIFO_TOPIC_NAME"
$fifoTopicArn = Get-TopicArn $FIFO_TOPIC_NAME
if (-not $fifoTopicArn) {
    Write-Host "Creating FIFO Topic..."
    
    $attributes = @{
        FifoTopic = "true"
        ContentBasedDeduplication = "true"
    } | ConvertTo-Json
    
    $attributes | Out-File "fifo-topic-attrs.json" -Encoding ascii
    
    $fifoTopicArn = aws sns create-topic --name $FIFO_TOPIC_NAME --attributes file://fifo-topic-attrs.json --output text
    Remove-Item "fifo-topic-attrs.json"
    
    Write-Host "Created FIFO Topic: $fifoTopicArn" -ForegroundColor Green
} else {
    Write-Host "FIFO Topic exists: $fifoTopicArn"
}

Write-Host "SNS setup complete!" -ForegroundColor Green
