$ErrorActionPreference = "Stop"

$INSTANCE_NAME = "myna-test-instance"

Write-Host "Starting EC2 Destruction..." -ForegroundColor Cyan

# 1. Find Instance
$instanceId = aws ec2 describe-instances --filters "Name=tag:Name,Values=$INSTANCE_NAME" "Name=instance-state-name,Values=pending,running,stopped,stopping" --query "Reservations[].Instances[].InstanceId" --output text

if ($instanceId) {
    Write-Host "Terminating Instance: $instanceId"
    aws ec2 terminate-instances --instance-ids $instanceId | Out-Null
    
    Write-Host "Waiting for termination..."
    aws ec2 wait instance-terminated --instance-ids $instanceId
    Write-Host "Instance terminated." -ForegroundColor Green
} else {
    Write-Host "Instance not found."
}

Write-Host "EC2 destruction complete!" -ForegroundColor Green
