$TableUsers = "myna-test-users"
$TableOrders = "myna-test-orders"

Write-Host "Starting DynamoDB Setup..."

# 1. Users Table
Write-Host "Checking Table: $TableUsers"
$table = aws dynamodb describe-table --table-name $TableUsers 2>$null | ConvertFrom-Json
if ($null -eq $table) {
    Write-Host "Creating Users Table..."
    aws dynamodb create-table `
        --table-name $TableUsers `
        --attribute-definitions `
            AttributeName=UserId,AttributeType=S `
        --key-schema `
            AttributeName=UserId,KeyType=HASH `
        --provisioned-throughput `
            ReadCapacityUnits=5,WriteCapacityUnits=5
    aws dynamodb wait table-exists --table-name $TableUsers
    Write-Host "Created Users Table."
} else {
    Write-Host "Users Table exists."
}

# 2. Orders Table
Write-Host "Checking Table: $TableOrders"
$table = aws dynamodb describe-table --table-name $TableOrders 2>$null | ConvertFrom-Json
if ($null -eq $table) {
    Write-Host "Creating Orders Table..."
    aws dynamodb create-table `
        --table-name $TableOrders `
        --attribute-definitions `
            AttributeName=UserId,AttributeType=S `
            AttributeName=OrderId,AttributeType=S `
        --key-schema `
            AttributeName=UserId,KeyType=HASH `
            AttributeName=OrderId,KeyType=RANGE `
        --provisioned-throughput `
            ReadCapacityUnits=5,WriteCapacityUnits=5
    aws dynamodb wait table-exists --table-name $TableOrders
    Write-Host "Created Orders Table."
} else {
    Write-Host "Orders Table exists."
}

Write-Host "DynamoDB setup complete!"
