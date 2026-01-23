$ErrorActionPreference = "Stop"

$BUS_NAME = "myna-test-bus"

# Helper function to get event bus ARN
function Get-EventBusArn ($name) {
    try {
        $arn = aws events describe-event-bus --name $name --query "Arn" --output text 2>$null
        if ($LASTEXITCODE -eq 0) { return $arn }
        return $null
    } catch {
        return $null
    }
}

Write-Host "Starting EventBridge Setup..." -ForegroundColor Cyan

# 1. Custom Event Bus
Write-Host "Checking Event Bus: $BUS_NAME"
$busArn = Get-EventBusArn $BUS_NAME
if (-not $busArn) {
    Write-Host "Creating Event Bus..."
    $busArn = aws events create-event-bus --name $BUS_NAME --query "EventBusArn" --output text
    Write-Host "Created Event Bus: $busArn" -ForegroundColor Green
} else {
    Write-Host "Event Bus exists."
}

Write-Host "EventBridge setup complete!" -ForegroundColor Green
