$ErrorActionPreference = "Stop"

$DB_ID = "myna-test-db"

Write-Host "Starting RDS Destruction..." -ForegroundColor Cyan

try {
    Write-Host "Deleting DB Instance: $DB_ID"
    aws rds delete-db-instance `
        --db-instance-identifier $DB_ID `
        --skip-final-snapshot `
        --delete-automated-backups `
        | Out-Null
        
    Write-Host "Deletion initiated." -ForegroundColor Green
} catch {
    Write-Host "DB Instance not found or already deleting."
}

Write-Host "RDS destruction initiated!" -ForegroundColor Green
