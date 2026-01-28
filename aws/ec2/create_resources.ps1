$ErrorActionPreference = "Stop"

$INSTANCE_NAME = "myna-test-instance"
$KEY_NAME = "myna-test-key"

Write-Host "Starting EC2 Setup..." -ForegroundColor Cyan

# 1. Get AMI (Amazon Linux 2023)
Write-Host "Fetching latest Amazon Linux 2023 AMI..."
$amiId = aws ssm get-parameters --names /aws/service/ami-amazon-linux-latest/al2023-ami-kernel-default-x86_64 --query "Parameters[0].Value" --output text
Write-Host "Using AMI: $amiId"

# 2. Check for existing instance
$instanceId = aws ec2 describe-instances --filters "Name=tag:Name,Values=$INSTANCE_NAME" "Name=instance-state-name,Values=pending,running,stopped,stopping" --query "Reservations[].Instances[].InstanceId" --output text

if (-not $instanceId) {
    Write-Host "Creating Instance..."
    $instanceId = aws ec2 run-instances `
        --image-id $amiId `
        --instance-type t3.micro `
        --tag-specifications "ResourceType=instance,Tags=[{Key=Name,Value=$INSTANCE_NAME}]" `
        --query "Instances[0].InstanceId" `
        --output text
    
    Write-Host "Created Instance: $instanceId" -ForegroundColor Green
    
    Write-Host "Waiting for instance to be running..."
    aws ec2 wait instance-running --instance-ids $instanceId
    Write-Host "Instance is running." -ForegroundColor Green
} else {
    Write-Host "Instance exists: $instanceId"
}

Write-Host "EC2 setup complete!" -ForegroundColor Green
