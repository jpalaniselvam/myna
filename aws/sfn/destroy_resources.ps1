$ErrorActionPreference = "Stop"

$SM_NAME = "myna-test-state-machine"
$ROLE_NAME = "myna-test-sfn-role"

Write-Host "Starting Step Functions Destruction..." -ForegroundColor Cyan

# 1. Delete State Machine
$smArn = aws stepfunctions list-state-machines --query "stateMachines[?name=='$SM_NAME'].stateMachineArn" --output text
if ($smArn) {
    Write-Host "Deleting State Machine: $smArn"
    aws stepfunctions delete-state-machine --state-machine-arn $smArn
}

# 2. Delete Role
try {
    Write-Host "Deleting Role: $ROLE_NAME"
    aws iam delete-role --role-name $ROLE_NAME 2>$null
} catch {
    Write-Host "Role not found or error deleting."
}

Write-Host "Step Functions destruction complete!" -ForegroundColor Green
