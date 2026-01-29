param (
    [Parameter(Mandatory = $true)]
    [string]$EmailIdentity
)

$ErrorActionPreference = "Stop"

$EMAIL_IDENTITY = $EmailIdentity

Write-Host "Starting SES Cleanup..." -ForegroundColor Cyan
Write-Host "Deleting email identity: $EMAIL_IDENTITY"

aws ses delete-identity --identity $EMAIL_IDENTITY

Write-Host "SES cleanup complete!" -ForegroundColor Green
