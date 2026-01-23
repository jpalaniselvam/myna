$ErrorActionPreference = "Stop"

$QUEUE_NAME = "myna-test-queue"
$DLQ_NAME = "myna-test-dlq"
$FIFO_QUEUE_NAME = "myna-test-fifo.fifo"
$FIFO_DLQ_NAME = "myna-test-fifo-dlq.fifo"

# Helper function to get queue URL
function Get-QueueUrl ($name) {
    try {
        $url = aws sqs get-queue-url --queue-name $name --output text 2>$null
        if ($LASTEXITCODE -eq 0) { return $url }
        return $null
    } catch {
        return $null
    }
}

# Helper function to get queue ARN
function Get-QueueArn ($url) {
    if (-not $url) { return $null }
    $arn = aws sqs get-queue-attributes --queue-url $url --attribute-names QueueArn --query Attributes.QueueArn --output text
    return $arn
}

Write-Host "Starting SQS Setup..." -ForegroundColor Cyan

# 1. Standard DLQ
Write-Host "Checking Standard DLQ: $DLQ_NAME"
$dlqUrl = Get-QueueUrl $DLQ_NAME
if (-not $dlqUrl) {
    Write-Host "Creating Standard DLQ..."
    $dlqUrl = aws sqs create-queue --queue-name $DLQ_NAME --output text
    Write-Host "Created DLQ: $dlqUrl" -ForegroundColor Green
} else {
    Write-Host "Standard DLQ exists."
}
$dlqArn = Get-QueueArn $dlqUrl

# 2. Standard Queue
Write-Host "Checking Standard Queue: $QUEUE_NAME"
$queueUrl = Get-QueueUrl $QUEUE_NAME
if (-not $queueUrl) {
    Write-Host "Creating Standard Queue..."
    
    # Construct attributes
    $redrivePolicy = @{
        deadLetterTargetArn = $dlqArn
        maxReceiveCount = "3"
    } | ConvertTo-Json -Compress

    $attributes = @{
        RedrivePolicy = $redrivePolicy
    } | ConvertTo-Json
    
    $attributes | Out-File "std-queue-attrs.json" -Encoding ascii
    
    $queueUrl = aws sqs create-queue --queue-name $QUEUE_NAME --attributes file://std-queue-attrs.json --output text
    Remove-Item "std-queue-attrs.json"
    
    Write-Host "Created Standard Queue: $queueUrl" -ForegroundColor Green
} else {
    Write-Host "Standard Queue exists."
}

# 3. FIFO DLQ
Write-Host "Checking FIFO DLQ: $FIFO_DLQ_NAME"
$fifoDlqUrl = Get-QueueUrl $FIFO_DLQ_NAME
if (-not $fifoDlqUrl) {
    Write-Host "Creating FIFO DLQ..."
    
    $attributes = @{
        FifoQueue = "true"
    } | ConvertTo-Json

    $attributes | Out-File "fifo-dlq-attrs.json" -Encoding ascii
    
    $fifoDlqUrl = aws sqs create-queue --queue-name $FIFO_DLQ_NAME --attributes file://fifo-dlq-attrs.json --output text
    Remove-Item "fifo-dlq-attrs.json"
    
    Write-Host "Created FIFO DLQ: $fifoDlqUrl" -ForegroundColor Green
} else {
    Write-Host "FIFO DLQ exists."
}
$fifoDlqArn = Get-QueueArn $fifoDlqUrl

# 4. FIFO Queue
Write-Host "Checking FIFO Queue: $FIFO_QUEUE_NAME"
$fifoQueueUrl = Get-QueueUrl $FIFO_QUEUE_NAME
if (-not $fifoQueueUrl) {
    Write-Host "Creating FIFO Queue..."

    $redrivePolicy = @{
        deadLetterTargetArn = $fifoDlqArn
        maxReceiveCount = "3"
    } | ConvertTo-Json -Compress

    $attributes = @{
        FifoQueue = "true"
        ContentBasedDeduplication = "true"
        RedrivePolicy = $redrivePolicy
    } | ConvertTo-Json
    
    $attributes | Out-File "fifo-queue-attrs.json" -Encoding ascii
    
    $fifoQueueUrl = aws sqs create-queue --queue-name $FIFO_QUEUE_NAME --attributes file://fifo-queue-attrs.json --output text
    Remove-Item "fifo-queue-attrs.json"
    
    Write-Host "Created FIFO Queue: $fifoQueueUrl" -ForegroundColor Green
} else {
    Write-Host "FIFO Queue exists."
}

Write-Host "SQS setup complete!" -ForegroundColor Green
