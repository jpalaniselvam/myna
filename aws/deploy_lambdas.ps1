$ErrorActionPreference = "Stop"

# Configuration
$ROLE_NAME = "myna-test-lambda-role"
$LAMBDA_DIR = ".\lambda"
$RUNTIME = "nodejs20.x"

Write-Host "Starting deployment of Myna test lambdas..." -ForegroundColor Cyan

# 1. Check/Create IAM Role
Write-Host "Checking IAM Role: $ROLE_NAME"
$roleExists = $false
try {
    $role = aws iam get-role --role-name $ROLE_NAME 2>$null
    if ($LASTEXITCODE -eq 0) { $roleExists = $true }
} catch {
    $roleExists = $false
}

if (-not $roleExists) {
    Write-Host "Creating role..."
    $trustPolicy = @{
        Version = "2012-10-17"
        Statement = @(
            @{
                Effect = "Allow"
                Principal = @{ Service = "lambda.amazonaws.com" }
                Action = "sts:AssumeRole"
            }
        )
    } | ConvertTo-Json -Depth 3
    
    $trustPolicy | Out-File "trust-policy.json" -Encoding ascii
    
    aws iam create-role --role-name $ROLE_NAME --assume-role-policy-document file://trust-policy.json | Out-Null
    aws iam attach-role-policy --role-name $ROLE_NAME --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
    
    Remove-Item "trust-policy.json"
    Write-Host "Role created. Waiting 10s for propagation..." -ForegroundColor Green
    Start-Sleep -Seconds 10
} else {
    Write-Host "Role exists."
}

$roleArn = (aws iam get-role --role-name $ROLE_NAME --query 'Role.Arn' --output text).Trim()

# 2. Deploy Functions
$files = Get-ChildItem -Path $LAMBDA_DIR -Filter "*.js"

foreach ($file in $files) {
    $funcName = $file.BaseName
    $zipFile = "$funcName.zip"
    
    Write-Host "`nProcessing $funcName..." -ForegroundColor Cyan
    
    # Compress-Archive requires full paths usually to be safe, or relative to current location
    if (Test-Path $zipFile) { Remove-Item $zipFile }
    Compress-Archive -Path $file.FullName -DestinationPath $zipFile
    
    # Check if function exists
    $funcExists = $false
    try {
        aws lambda get-function --function-name $funcName 2>$null | Out-Null
        if ($LASTEXITCODE -eq 0) { $funcExists = $true }
    } catch {
        $funcExists = $false
    }

    if (-not $funcExists) {
        Write-Host "Creating function..."
        # Node handler format is usually filename.exportName
        # Here we assume the file 'context-handler.js' exports 'handler', so 'context-handler.handler'
        aws lambda create-function `
            --function-name $funcName `
            --runtime $RUNTIME `
            --role $roleArn `
            --handler "$funcName.handler" `
            --zip-file "fileb://$zipFile" `
            --timeout 15 `
            | Out-Null
        Write-Host "Created $funcName" -ForegroundColor Green
    } else {
        Write-Host "Updating function code..."
        aws lambda update-function-code `
            --function-name $funcName `
            --zip-file "fileb://$zipFile" `
            | Out-Null
        Write-Host "Updated $funcName" -ForegroundColor Green
    }
    
    Remove-Item $zipFile
}

Write-Host "`nDeployment Complete!" -ForegroundColor Green
