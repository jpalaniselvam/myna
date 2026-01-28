$ErrorActionPreference = "Stop"

$DB_ID = "myna-test-db"
$DB_USER = "postgres"
$DB_PASS = "mypassword123"

Write-Host "Starting RDS Setup..." -ForegroundColor Cyan

# 1. Check DB
$dbStatus = $null
try {
    $dbStatus = aws rds describe-db-instances --db-instance-identifier $DB_ID --query "DBInstances[0].DBInstanceStatus" --output text 2>$null
} catch {
    $dbStatus = $null
}

if (-not $dbStatus) {
    Write-Host "Creating DB Instance (this will take several minutes)..."
    aws rds create-db-instance `
        --db-instance-identifier $DB_ID `
        --db-instance-class db.t3.micro `
        --engine postgres `
        --master-username $DB_USER `
        --master-user-password $DB_PASS `
        --allocated-storage 20 `
        --publicly-accessible `
        --tags "Key=Project,Value=myna-test" `
        | Out-Null
        
    Write-Host "DB Instance creation initiated. use 'aws rds wait db-instance-available' to wait." -ForegroundColor Yellow
} else {
    Write-Host "DB Instance exists: $dbStatus"
}

Write-Host "RDS setup initiated!" -ForegroundColor Green
