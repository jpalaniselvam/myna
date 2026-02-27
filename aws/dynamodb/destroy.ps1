$TableUsers = "myna-test-users"
$TableOrders = "myna-test-orders"

Write-Host "Destroying DynamoDB resources..."

aws dynamodb delete-table --table-name $TableUsers 2>$null
aws dynamodb delete-table --table-name $TableOrders 2>$null

Write-Host "DynamoDB resources destroyed."
