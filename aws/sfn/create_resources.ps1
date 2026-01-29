$ErrorActionPreference = "Stop"

$SM_NAME = "myna-test-state-machine"
$ROLE_NAME = "myna-test-sfn-role"

Write-Host "Starting Step Functions Setup..." -ForegroundColor Cyan

# 1. Create Role
Write-Host "Checking Role: $ROLE_NAME"
$roleArn = $null
try {
    $roleArn = aws iam get-role --role-name $ROLE_NAME --query "Role.Arn" --output text 2>$null
} catch { }

if (-not $roleArn) {
    Write-Host "Creating Role..."
    $trustPolicy = @{
        Version = "2012-10-17"
        Statement = @(
            @{
                Effect = "Allow"
                Principal = @{ Service = "states.amazonaws.com" }
                Action = "sts:AssumeRole"
            }
        )
    } | ConvertTo-Json -Depth 3
    
    $trustPolicy | Out-File "sfn-trust.json" -Encoding ascii
    $roleArn = aws iam create-role --role-name $ROLE_NAME --assume-role-policy-document file://sfn-trust.json --query "Role.Arn" --output text
    Remove-Item "sfn-trust.json"
    Write-Host "Role Created: $roleArn"
}

# 2. Create State Machine
Write-Host "Checking State Machine: $SM_NAME"
$smArn = aws stepfunctions list-state-machines --query "stateMachines[?name=='$SM_NAME'].stateMachineArn" --output text

if (-not $smArn) {
    Write-Host "Creating State Machine..."
    $definition = @{
        Comment = "A Hello World example of the Amazon States Language using a Pass state"
        StartAt = "HelloWorld"
        States = @{
            HelloWorld = @{
                Type = "Pass"
                End = $true
            }
        }
    } | ConvertTo-Json -Depth 5
    
    $definition | Out-File "sfn-definition.json" -Encoding ascii
    
    $smArn = aws stepfunctions create-state-machine `
        --name $SM_NAME `
        --definition file://sfn-definition.json `
        --role-arn $roleArn `
        --query "stateMachineArn" `
        --output text
        
    Remove-Item "sfn-definition.json"
        
    Write-Host "Created State Machine: $smArn" -ForegroundColor Green
} else {
    Write-Host "State Machine exists: $smArn"
}

Write-Host "Step Functions setup complete!" -ForegroundColor Green
