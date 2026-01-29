param (
    [Parameter(Mandatory = $true)]
    [string]$EmailIdentity
)

$ErrorActionPreference = "Stop"

$EMAIL_IDENTITY = $EmailIdentity

Write-Host "Starting SES Setup..." -ForegroundColor Cyan
Write-Host "Using email identity: $EMAIL_IDENTITY"

# Helper function to check identity existence
function Get-Identity ($email) {
    try {
        # This returns JSON list of identities
        $ids = aws ses list-identities --identity-type EmailAddress --query "Identities" --output json 2>$null | ConvertFrom-Json
        if ($ids -contains $email) {
            return $true
        }
        return $false
    } catch {
        return $false
    }
}

if (-not (Get-Identity $EMAIL_IDENTITY)) {
    Write-Host "Creating/Verifying Email Identity: $EMAIL_IDENTITY"
    aws ses verify-email-identity --email-address $EMAIL_IDENTITY
    Write-Host "Verification email sent to $EMAIL_IDENTITY. Please verify it to send emails in Sandbox mode." -ForegroundColor Yellow
} else {
    Write-Host "Identity $EMAIL_IDENTITY already exists."
}

Write-Host "SES setup complete!" -ForegroundColor Green
